package tcpsrv

import (
	"github.com/seanmcadam/octovpn/internal/chanconn/tcp"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (t *TcpServerStruct) goNewConn(tcp *tcp.TcpStruct) {

	select {
	case <-tcp.Closech:
		log.Info("TCPSrv connection closing down")
	}

}
