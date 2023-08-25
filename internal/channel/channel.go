package channel

import (
	"log"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/octolib/counter"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/packet/packetchan"
	"github.com/seanmcadam/octovpn/octolib/tracker"
)

type ChannelStruct struct {
	channel interfaces.ChannelInterface
	counter *counter.Counter64Struct
	tracker *tracker.TrackerStruct
	closech chan interface{}
}

func NewChannel(ci interfaces.ChannelInterface) (cs *ChannelStruct, err error) {

	closech := make(chan interface{})

	cs = &ChannelStruct{
		channel: ci,
		counter: counter.NewCounter64(),
		tracker: tracker.NewTracker(closech),
		closech: closech,
	}
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

	cs.tracker.Push(packet)
	return cs.channel.Send(packet.ToByte())
}

func (cs *ChannelStruct) Recv() (b []byte, err error) {
	if !cs.channel.Active() {
		return nil, errors.ErrNetChannelDown
	}

	recv, err := cs.channel.Recv()

	cp, err := packetchan.MakePacket(recv)
	cs.tracker.Ack(counter.Counter64(cp.GetCounter()))

	return cp.ToByte(), err
}

func (cs *ChannelStruct) Reset() {
	cs.channel.Reset()
}

func (cs *ChannelStruct) Close() {
	cs.channel.Close()
	cs.close()
}

func (cs *ChannelStruct) close() {
	close(cs.closech)
}
