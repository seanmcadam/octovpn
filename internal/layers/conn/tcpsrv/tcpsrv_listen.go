package tcpsrv

import (
	"strings"
	"time"

	"github.com/seanmcadam/octovpn/internal/layers/network/tcp"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (t *TcpServerStruct) goListen() {
	if t == nil {
		return
	}

	defer t.Cancel()

	t.link.Listen()

	for {
		conn, err := t.tcplistener.AcceptTCP()
		if err != nil {
			// Assumed closed due to Cancel()
			if !strings.Contains(err.Error(), "use of closed network connection") {
				log.Warnf("AcceptTCP Closing Err:%s, but keep trying", err)
				time.Sleep(time.Second)
			} else {
				return
			}
			if t.closed() {
				return
			}
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
