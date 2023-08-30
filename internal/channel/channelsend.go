package channel

import (
	"log"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/internal/packet/packetchan"
	"github.com/seanmcadam/octovpn/octolib/errors"
)

func (cs *ChannelStruct) Send(sp interfaces.PacketInterface) error {
	if !cs.channel.Active() {
		return errors.ErrNetChannelDown
	}

	packet, err := packetchan.NewPacket(packet.CHAN_TYPE_PARENT, sp)
	if err != nil {
		log.Fatalf("NewPacket Err:%s", err)
	}

	cs.tracker.Send(packet)
	return cs.channel.Send(packet)
}
