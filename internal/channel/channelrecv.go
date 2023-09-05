package channel

import (
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (cs *ChannelStruct) RecvChan() <-chan *packet.PacketStruct {
	return cs.recvch
}

func (cs *ChannelStruct) goRecv() {

	//defer

	for {
		select {
		case <-cs.cx.DoneChan():
			return
		case p := <-cs.channel.RecvChan():
			cs.recv(p)
		}
	}
}

func (cs *ChannelStruct) recv(p *packet.PacketStruct) {
	t := p.Sig()
	switch t {
	case packet.SIG_CHAN_32_PACKET:
		fallthrough
	case packet.SIG_CHAN_64_PACKET:

		copyack, err := p.CopyAck()
		log.FatalfStack("CopyAck() Err:%s", err)
		copy, err := p.Copy()
		log.FatalfStack("Copy() Err:%s", err)

		cs.channel.Send(copyack)
		cs.tracker.Recv(copy)
		cs.recvch <- p

	case packet.SIG_CHAN_32_ACK:
		fallthrough
	case packet.SIG_CHAN_64_ACK:

		cs.tracker.Ack(p.Counter())

	case packet.SIG_CHAN_32_NAK:
	case packet.SIG_CHAN_64_NAK:
		cs.tracker.Nak(p.Counter())

	default:
		log.Fatalf("Unhandled CHAN TYPE:%d", t)
	}
}
