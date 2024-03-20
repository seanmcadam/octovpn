package connection

import (
	"errors"
	"io"
	"net"

	"github.com/seanmcadam/bufferpool"
	"github.com/seanmcadam/counter"
	"github.com/seanmcadam/ctx"
	"github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/common"
	"github.com/seanmcadam/octovpn/interfaces"
)

var connCount counter.Counter

func init() {
	connCount = counter.New(ctx.New(), counter.BIT32)
}

type Conn struct {
	serial   counter.Count
	cx       *ctx.Ctx
	conn     net.Conn
	recvch   chan *bufferpool.Buffer
	sendch   chan *bufferpool.Buffer
	statusch chan common.LayerStatus
}

var pool bufferpool.Pool

func init() {
	pool = *bufferpool.New()
}

func Connection(cx *ctx.Ctx, conn net.Conn) interfaces.LayerInterface {

	t := &Conn{
		serial:   connCount.Next(),
		cx:       cx,
		conn:     conn,
		recvch:   make(chan *bufferpool.Buffer, 5),
		sendch:   make(chan *bufferpool.Buffer, 5),
		statusch: make(chan common.LayerStatus, 2),
	}

	loggy.Debugf("[%d] NEW %s:%s:%s", t.serial.Uint(), t.conn.LocalAddr().Network(), t.conn.LocalAddr().String(), t.conn.RemoteAddr().String())

	//err := t.conn.SetDeadline(time.Now())
	//if err != nil {
	//	loggy.Debugf("[%d] SetDeadLine() Error %s:%s:%s %s", t.serial.Uint(), t.conn.LocalAddr().Network(), t.conn.LocalAddr().String(), t.conn.RemoteAddr().String(), err)
	//	return nil
	//}

	go func(t *Conn) {
		defer func() {
			loggy.Debugf("[%d] Calling Closing Defer() %s:%s:%s", t.serial.Uint(), t.conn.LocalAddr().Network(), t.conn.LocalAddr().String(), t.conn.RemoteAddr().String())
			cx.Cancel()
			t.conn.Close()
			close(t.sendch)
			close(t.statusch)
		}()

		for {
			loggy.Debugf("[%d] mainloop called %s:%s:%s", t.serial.Uint(), t.conn.LocalAddr().Network(), t.conn.LocalAddr().String(), t.conn.RemoteAddr().String())

			select {
			case <-t.cx.DoneChan():
				loggy.Debugf("[%d] DoneChan() called %s:%s:%s", t.serial.Uint(), t.conn.LocalAddr().Network(), t.conn.LocalAddr().String(), t.conn.RemoteAddr().String())
				return
			case buf := <-t.sendch:
				if buf == nil {
					loggy.Debugf("[%d] sendch closed %s:%s:%s", t.serial.Uint(), t.conn.LocalAddr().Network(), t.conn.LocalAddr().String(), t.conn.RemoteAddr().String())
					return
				}

				loggy.Debugf("[%d] Got Send() Buffer[%d]", t.serial.Uint(), buf.Serial())

				if buf.Size() > 0 {
					n, err := t.conn.Write(buf.Data())
					if err != nil {
						loggy.Debugf("[%s] Write Error: %s", t.conn.RemoteAddr(), err)
						buf.ReturnToPool()
						return
					}
					if n != buf.Size() {
						loggy.Debugf("[%s] Size mismatch: Want:%d Size:%d", t.conn.RemoteAddr(), buf.Size(), n)
						buf.ReturnToPool()
						return
					}
				} else {
					loggy.Debugf("[%s]Send Buffer Size Zero", t.conn.RemoteAddr())
				}
				buf.ReturnToPool()
			}
		}
	}(t)

	go t.goRecv()

	return interfaces.LayerInterface(t)
}

func (t *Conn) goRecv() {
	loggy.Debugf("[%d] goRecv() %s:%s:%s", t.serial.Uint(), t.conn.LocalAddr().Network(), t.conn.LocalAddr().String(), t.conn.RemoteAddr().String())

	defer func() {
		loggy.Debugf("[%d] Recv Closing Defer() %s:%s:%s", t.serial.Uint(), t.conn.LocalAddr().Network(), t.conn.LocalAddr().String(), t.conn.RemoteAddr().String())
		t.cx.Cancel()
		close(t.recvch)
	}()

	//var buffer bytes.Buffer
	//tmp := make([]byte, 2048)

INNERLOOP:
	for {
		//err := t.conn.SetReadDeadline(time.Now().Add(time.Second))
		//if err != nil {
		//	loggy.FatalfStack("[%d] SetDeadLine() Error %s:%s:%s %s", t.serial.Uint(), t.conn.LocalAddr().Network(), t.conn.LocalAddr().String(), t.conn.RemoteAddr().String(), err)
		//	//return
		//}
		//
		// Load the receive buffer
		//
		b := make([]byte, 2048)
		n, err := t.conn.Read(b)
		if err != nil {
			if err == io.EOF || isNetClosingError(err) {
				loggy.Debugf("[%d] Read() %s:%s:%s Error: %s", t.serial.Uint(), t.conn.LocalAddr().Network(), t.conn.LocalAddr().String(), t.conn.RemoteAddr(), err)
			} else {
				loggy.FatalfStack("Read() Error:%s on %s", err, t.conn.RemoteAddr())
			}
			loggy.Debugf("[%d]ERROR Read() %s:%s:%s >> %s", t.serial.Uint(), t.conn.LocalAddr().Network(), t.conn.LocalAddr().String(), t.conn.RemoteAddr(), err)
			return
		}

		if n == 0 {
			loggy.Debugf("[%d] Zero Read() %s:%s:%s", t.serial.Uint(), t.conn.LocalAddr().Network(), t.conn.LocalAddr().String(), t.conn.RemoteAddr())
			continue INNERLOOP
		}

		loggy.Debugf("[%d] Read() %s:%s:%s Size=%d", t.serial.Uint(), t.conn.LocalAddr().Network(), t.conn.LocalAddr().String(), t.conn.RemoteAddr(), n)
		buf := pool.Get()
		buf.Append(b[:n])

		select {
		case <-t.cx.DoneChan():
			loggy.Debugf("[%d] DoneCh() Read() %s:%s:%s Size=%d", t.serial.Uint(), t.conn.LocalAddr().Network(), t.conn.LocalAddr().String(), t.conn.RemoteAddr(), n)
			return
		case t.recvch <- buf:
		default:
			loggy.Debugf("[%d] Failed recvch Read() %s:%s:%s Size=%d", t.serial.Uint(), t.conn.LocalAddr().Network(), t.conn.LocalAddr().String(), t.conn.RemoteAddr(), n)
			buf.ReturnToPool()
		}
	}
}

func (t *Conn) Send(b *bufferpool.Buffer) {
	if t == nil {
		loggy.Panicf("nil method pointer")
	}
	if b == nil {
		loggy.Panicf("nil data pointer")
	}

	loggy.Debugf("[%d] Send() Buffer[%d]", t.serial.Uint(), b.Serial())
	select {

	case <-t.cx.DoneChan():
		loggy.Debugf("[%d] Closed, Dropping Buffer %s:%s:%s", t.serial.Uint(), t.conn.LocalAddr().Network(), t.conn.LocalAddr().String(), t.conn.RemoteAddr().String())
		return
	default:
	}

	select {
	case <-t.cx.DoneChan():
		loggy.Debugf("[%d] DoneCh() Send() %s:%s:%s", t.serial.Uint(), t.conn.LocalAddr().Network(), t.conn.LocalAddr().String(), t.conn.RemoteAddr())
		return
	case t.sendch <- b:
		loggy.Debugf("Send() Data:'%s'", string(b.Data()))
	default:
		loggy.Debugf("[%d] Failed Send() Buffer[%d]", t.serial.Uint(), b.Serial())
		b.ReturnToPool()
	}

}

func (t *Conn) Reset() {
	t.cx.Cancel()
}

func (t *Conn) RecvCh() (ch chan *bufferpool.Buffer) {
	return t.recvch
}

func (t *Conn) StatusCh() (statusch chan common.LayerStatus) {
	return t.statusch
}

func isNetClosingError(err error) bool {
	var netOpError *net.OpError
	if errors.As(err, &netOpError) {
		return netOpError.Err.Error() == "use of closed network connection"
	}
	return false
}
