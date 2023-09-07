package chanconn

import (
	"github.com/seanmcadam/octovpn/internal/link"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (cs *ChanconnStruct) RecvChan() <-chan *packet.PacketStruct {
	if cs == nil {
		log.FatalStack("nil ChanconnStruct")
		return nil
	}
	if cs.recvch == nil {
		log.Error("Nil recvch pointer")
		return nil
	}
	return cs.recvch
}

func (cs *ChanconnStruct) goRecv() {
	//cs.link.ToggleState(cs.conn.GetState())
	for {
		select {
		case state := <-cs.conn.StateToggleCh():
			if state == link.LinkStateDown {
				cs.link.ToggleState(state)
				return
			}
			if state == link.LinkStateClose {
				return
			}

		case <-cs.cx.DoneChan():
			return

		case p := <-cs.conn.RecvChan():
			if p == nil {
				log.Debug("Got nil packet")
				return
			}

			switch p.Sig() {
			case packet.SIG_CONN_32_PACKET:
			case packet.SIG_CONN_64_PACKET:
				if p.Packet() == nil {
					log.FatalfStack("nil Packet(): %v", p)
				}
				cs.recvch <- p.Packet()

			case packet.SIG_CONN_32_AUTH:
			case packet.SIG_CONN_64_AUTH:
				log.Fatal("Unhandled Sig")

			case packet.SIG_CONN_32_PING:
			case packet.SIG_CONN_64_PING:
				pong, err := p.CopyPong()
				if err != nil {
					log.FatalfStack("CopyPong() Err:%s", err)
				}
				cs.send(pong)

			case packet.SIG_CONN_32_PONG:
			case packet.SIG_CONN_64_PONG:
				cs.pinger.RecvPong(p.Pong())

			case packet.SIG_CHAN_32_RAW:
			case packet.SIG_CHAN_64_RAW:
				log.Debug("Discarded Sig")

			default:
				log.Fatalf("Unhandled Packet Type:%d", p.Sig())
			}
		}
	}
}
