package udpsrv

import (
	"github.com/seanmcadam/octovpn/interfaces"
)

func (t *UdpServerStruct) RecvChan() <-chan interfaces.PacketInterface {
	return t.udpconn.RecvChan()
}
