package tcpcli

import (
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (t *TcpClientStruct) RecvChan() <-chan *packet.PacketStruct {

	if t == nil {
		log.FatalStack("nil TcpClientStruct")
		return nil
	}
	if t.tcpconn == nil {
		return nil
	}

	return t.tcpconn.RecvChan()
}
