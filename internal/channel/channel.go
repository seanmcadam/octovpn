package channel

import (
	"time"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/internal/pinger"
	"github.com/seanmcadam/octovpn/internal/tracker"
	"github.com/seanmcadam/octovpn/octolib/ctx"
)

type ChannelStruct struct {
	cx      *ctx.Ctx
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
		width:   packet.PacketWidth32,
		channel: ci,
		pinger:  pinger.NewPinger32(ctx, 1, 2),
		counter: counter.NewCounter32(ctx),
		tracker: tracker.NewTracker(ctx, time.Second),
		recvch:  make(chan *packet.PacketStruct, 16),
	}

	go cs.goRecv()
	return cs, err
}

func NewChannel64(ctx *ctx.Ctx, ci interfaces.ChannelInterface) (cs *ChannelStruct, err error) {

	cs = &ChannelStruct{
		cx:      ctx,
		width:   packet.PacketWidth64,
		channel: ci,
		pinger:  pinger.NewPinger64(ctx, 1, 2),
		counter: counter.NewCounter64(ctx),
		tracker: tracker.NewTracker(ctx, time.Second),
		recvch:  make(chan *packet.PacketStruct, 16),
	}

	go cs.goRecv()
	return cs, err
}

func (cs *ChannelStruct) Width() packet.PacketWidth {
	return cs.width
}

func (cs *ChannelStruct) Reset() {
	cs.channel.Reset()
}

func (cs *ChannelStruct) Close() {
	cs.close()
}

func (cs *ChannelStruct) close() {
	cs.cx.Cancel()
}
