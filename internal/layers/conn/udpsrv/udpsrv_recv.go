package udpsrv

import (
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
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
