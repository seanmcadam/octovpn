package udpcli

import "github.com/seanmcadam/octovpn/internal/packet"

func (u *UdpClientStruct) RecvChan() <-chan *packet.PacketStruct {
	return u.udpconn.RecvChan()
}
