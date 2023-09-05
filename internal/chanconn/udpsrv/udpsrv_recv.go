package udpsrv

import (
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (t *UdpServerStruct) RecvChan() <-chan *packet.PacketStruct {
	if t.udpconn == nil {
		log.Error("udpconn is nil")
		return nil
	}
	return t.udpconn.RecvChan()
}
