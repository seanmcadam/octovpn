package udpcli

import (
	"github.com/seanmcadam/octovpn/interfaces"
)

func (u *UdpClientStruct) RecvChan() <-chan interfaces.PacketInterface {
	return u.udpconn.RecvChan()
}
