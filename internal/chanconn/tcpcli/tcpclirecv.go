package tcpcli

import (
	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (t *TcpClientStruct) RecvChan() <-chan interfaces.PacketInterface {

	if t == nil {
		log.FatalStack("nil TcpClientStruct")
		return nil
	}
	if t.tcpconn == nil {
		return nil
	}

	return t.tcpconn.RecvChan()
}
