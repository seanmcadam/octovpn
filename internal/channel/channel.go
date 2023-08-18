package channel

import (
	"github.com/seanmcadam/octovpn/interfaces"
)

type ChannelStruct struct {
	channel interfaces.ChannelInterface
}

func NewChannel(ci interfaces.ChannelInterface) (cs *ChannelStruct, err error) {
	return &ChannelStruct{channel: ci}, err
}

func (cs *ChannelStruct) Send(b []byte) error {
	return cs.channel.Send(b)
}

func (cs *ChannelStruct) Recv() (b []byte, err error) {
	return cs.channel.Recv()
}

func (cs *ChannelStruct) Reset() (err error) {
	return cs.channel.Reset()
}
