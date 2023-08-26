package tcpsrv

import (
	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/packet/packetchan"
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

func (t *TcpServerStruct) RecvChan() <-chan *packetchan.ChanPacket {
	return t.recvch
}

func (t *TcpServerStruct) goRecv() {
	for {
		select {
		case <-t.cx.DoneChan():
			return
		case data := <-t.tcpconn.RecvChan():
			if data == nil {
				// Closed connection...
				return
			}

			packettype := data.GetType()
			switch packettype {
			case packetconn.PACKET_TYPE_TCP:
				pc, err := packetchan.MakePacket(data.GetPayload())
				if err != nil {
					log.Fatalf("MakePacket Err:%s", err)
				}
				t.recvch <- pc
			case packetconn.PACKET_TYPE_PING:
				log.Debug("Ignore PING")
			case packetconn.PACKET_TYPE_PONG:
				log.Debug("Ignore PONG")
			}
		}
	}
}
