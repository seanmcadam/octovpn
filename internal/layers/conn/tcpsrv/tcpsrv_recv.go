package tcpsrv

import (
	"time"

	"github.com/seanmcadam/octovpn/internal/layers/network/tcp"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (t *TcpServerStruct) RecvChan() <-chan *packet.PacketStruct {
	if t == nil {
		return nil
	}

	if len(t.tcpconn) == 0 {
		log.Debug("No Active TCP Srv connections")
		return nil
	}

	if t.link.IsDown() {
		log.Debugf("TCP Srv Recv state:%s", t.link.GetState())
		return nil
	}

	return t.recvch
}

// -
// goTcpStart()
// Wait for Timeout to receieve start packet, and send one back
// Move to auth or exit
// -
func (t *TcpServerStruct) goTcpStart(tcp *tcp.TcpStruct) {
	tcp.Run()
	t.Link().Start()

	log.Debug("TCPSrv START")

	for {
		select {
		case p := <-tcp.RecvChan():
			if p.Sig().Start() {
				if p, err := packet.NewPacket(packet.SIG_CONN_START); err != nil {
					t.removeConnection(tcp)
				} else {
					log.Debug("TCPSrv Send START")
					tcp.Send(p)
					t.goTcpRun(tcp)
				}
				return
			} else if p.Sig().Close() { // Could happen
				log.Debug("TCPSrv got Close()")
				t.removeConnection(tcp)
				return
			} else {
				log.Debug("TCPSrv got %s", p.Sig() )
				t.removeConnection(tcp)
				return
			}
		case <-time.After(1 * time.Second):
			log.Debug("TCPSrv Start Timeout")
			t.removeConnection(tcp)
			return
		}
	}
}

// -
//
// -
func (t *TcpServerStruct) goTcpRun(tcp *tcp.TcpStruct) {
	defer t.removeConnection(tcp)
	defer tcp.Cancel()
	t.Link().Connected()

	log.Debug("TCPSrv Running")

	for {
		select {
		case p := <-tcp.RecvChan():
			if p.Sig().Close() {
				return
			}

			t.recvch <- p

		case <-t.doneChan():
			return
		}
	}
}
