package chanconn

import (
	"github.com/seanmcadam/octovpn/octolib/counter"
	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/packet/packetchan"
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

func (cs *ChanconnStruct) RecvChan() <-chan *packetchan.ChanPacket {
	return cs.recvch
}

func (cs *ChanconnStruct) goRecv() {
	for {
		select {
		case <-cs.cx.DoneChan():
			return

		case packet := <-cs.conn.RecvChan():
			if packet == nil {
				log.Debug("Got nil packet")
				return
			}

			switch packet.GetType() {
			case packetconn.CONN_TYPE_CHAN:
				cs.recvch <- packet.GetPayload().(*packetchan.ChanPacket)

			case packetconn.CONN_TYPE_PING64:
				pong, err := packetconn.NewPacket(packetconn.CONN_TYPE_PONG64, packet.GetPayload())
				if err != nil {
					log.Errorf("NewPacket() PONG64 Err:%s", err)
				}
				cs.send(pong)

			case packetconn.CONN_TYPE_PONG64:
				cs.pinger.Pongch <- packet.GetPayload().(counter.Counter64)

			default:
				log.Fatalf("Unhandled Packet Type:%d", packet.GetType())
			}
		}
	}
}
