package channel

import (
	"log"

	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/errors"
)

func (cs *ChannelStruct) Send(sp *packet.PacketStruct) error {
	var sig packet.PacketSigType
	if !cs.channel.Active() {
		return errors.ErrNetChannelDown
	}

	if cs.Width() == 32 {
		sig = packet.SIG_CHAN_32_PACKET
	} else {
		sig = packet.SIG_CHAN_32_PACKET
	}
	packet, err := packet.NewPacket(sig,sp)
	if err != nil {
		log.Fatalf("NewPacket Err:%s", err)
	}

	// cs.tracker.Send(packet)
	return cs.channel.Send(packet)
}
