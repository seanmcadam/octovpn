package udpsrv

import (
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

func (t *UdpServerStruct) RecvChan() <-chan *packetconn.ConnPacket {
	return t.udpconn.RecvChan()
}
