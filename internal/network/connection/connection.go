package connection

import (
	"io"
	"net"

	"github.com/seanmcadam/bufferpool"
	"github.com/seanmcadam/ctx"
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/common"
	"github.com/seanmcadam/octovpn/interfaces"
)

var pool bufferpool.Pool

type Conn struct {
	cx       *ctx.Ctx
	conn     net.Conn
	recvch   chan *bufferpool.Buffer
	sendch   chan *bufferpool.Buffer
	statusch chan common.LayerStatus
}

func init() {
	pool = *bufferpool.New()
}

func connection(cx *ctx.Ctx, conn net.Conn) (t interfaces.LayerInterface) {
	t = &Conn{
		cx:       cx,
		conn:     conn,
		recvch:   make(chan *bufferpool.Buffer, 5),
		sendch:   make(chan *bufferpool.Buffer, 5),
		statusch: make(chan common.LayerStatus, 1),
	}

	go t.goRecv()

	go func(t *Conn) {
		defer func() {
			cx.Cancel()
			close(t.sendch)
			close(t.statusch)
		}()

		select {
		case <-t.cx.DoneChan():
			return
		case buf := <-t.sendch:
			n, err := t.conn.Write(buf.Data())
			if err != nil {
				log.Errorf("Conn[%s] Write Err: %s", t.conn.RemoteAddr(), err)
				buf.ReturnToPool()
				return
			}
			if n != buf.Size() {
				log.Errorf("Conn[%s] Size mismatch: Want:%d Sena:t%d", t.conn.RemoteAddr(), buf.Size(), n)
				buf.ReturnToPool()
				return
			}
			buf.ReturnToPool()
		}
	}(t.(*Conn))

	return t
}

func (t *Conn) goRecv() {
	for {
		func() {
			t.cx.Cancel()
			close(t.recvch)
		}()

		//var buffer bytes.Buffer
		//tmp := make([]byte, 2048)

		for {
			//
			// Load the receive buffer
			//
			buf := pool.Get()
			_, err := t.conn.Read(buf.Data())
			if err != nil {
				if err == io.EOF {
					log.Errorf("Read() connection closed %s", t.conn.RemoteAddr())
				} else {
					log.Errorf("Read() Error:%s on %s", err, t.conn.RemoteAddr())
				}
				return
			}

			//buffer.Write(tmp[:n])
			// buf := pool.Get()
			//buf.Append(buffer.Bytes())
			// buf.Append(tmp[:n])
			t.recvch <- buf

			//
			// Does the buffer have enough data to assemble a packet?
			//
			// sig, length, err := packet.ReadPacketBuffer(buffer.Bytes()[:6])

			// log.Debugf("RecvBuffer:%v",buffer.Bytes()[:n])
			//
			// Error checking types here
			//
			// if err != nil {
			// 	log.Errorf("MakePacket() Err:%s on %s", err, t.conn.RemoteAddr())
			// 	return
			// }

			//
			// Only receive CONN layer packets here
			//
			// if !sig.ConnLayer() {
			// 	log.Errorf("Bad SIG Layer Received:%s, on %s", sig, t.conn.RemoteAddr())
			// 	return
			// }

			//
			// Is there enough data?
			//
			// if buffer.Len() < int(length) {
			// 	log.Warnf("Not Enough Buffer Data %d < %d", buffer.Len(), int(length))
			// 	continue
			// }

			//
			// Extract a packet
			//
			// newpacketbuf := buffer.Next(int(length))
			// log.Debugf("Raw TCP Recv:%v", newpacketbuf)
			// p, err := packet.MakePacket(newpacketbuf)
			// if err != nil {
			// 	log.Errorf("MakePacket() Err:%s on %s", err, t.conn.RemoteAddr())
			// 	return
			// }

			// if p == nil {
			// 	log.Errorf("MakePacket() returned Nil Packet")
			// }

			// if p.Sig().Close() {
			// 	log.Debug("TCP received SOFT CLOSE")
			// 	return
			// }

			// packet := msg.NewPacket(t.me, p)
			// log.Debug("TCP Recv PAcket %v", packet)
		}
	}
}

func (t *Conn) Send(b *bufferpool.Buffer) {
	t.sendch <- b
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
