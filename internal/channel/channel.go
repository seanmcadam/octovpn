package channel

import (
	"log"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/octolib/counter"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/packet/packetchan"
	"github.com/seanmcadam/octovpn/octolib/tracker"
)

type ChannelStruct struct {
	cx      *ctx.Ctx
	channel interfaces.ChannelInterface
	counter *counter.Counter64Struct
	tracker *tracker.TrackerStruct
	recvch  chan packetchan.ChanPayload
}

func NewChannel(ctx *ctx.Ctx, ci interfaces.ChannelInterface) (cs *ChannelStruct, err error) {

	cs = &ChannelStruct{
		cx:      ctx,
		channel: ci,
		counter: counter.NewCounter64(),
		tracker: tracker.NewTracker(),
		recvch:  make(packetchan.ChanPayload, 16),
	}

	go cs.goRecv()
	return cs, err
}

func (cs *ChannelStruct) Send(b []byte) error {
	if !cs.channel.Active() {
		return errors.ErrNetChannelDown
	}

	packet, err := packetchan.NewPacket(packetchan.CHAN_TYPE_DATA, b)
	if err != nil {
		log.Fatalf("NewPacket Err:%s", err)
	}

	cs.tracker.Send(packet)
	return cs.channel.Send(packet.ToByte())
}

func (cs *ChannelStruct) RecvChan() <-chan packetchan.ChanPayload {
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
