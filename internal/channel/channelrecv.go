package channel

import (
	"github.com/seanmcadam/octovpn/interfaces"
)

func (cs *ChannelStruct) RecvChan() <-chan interfaces.PacketInterface {
	return cs.recvch
}
