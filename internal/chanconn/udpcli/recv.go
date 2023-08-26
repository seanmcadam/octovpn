package udpcli

import (
	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/packet/packetchan"
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

func (t *UdpClientStruct) RecvChan() <-chan *packetchan.ChanPacket {
	return t.recvch
}

func (t *UdpClientStruct) goRecv() {
	for {
		select {
		case <-t.cx.DoneChan():
			return
		case data := <-t.udpconn.RecvChan():
			if data == nil {
				// Closed connection...
				return
			}

			packettype := data.GetType()
			switch packettype {
			case packetconn.PACKET_TYPE_UDP:
				payload := data.GetPayload()
				switch payload.(type) {
				case *packetchan.ChanPacket:
					t.recvch <- payload.(*packetchan.ChanPacket)
				default:
					log.Fatalf("Unexpected type:%t", payload)
				}
			case packetconn.PACKET_TYPE_PING:
				log.Debug("Ignore PING")
			case packetconn.PACKET_TYPE_PONG:
				log.Debug("Ignore PONG")
			default:
				log.Fatalf("Unhandled Type: %s", packettype)
			}
		}
	}
}
