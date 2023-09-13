package channel

import (
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (cs *ChannelStruct) RecvChan() <-chan *packet.PacketStruct {
	if cs == nil {
		return nil
	}
	return cs.recvch
}

func (cs *ChannelStruct) goRecv() {
	if cs == nil {
		return
	}

	defer cs.Cancel()

	for {
		select {
		case <-cs.doneChan():
			return
		case p := <-cs.channel.RecvChan():
			if err := cs.recv(p); err != nil {
				log.Error("%s", err)
				return
			}
		}
	}
}

func (cs *ChannelStruct) recv(p *packet.PacketStruct) (err error) {
	if cs == nil || p == nil {
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
			return errors.ErrChanRecv(log.Errf("CopyAck() Err:%s", err))
		}
		copy, err := p.Copy()
		if err != nil {
			return errors.ErrChanRecv(log.Errf("Copy() Err:%s", err))
		}

		cs.channel.Send(copyack)
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
		return errors.ErrChanRecv(log.Errf("Default Reached CHAN TYPE:%d", t))
	}
	return nil
}
