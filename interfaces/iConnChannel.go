package interfaces

import (
	"github.com/seanmcadam/octovpn/internal/link"
	"github.com/seanmcadam/octovpn/internal/packet"
)

// Interface between the Channel and ChanConn Layers
// All ChanConn objects are ChannleInterfaces
//
// Conn Packet implements this for Channel
type ChannelInterface interface {
	Name() string
	Send(*packet.PacketStruct) error
	RecvChan() <-chan *packet.PacketStruct
	Link() *link.LinkStateStruct
	Reset() error
	MaxLocalMtu() packet.PacketSizeType
	//Stats() TrackerData
}
