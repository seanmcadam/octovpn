package tcpsrv

import (
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

func (t *TcpServerStruct) RecvChan() <-chan *packetconn.ConnPacket {
	return t.tcpconn.RecvChan()
}
