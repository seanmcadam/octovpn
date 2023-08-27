package udpcli

import (
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

func (t *UdpClientStruct) RecvChan() <-chan *packetconn.ConnPacket {
	return t.udpconn.RecvChan()
}
