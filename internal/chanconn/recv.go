package chanconn

import (
	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/internal/packet/packetchan"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (cs *ChanconnStruct) RecvChan() <-chan interfaces.PacketInterface {
	return cs.recvch
}

func (cs *ChanconnStruct) goRecv() {
	for {
		select {
		case <-cs.cx.DoneChan():
			return

		case p := <-cs.conn.RecvChan():
			if p == nil {
				log.Debug("Got nil packet")
				return
			}

			switch p.Type() {
			case packet.CONN_TYPE_PARENT:
				payload, err := packetchan.MakePacket(p.Payload().([]byte))
				if err != nil {
					log.Errorf("packetchan.MakePacket() Err:%s", err)
					continue
				}

				cs.recvch <- payload

			case packet.CONN_TYPE_PING64:
				cs.send(p.CopyPong64())

			case packet.CONN_TYPE_PONG64:
				cs.pinger.Pongch <- p.Payload().(counter.Counter64)

			default:
				log.Fatalf("Unhandled Packet Type:%d", p.Type())
			}
		}
	}
}
