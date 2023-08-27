package tcpcli

import (
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

func (t *TcpClientStruct) RecvChan() <-chan *packetconn.ConnPacket {
	return t.tcpconn.RecvChan()
}
