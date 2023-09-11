package site

import (
	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (cs *SiteStruct) RecvChan() <-chan *packet.PacketStruct {
	return cs.recvch
}

func (cs *SiteStruct) goRecv(channel interfaces.ChannelSiteInterface) {

	//defer

	for {
		select {
		case <-cs.cx.DoneChan():
			return
		case p := <-channel.RecvChan():
			cs.recv(p, channel)
		}
	}
}

func (cs *SiteStruct) recv(p *packet.PacketStruct, channel interfaces.ChannelSiteInterface) {

	p.DebugPacket("CHAN RECV")

	t := p.Sig()
	switch t {
	case packet.SIG_CHAN_32_PACKET:
		fallthrough
	case packet.SIG_CHAN_64_PACKET:

		copyack, err := p.CopyAck()
		if err != nil {
			log.FatalfStack("CopyAck() Err:%s", err)
		}
		copy, err := p.Copy()
		if err != nil {
			log.FatalfStack("Copy() Err:%s", err)
		}

		channel.Send(copyack)
		cs.tracker.Recv(copy)
		cs.recvch <- p

	case packet.SIG_CHAN_32_ACK:
		fallthrough
	case packet.SIG_CHAN_64_ACK:
		cs.tracker.Ack(p.Counter())

	case packet.SIG_CHAN_32_NAK:
		fallthrough
	case packet.SIG_CHAN_64_NAK:
		cs.tracker.Nak(p.Counter())

	default:
		log.Fatalf("Unhandled CHAN TYPE:%d", t)
	}
}
