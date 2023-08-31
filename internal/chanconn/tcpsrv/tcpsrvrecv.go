package tcpsrv

import (
	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (t *TcpServerStruct) RecvChan() <-chan interfaces.PacketInterface {

	if t == nil {
		log.FatalStack("nil TcpServerStruct")
		return nil
	}
	if t.tcpconn == nil {
		return nil
	}

	return t.tcpconn.RecvChan()
}
