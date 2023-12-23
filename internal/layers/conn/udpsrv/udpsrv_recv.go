package udpsrv

import (
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/internal/packet"
)

func (u *UdpServerStruct) RecvChan() <-chan *packet.PacketStruct {
	if u == nil || u.udpconn == nil {
		return nil
	}

	if u.link.IsDown() {
		log.Debugf("UDP Srv Recv state:%s", u.link.GetState())
		return nil
	}

	return u.udpconn.RecvChan()
}
