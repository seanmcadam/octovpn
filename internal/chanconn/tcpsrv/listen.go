package tcpsrv

import (
	"github.com/seanmcadam/octovpn/internal/chanconn/tcp"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (t *TcpServerStruct) goListen() {
	defer t.Close()

	for {
		conn, err := t.tcplistener.AcceptTCP()
		if err != nil {
			log.Debugf("AcceptTCP Error:%s", err)
			return
		}

		log.Debug("TCP New connection")
		newconn := tcp.NewTCP(conn)

		select {
		case t.tcpconnch <- newconn:
		case <-t.closech:
			log.Debug("TCPSrv listener close ch")
			return
		}
	}
}
