package channel

import (
	"time"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/internal/pinger"
	"github.com/seanmcadam/octovpn/internal/tracker"
	"github.com/seanmcadam/octovpn/octolib/ctx"
)

type ChannelStruct struct {
	cx      *ctx.Ctx
	channel interfaces.ChannelInterface
	pinger  *pinger.Pinger32Struct
	tracker *tracker.TrackerStruct
	counter *counter.Counter32Struct
	recvch  chan interfaces.PacketInterface
}

func NewChannel(ctx *ctx.Ctx, ci interfaces.ChannelInterface) (cs *ChannelStruct, err error) {

	cs = &ChannelStruct{
		cx:      ctx,
		channel: ci,
		pinger:  pinger.NewPinger32(ctx, 1, 2),
		counter: counter.NewCounter32(ctx),
		tracker: tracker.NewTracker(ctx, time.Second),
		recvch:  make(chan interfaces.PacketInterface, 16),
	}

	go cs.goRecv()
	return cs, err
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
