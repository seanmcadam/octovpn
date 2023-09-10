package tcpsrv

import (
	"github.com/seanmcadam/octovpn/internal/link"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (t *TcpServerStruct) RecvChan() <-chan *packet.PacketStruct {
	if t == nil {
		log.ErrorStack("TCP Srv Recv Nil")
		return nil

	}
	if t.tcpconn == nil {
		log.Debug("TCP Srv nil tcpconn")
		return nil
	}

	if t.link.GetState() != link.LinkStateUP {
		log.Debugf("TCP Srv Recv state:%s", t.link.GetState())
		return nil
	}

	log.Debugf("TCP Srv Recv state:%s", t.link.GetState())
	return t.tcpconn.RecvChan()
}
