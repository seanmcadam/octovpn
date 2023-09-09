package interfaces

import "github.com/seanmcadam/octovpn/internal/packet"

type ChannelSiteInterface interface {
	Send(*packet.PacketStruct) error
	RecvChan() <-chan *packet.PacketStruct
	Reset() error
	MaxLocalMtu() packet.PacketSizeType
}
