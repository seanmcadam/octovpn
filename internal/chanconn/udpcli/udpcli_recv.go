package udpcli

import (
	"github.com/seanmcadam/octovpn/internal/packet"
)

// RecvChan could be called when the pointer is nil - return a nil
func (u *UdpClientStruct) RecvChan() <-chan *packet.PacketStruct {
	if u == nil {
		return nil
	}
	return u.udpconn.RecvChan()
}
