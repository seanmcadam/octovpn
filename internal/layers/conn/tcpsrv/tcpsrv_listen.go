package tcpsrv

import (
	"strings"
	"time"

	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/internal/layers/network/tcp"
	"github.com/seanmcadam/octovpn/internal/msg"
)

func (t *TcpServerStruct) goListen() {
	if t == nil {
		return
	}

	defer t.Cancel()

	//t.link.Listen()
	t.setState(msg.StateLISTEN)

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
			cx := t.cx.NewWithCancel()
			log.Debug("TCPSrv New incoming connection")
			newconn := tcp.New(cx, conn)
			if newconn == nil {
				log.FatalStack("NewTCP is Nil")
			}

			if t.closed() {
				return
			}
		}
	}
}
