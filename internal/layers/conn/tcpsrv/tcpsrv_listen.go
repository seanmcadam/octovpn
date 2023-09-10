package tcpsrv

import (
	"github.com/seanmcadam/octovpn/internal/layers/network/tcp"
	"github.com/seanmcadam/octovpn/internal/link"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (t *TcpServerStruct) goListen() {
	if t == nil {
		return
	}

	defer t.cx.Done()

	for {
		conn, err := t.tcplistener.AcceptTCP()
		if err != nil {
			log.Debugf("AcceptTCP Error:%s", err)
			return
		}

		log.Debug("TCP New connection")
		newconn := tcp.NewTCP(t.cx.NewWithCancel(), conn)
		if newconn == nil {
			log.Debugf("NewTCP is Nil")
			return
		}
		t.tcpconnch <- newconn

		for {
			tcplink := newconn.Link().LinkStateCh()
			select {
			case state := <-tcplink:
				log.Debug("TCPSrv listener got State %v", state)
				if state.State() == link.LinkStateDOWN {
					return
				}
			case <-t.cx.DoneChan():
				log.Debug("TCPSrv listener close ch")
				return
			}
		}
	}
}
