package chanconn

import (
	"time"

	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (cs *ChanconnStruct) goStart() {
	if cs == nil {
		return
	}

	log.Debugf("starting %s", cs.name)

	linkupch := cs.conn.Link().LinkUpCh()

	log.Debugf("conn link status:%s", cs.conn.Link().GetState())

	if cs.conn.Link().IsUp() {
		go cs.goAuth()
		return
	}

	for {
		select {
		case <-cs.doneChan():
			return

		case <-linkupch:
			log.Debug("Channel Link Up")
			go cs.goAuth()
			return

		case <-time.After(5 * time.Second):
			log.Debugf("conn link status:%s", cs.conn.Link().GetState())
			log.Debug("Channel Link Up - Timeout...")
			return
		}
	}
}

func (cs *ChanconnStruct) goAuth() {
	if cs == nil {
		return
	}

	log.Debugf("starting %s", cs.name)

	cs.auth.Run()

	for {
		select {
		case <-cs.doneChan():
			return

		case <-cs.auth.Link().LinkUpCh():
			log.Debug("Channel Authenticated")
			go cs.goRecv()
			return

		case ap := <-cs.auth.GetSendCh():
			var p *packet.PacketStruct
			var err error
			if cs.width == 32 {
				p, err = packet.NewPacket(packet.SIG_CONN_32_AUTH, ap, cs.counter.Next())
			} else {
				p, err = packet.NewPacket(packet.SIG_CONN_64_AUTH, ap, cs.counter.Next())
			}
			if err != nil {
				return
			}
			cs.send(p)

		case p := <-cs.conn.RecvChan():
			if p == nil {
				log.Debug("Got nil packet")
				return
			}

			//log.Debugf("Conn Recv:%v", p)
			p.DebugPacket("Chanconn Recv")

			switch p.Sig() {
			case packet.SIG_CONN_32_AUTH:
			case packet.SIG_CONN_64_AUTH:
				cs.auth.GetRecvCh() <- p.Auth()

			default:
				log.Fatalf("Unhandled Packet Type:%d", p.Sig())
			}
		}
	}
}
