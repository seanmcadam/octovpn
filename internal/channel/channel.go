package channel

import (
	"log"
	"time"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/internal/packet/packetchan"
	"github.com/seanmcadam/octovpn/internal/tracker"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/errors"
)

type ChannelStruct struct {
	cx      *ctx.Ctx
	channel interfaces.ChannelInterface
	counter *counter.Counter64Struct
	tracker *tracker.TrackerStruct
	recvch  chan interfaces.PacketInterface
}

func NewChannel(ctx *ctx.Ctx, ci interfaces.ChannelInterface) (cs *ChannelStruct, err error) {

	cs = &ChannelStruct{
		cx:      ctx,
		channel: ci,
		counter: counter.NewCounter64(ctx),
		tracker: tracker.NewTracker(ctx, time.Second),
		recvch:  make(chan interfaces.PacketInterface, 16),
	}

	go cs.goRecv()
	return cs, err
}

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

func (cs *ChannelStruct) RecvChan() <-chan interfaces.PacketInterface {
	return cs.recvch
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
