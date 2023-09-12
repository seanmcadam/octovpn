package tcpcli

import (
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (t *TcpClientStruct) RecvChan() <-chan *packet.PacketStruct {

	if t == nil || t.tcpconn == nil {
		log.Debugf("TCP Cli Recv Nil")
		return nil
	}

	if t.link.IsDown() {
		log.Debugf("TCP Cli Recv state:%s", t.link.GetState())
		return nil
	}

	return t.tcpconn.RecvChan()
}
