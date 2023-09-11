package interfaces

import (
	"github.com/seanmcadam/octovpn/internal/link"
	"github.com/seanmcadam/octovpn/internal/packet"
)

//
// Channel Packet implements this for Site
//
type ChannelSiteInterface interface {
	Send(*packet.PacketStruct) error
	RecvChan() <-chan *packet.PacketStruct
	Link() *link.LinkStateStruct
	Reset() error
	MaxLocalMtu() packet.PacketSizeType
}
