package tcpsrv

import (
	"github.com/seanmcadam/octovpn/interfaces"
)

func (t *TcpServerStruct) RecvChan() <-chan interfaces.PacketInterface {
	return t.tcpconn.RecvChan()
}
