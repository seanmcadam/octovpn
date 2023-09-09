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

	if cs.link.GetState() != link.LinkStateUP {
		return nil
	}

	return cs.recvch
}

func (cs *ChanconnStruct) goRecv() {
	for {
		select {
		case <-cs.link.LinkDownCh():
			return
		case <-cs.cx.DoneChan():
			return

		case <-cs.link.LinkUpCh():
			continue
		case <-cs.link.LinkLinkCh():
			continue

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

			log.Debugf("Conn Recv:%v", p)

			switch p.Sig() {
			case packet.SIG_CONN_32_PACKET:
			case packet.SIG_CONN_64_PACKET:
				if p.Packet() == nil {
					log.FatalfStack("nil Packet(): %v", p)
				}
				cs.recvch <- p.Packet()

			case packet.SIG_CONN_32_AUTH:
			case packet.SIG_CONN_64_AUTH:
				cs.auth.GetRecvCh() <- p.Auth()

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
