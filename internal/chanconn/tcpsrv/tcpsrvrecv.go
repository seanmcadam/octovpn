package tcpsrv

import (
	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (t *TcpServerStruct) RecvChan() <-chan interfaces.PacketInterface {

	if t == nil {
		log.FatalStack("nil TcpStruct")
		return nil
	}
	if t.recvch == nil {
		log.Error("Nil recvch pointer")
		return nil
	}

	return t.tcpconn.RecvChan()
}
