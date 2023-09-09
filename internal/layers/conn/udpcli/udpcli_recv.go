package udpcli

import (
	"github.com/seanmcadam/octovpn/internal/link"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// RecvChan could be called when the pointer is nil - return a nil
func (u *UdpClientStruct) RecvChan() <-chan *packet.PacketStruct {
	if u == nil || u.udpconn == nil {
		return nil
	}

	if u.link.GetState() != link.LinkStateUP {
		log.Debugf("UDP Cli Recv state:%s", u.link.GetState())
		return nil
	}

	return u.udpconn.RecvChan()
}
