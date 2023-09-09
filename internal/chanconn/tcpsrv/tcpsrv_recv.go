package tcpsrv

import (
	"github.com/seanmcadam/octovpn/internal/link"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (t *TcpServerStruct) RecvChan() <-chan *packet.PacketStruct {

	if t == nil || t.tcpconn == nil {
		log.Debugf("TCP Srv Recv Nil")
		log.Debugf("TCP Srv Recv state:%s", t.link.GetState())
		return nil
	}

	if t.link.GetState() != link.LinkStateUP {
		log.Debugf("TCP Srv Recv state:%s", t.link.GetState())
		return nil
	}

	log.Debugf("TCP Srv Recv state:%s", t.link.GetState())
	return t.tcpconn.RecvChan()
}
