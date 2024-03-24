package udp

import (
	"errors"
	"fmt"
	"net"
	"net/netip"

	"github.com/seanmcadam/bufferpool"
	"github.com/seanmcadam/counter"
	"github.com/seanmcadam/ctx"
	"github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/common/status"
)

var connCount counter.Counter

func init() {
	connCount = counter.New(ctx.New(), counter.BIT32)
}

//type ConnectionInterface interface {
//	Send()
//	Recv()
//}

type UDPConnection struct {
	serial   counter.Count
	id       string
	cx       *ctx.Ctx
	conn     *net.UDPConn
	raddr    *net.UDPAddr
	addrport netip.AddrPort
	recvch   chan *bufferpool.Buffer
	sendch   chan *bufferpool.Buffer
	status   *status.LayerStatusStruct
	remove   func(netip.AddrPort)
}

var pool bufferpool.Pool

func init() {
	pool = *bufferpool.New()
}

func NewConnection(cx *ctx.Ctx, conn *net.UDPConn, udpport *net.UDPAddr, remove func(netip.AddrPort)) (uc *UDPConnection) {

	uc = &UDPConnection{
		serial:   connCount.Next(),
		cx:       cx,
		conn:     conn,
		raddr:    udpport,
		addrport: udpport.AddrPort(),
		recvch:   make(chan *bufferpool.Buffer, 5),
		sendch:   make(chan *bufferpool.Buffer, 5),
		status:   status.New(cx),
		remove:   remove,
	}

	var connType string = "Server"
	if remove == nil {
		connType = "Client"
	}

	uc.id = fmt.Sprintf("[%s[%d]%s:%s->%s]", connType, uc.serial.Uint(), uc.conn.LocalAddr().Network(), uc.conn.LocalAddr().String(), uc.raddr.String())

	loggy.Debugf("%s: connection()", uc.id)

	go uc.goSend()
	// go uc.goRecv()

	return uc
}

func (uc *UDPConnection) goSend() {
	if uc == nil {
		loggy.Panicf("nil method pointer")
	}

	var n int
	var err error

	defer func() {
		loggy.Debugf("%s: Defer()", uc.id)
		uc.cx.Cancel()
		if uc.remove != nil {
			uc.remove(uc.addrport)
		}
		close(uc.sendch)
	}()

	for {
		loggy.Debugf("%s: mainloop", uc.id)

		select {
		case <-uc.cx.DoneChan():
			loggy.Debugf("%s: DoneChan()", uc.id)
			return
		case buf := <-uc.sendch:
			if buf == nil {
				loggy.Debugf("%s: sendch closed", uc.id)
				return
			}

			loggy.Debugf("%s: Send() Buf:%d", uc.id, buf.Serial())

			if uc.conn.RemoteAddr() == nil {
				n, err = uc.conn.WriteToUDP(buf.Data(), uc.raddr)
				if err != nil {
					loggy.FatalfStack("%s: WriteToUDP() Error:%s", uc.id, err)
				}
				loggy.Debugf("%s: WriteToUDP() Buf:%d", uc.id, buf.Serial())
			} else {
				n, err = uc.conn.Write(buf.Data())
				if err != nil {
					loggy.FatalfStack("%s: Write() Error:%s", uc.id, err)
				}
				loggy.Debugf("%s: Write() Buf:%d", uc.id, buf.Serial())
			}
			if n != buf.Size() {
				loggy.FatalfStack("%s: WriteToUDP() Size Error:%d/%d", uc.id, n, buf.Size())
			}

			if buf.Used() == false {
				loggy.ErrorfStack("%s: Buf Used = false", uc.id)
			}

			buf.ReturnToPool()
		}
	}
}

func (uc *UDPConnection) goRecv() {
	if uc == nil {
		loggy.Panicf("nil method pointer")
	}

	defer func() {
		loggy.Debugf("%s: Defer()", uc.id)
		uc.cx.Cancel()
	}()

	for {
		loggy.Debugf("%s: for loop", uc.id)
		var buf *bufferpool.Buffer = nil
		select {
		case buf = <-uc.recvch:
		case <-uc.cx.DoneChan():
			loggy.Debugf("%s: DoneCh()", uc.id)
			return
		}

		if buf == nil {
			loggy.Debugf("%s: NIL Buf", uc.id)
		}

		select {
		case uc.recvch <- buf:
		case <-uc.cx.DoneChan():
			loggy.Debugf("%s: DoneCh()", uc.id)
			return
		}
	}

}

func (uc *UDPConnection) PushRecv(b *bufferpool.Buffer) {
	if uc == nil {
		loggy.Panicf("nil method pointer")
	}
	if b == nil {
		loggy.Panicf("nil data pointer")
	}

	// loggy.Debugf("%s: Buf[%d] %v", uc.id, b.Serial(), b)
	select {
	case uc.recvch <- b:
	case <-uc.cx.DoneChan():
		loggy.Debugf("%s: DoneChan()", uc.id)
		b.ReturnToPool()
	default:
		loggy.Errorf("%s: Dropped Buf[%d]", uc.id, b.Serial())
		b.ReturnToPool()
	}
}

func (uc *UDPConnection) Send(b *bufferpool.Buffer) {
	if uc == nil {
		loggy.Panicf("nil method pointer")
	}
	if b == nil {
		loggy.Panicf("nil data pointer")
	}

	loggy.Debugf("%s: Buf[%d]", uc.id, b.Serial())
	select {

	case <-uc.cx.DoneChan():
		loggy.Debugf("%s: DoneChan()", uc.id)
		return
	default:
	}

	select {
	case <-uc.cx.DoneChan():
		loggy.Debugf("%s: DoneCh", uc.id)
		return
	case uc.sendch <- b:
	default:
		loggy.Debugf("%s: Failed to send to sendch", uc.id)
		b.ReturnToPool()
	}

}

func (uc *UDPConnection) Reset() {
	uc.cx.Cancel()
}

func (uc *UDPConnection) RecvCh() (ch chan *bufferpool.Buffer) {
	return uc.recvch
}

func (uc *UDPConnection) Status() (status status.LayerStatus) {
	return uc.status.Get()
}

func (uc *UDPConnection) StatusCh() (statusch chan status.LayerStatus) {
	return uc.status.GetCh()
}

func isNetClosingError(err error) bool {
	var netOpError *net.OpError
	if errors.As(err, &netOpError) {
		return netOpError.Err.Error() == "use of closed network connection"
	}
	return false
}
