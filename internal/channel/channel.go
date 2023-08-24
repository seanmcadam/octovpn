package channel

import (
	"github.com/seanmcadam/octovpn/constants"
	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/octolib"
)

type ChannelStruct struct {
	channel interfaces.ChannelInterface
	counter octolib.Counter64
}

func NewChannel(ci interfaces.ChannelInterface) (cs *ChannelStruct, err error) {
	cs = &ChannelStruct{
		channel: ci,
		counter: *octolib.NewCounter64(),
	}
	return cs	, err
}

func (cs *ChannelStruct) Send(b []byte) (error) {
	if !cs.channel.Active(){
		return constants.ErrorChannelDown
	}
	return cs.channel.Send(b)
}

func (cs *ChannelStruct) Recv() (b []byte, err error) {
	if !cs.channel.Active(){
		return nil, constants.ErrorChannelDown
	}
	return cs.channel.Recv()
}

func (cs *ChannelStruct) Reset() {
	cs.channel.Reset()
}

func (cs *ChannelStruct) Close() {
	//cs.counter.Close()
	cs.channel.Close()
}
