package site

import (
	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (s *SiteStruct) RecvChan() <-chan *packet.PacketStruct {
	if s == nil {
		return nil
	}
	return s.recvch
}

func (s *SiteStruct) goRecv(channel interfaces.ChannelSiteInterface) {
	if s == nil {
		return
	}

	s.Cancel()

	for {
		select {
		case <-s.doneChan():
			return
		case p := <-channel.RecvChan():
			s.recv(p, channel)
		}
	}
}

func (s *SiteStruct) recv(p *packet.PacketStruct, channel interfaces.ChannelSiteInterface) {
	if s == nil {
		return
	}

	if p == nil {
		log.ErrorStack("recv() Nil Packet")
		return
	}

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
		s.tracker.Recv(copy)
		s.recvch <- p

	case packet.SIG_CHAN_32_ACK:
		fallthrough
	case packet.SIG_CHAN_64_ACK:
		s.tracker.Ack(p.Counter())

	case packet.SIG_CHAN_32_NAK:
		fallthrough
	case packet.SIG_CHAN_64_NAK:
		s.tracker.Nak(p.Counter())

	default:
		log.FatalfStack("Unhandled CHAN TYPE:%d", t)
	}
}
