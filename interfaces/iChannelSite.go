package interfaces

import "github.com/seanmcadam/octovpn/internal/packet"

type SiteInterface interface {
	Active() bool
	Send(*packet.PacketStruct) error
	RecvChan() <-chan *packet.PacketStruct
	Reset() error
	MaxLocalMtu() packet.PacketSizeType
}
