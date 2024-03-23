package tcp

import (
	"errors"
	"fmt"
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

type TCPConnection struct {
	serial   counter.Count
	id       string
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

func connection(cx *ctx.Ctx, conn net.Conn) interfaces.LayerInterface {

	tc := &TCPConnection{
		serial:   connCount.Next(),
		cx:       cx,
		conn:     conn,
		recvch:   make(chan *bufferpool.Buffer, 5),
		sendch:   make(chan *bufferpool.Buffer, 5),
		statusch: make(chan common.LayerStatus, 2),
	}

	tc.id = fmt.Sprintf("[%d]%s:%s->%s", tc.serial.Uint(), tc.conn.LocalAddr().Network(), tc.conn.LocalAddr().String(), tc.conn.RemoteAddr().String())

	loggy.Debugf("%s:connection()", tc.id)

	go tc.goSend()
	go tc.goRecv()

	return interfaces.LayerInterface(tc)
}

func (tc *TCPConnection) goSend() {
	if tc == nil {
		loggy.Panicf("nil method pointer")
	}

	defer func() {
		loggy.Debugf("%s:goSend() Defer()", tc.id)
		tc.cx.Cancel()
		tc.conn.Close()
		close(tc.sendch)
		close(tc.statusch)
	}()

	for {
		loggy.Debugf("%s:goSend() mainloop", tc.id)

		select {
		case <-tc.cx.DoneChan():
			loggy.Debugf("%s:goSend() DoneChan()", tc.id)
			return
		case buf := <-tc.sendch:
			if buf == nil {
				loggy.Debugf("%s:goSend() sendch closed", tc.id)
				return
			}

			loggy.Debugf("%s:goSend Send() Buf:%d", tc.id, buf.Serial())

			if buf.Size() > 0 {
				n, err := tc.conn.Write(buf.Data())
				if err != nil {
					loggy.Debugf("%s:goSend() ERROR Send() Write %s", tc.id, err)
					buf.ReturnToPool()
					return
				}
				if n != buf.Size() {
					loggy.Debugf("%s:goSend() ERROR Send() Size mismatch %s", tc.id, buf.Size(), n)
					buf.ReturnToPool()
					return
				}
			} else {
				loggy.Debugf("%s:goSend ERROR Send() Zero Buf", tc.id)
			}
			buf.ReturnToPool()
		}
	}
}

func (tc *TCPConnection) goRecv() {
	if tc == nil {
		loggy.Panicf("nil method pointer")
	}

	defer func() {
		loggy.Debugf("%s:goRecv Defer()", tc.id)
		tc.cx.Cancel()
		close(tc.recvch)
	}()

INNERLOOP:
	for {
		loggy.Debugf("%s:goRecv mainloop", tc.id)
		//
		// Load the receive buffer
		//
		b := make([]byte, 2048)
		n, err := tc.conn.Read(b)
		if err != nil {
			if err == io.EOF || isNetClosingError(err) {
				loggy.Debugf("%s:goRecv Read Close Error %s", tc.id, err)
			} else {
				loggy.FatalfStack("%s:Read() Error:%s", tc.id, err)
			}
			loggy.Debugf("%s:goRecv Read Error %s", tc.id, err)
			return
		}

		if n == 0 {
			loggy.Debugf("%s:goRecv() Read Zero", tc.id)
			continue INNERLOOP
		}

		loggy.Debugf("%s:goRecv Read Size=%d", tc.id, n)
		buf := pool.Get()
		buf.Append(b[:n])

		select {
		case <-tc.cx.DoneChan():
			loggy.Debugf("%s:goRecv() DoneCh()", tc.id)
			buf.ReturnToPool()
			return
		case tc.recvch <- buf:
		default:
			loggy.Debugf("%s:goRecv() Failed to send buf()", tc.id)
			buf.ReturnToPool()
		}
	}
}

func (tc *TCPConnection) Send(b *bufferpool.Buffer) {
	if tc == nil {
		loggy.Panicf("nil method pointer")
	}
	if b == nil {
		loggy.Panicf("nil data pointer")
	}

	loggy.Debugf("%s:Send() Buf[%d]", tc.id, b.Serial())
	select {

	case <-tc.cx.DoneChan():
		loggy.Debugf("%s:Send() DoneCh", tc.id)
		return
	default:
	}

	select {
	case <-tc.cx.DoneChan():
		loggy.Debugf("%s:Send() DoneCh", tc.id)
		return
	case tc.sendch <- b:
	default:
		loggy.Debugf("%s:Send() Failed to send to sendch", tc.id)
		b.ReturnToPool()
	}

}

func (tc *TCPConnection) Reset() {
	tc.cx.Cancel()
}

func (tc *TCPConnection) RecvCh() (ch chan *bufferpool.Buffer) {
	return tc.recvch
}

func (tc *TCPConnection) StatusCh() (statusch chan common.LayerStatus) {
	return tc.statusch
}

func isNetClosingError(err error) bool {
	var netOpError *net.OpError
	if errors.As(err, &netOpError) {
		return netOpError.Err.Error() == "use of closed network connection"
	}
	return false
}
