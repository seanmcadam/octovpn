package channel

import (
	"time"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/internal/link"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/internal/pinger"
	"github.com/seanmcadam/octovpn/internal/tracker"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type ChannelStruct struct {
	cx      *ctx.Ctx
	name    string
	link    *link.LinkStateStruct
	width   packet.PacketWidth
	channel interfaces.ChannelInterface
	pinger  pinger.PingerStruct
	counter counter.CounterStruct
	tracker *tracker.TrackerStruct
	recvch  chan *packet.PacketStruct
}

func NewChannel32(ctx *ctx.Ctx, ci interfaces.ChannelInterface) (cs *ChannelStruct, err error) {

	cs = &ChannelStruct{
		cx:      ctx,
		name:    ci.Name(),
		link:    link.NewLinkState(ctx),
		width:   packet.PacketWidth32,
		channel: ci,
		pinger:  pinger.NewPinger32(ctx, 1, 2),
		counter: counter.NewCounter32(ctx),
		tracker: tracker.NewTracker(ctx, time.Second),
		recvch:  make(chan *packet.PacketStruct, 16),
	}

	cs.link.AddLinkStateCh(cs.channel.Link())
	go cs.goRecv()
	return cs, err
}

func NewChannel64(ctx *ctx.Ctx, ci interfaces.ChannelInterface) (cs *ChannelStruct, err error) {

	cs = &ChannelStruct{
		cx:      ctx,
		link:    link.NewLinkState(ctx),
		width:   packet.PacketWidth64,
		channel: ci,
		pinger:  pinger.NewPinger64(ctx, 1, 2),
		counter: counter.NewCounter64(ctx),
		tracker: tracker.NewTracker(ctx, time.Second),
		recvch:  make(chan *packet.PacketStruct, 16),
	}

	cs.link.AddLinkStateCh(cs.channel.Link())
	go cs.goRecv()
	return cs, err
}

func (cs *ChannelStruct) Width() packet.PacketWidth {
	if cs == nil {
		return 0
	}
	return cs.width
}

func (cs *ChannelStruct) Link() *link.LinkStateStruct {
	if cs == nil {
		return nil
	}
	return cs.link
}

func (cs *ChannelStruct) Reset() error {
	if cs == nil {
		return errors.ErrNetNilPointerMethod(log.Errf(""))
	}
	return cs.channel.Reset()
}

func (c *ChannelStruct) MaxLocalMtu() (size packet.PacketSizeType) {
	if c == nil {
		return 0
	}
	size = packet.PacketSigSize + packet.PacketSize16Size
	if c.width == packet.PacketWidth32 {
		size += packet.PacketCounter32Size
		size += packet.PacketPing32Size
		if c.width == packet.PacketWidth64 {
			size += packet.PacketCounter64Size
			size += packet.PacketPing64Size
		} else {
			log.FatalfStack("CannedStruct:%v", c)
		}

		size += c.channel.MaxLocalMtu()

	}
	return size
}
