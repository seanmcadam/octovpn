package tcpcli

import (
	"github.com/seanmcadam/octovpn/interfaces"
)

func (t *TcpClientStruct) RecvChan() <-chan interfaces.PacketInterface {
	return t.tcpconn.RecvChan()
}
