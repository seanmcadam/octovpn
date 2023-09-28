package tcpcli

import (
	"time"

	"github.com/seanmcadam/octovpn/internal/layers/network/tcp"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (t *TcpClientStruct) RecvChan() <-chan *packet.PacketStruct {

	if t == nil || t.tcpconn == nil {
		log.Debugf("TCP Cli Recv Nil")
		return nil
	}

	if t.link.IsDown() {
		log.Debugf("TCP Cli Recv state:%s", t.link.GetState())
		return nil
	}

	return t.tcpconn.RecvChan()
}

// -
// goTcpStart()
// Wait for Timeout to receieve start packet, and send one back
// Move to auth or exit
// -
func (t *TcpClientStruct) goTcpStart(tcp *tcp.TcpStruct) {
	tcp.Run()
	t.Link().Start()

	if p, err := packet.NewPacket(packet.SIG_CONN_START); err != nil {
		log.Fatalf("NewPacket() Err:%s", err)
	} else {
		log.Debug("TCPCli Send START")
		tcp.Send(p)
	}

	for {
		select {
		case p := <-tcp.RecvChan():
			if p.Sig().Start() {
				t.goTcpRun(tcp)
				return
			} else if p.Sig().Close() { // Could happen
				tcp.Cancel()
				return
			} else {
				tcp.Cancel()
				log.Errorf("TCPCli Recv out of sequence packet %v", p)
				return
			}

		case <-t.Link().LinkUpCh():
			log.Error("TCPCli TCP Up")

		case <-time.After(5 * time.Second):
			log.Error("TCPCli Startup timeout")
			tcp.Cancel()
			return
		}
	}
}

// -
//
// -
func (t *TcpClientStruct) goTcpRun(tcp *tcp.TcpStruct) {
	defer tcp.Cancel()
	t.Link().Connected()
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
