package tcpsrv

import (
	"github.com/seanmcadam/octovpn/internal/chanconn/tcp"
	"github.com/seanmcadam/octovpn/internal/link"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (t *TcpServerStruct) goListen() {
	defer t.cx.Done()

	for {
		conn, err := t.tcplistener.AcceptTCP()
		if err != nil {
			log.Debugf("AcceptTCP Error:%s", err)
			return
		}

		log.Debug("TCP New connection")
		newconn := tcp.NewTCP(t.cx.NewWithCancel(), conn)
		t.tcpconnch <- newconn

		t.link.ToggleState(link.LinkStateUp)

		for {
			tcplink := newconn.LinkToggleCh()
			select {
			case state := <-tcplink:
				log.Debug("TCPSrv listener got State %v", state)
				if state == link.LinkStateDown {
					return
				}
			case <-t.cx.DoneChan():
				log.Debug("TCPSrv listener close ch")
				return
			}
		}
	}
}
