package channel

import (
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (cs *ChannelStruct) Send(sp *packet.PacketStruct) (err error) {

	if cs == nil {
		return errors.ErrNetNilMethodPointer(log.Errf(""))
	}

	var sig packet.PacketSigType
	var p *packet.PacketStruct

	if sp.Sig().RouterLayer() {
		if cs.Width() == 32 {
			sig = packet.SIG_CHAN_32_PACKET
		} else {
			sig = packet.SIG_CHAN_64_PACKET
		}

		p, err = packet.NewPacket(sig, sp, <-cs.counter.GetCountCh())
		if err != nil {
			return errors.ErrChanSend(log.Errf("NewPacket Err:%s", err))
		}
	} else {
		p = sp
	}

	cs.tracker.Send(p)
	if err = cs.channel.Send(p); err != nil {
		return errors.ErrChanSend(log.Errf("Send() Err:%s", err))
	}

	return cs.channel.Send(p)
}
