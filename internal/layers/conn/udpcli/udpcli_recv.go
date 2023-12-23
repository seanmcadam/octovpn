package udpcli

import (
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/internal/packet"
)

// RecvChan could be called when the pointer is nil - return a nil
func (u *UdpClientStruct) RecvChan() <-chan *packet.PacketStruct {
	if u == nil || u.udpconn == nil {
		return nil
	}

	if u.link.IsDown() {
		log.Debugf("UDP Cli Recv state:%s", u.link.GetState())
		return nil
	}

	return u.udpconn.RecvChan()
}
