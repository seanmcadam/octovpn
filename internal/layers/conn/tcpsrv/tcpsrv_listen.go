package tcpsrv

import (
	"strings"

	"github.com/seanmcadam/octovpn/internal/layers/network/tcp"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (t *TcpServerStruct) goListen() {
	if t == nil {
		return
	}

	defer t.Cancel()

	for {
		conn, err := t.tcplistener.AcceptTCP()
		if err != nil {
			// Assumed closed due to Cancel()
			if !strings.Contains(err.Error(), "use of closed network connection") {
				log.Errorf("AcceptTCP Closing:%s", err)
			}
			return
		}

		if conn != nil {
			log.Debug("TCPSrv New incoming connection")
			newconn := tcp.NewTCP(t.cx.NewWithCancel(), conn)
			if newconn == nil {
				log.FatalStack("NewTCP is Nil")
			}
			t.tcpconnch <- newconn

			if t.closed() {
				return
			}
		}
	}
}
