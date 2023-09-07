package interfaces

import "github.com/seanmcadam/octovpn/internal/packet"

// Interface between the Channel and ChanConn Layers
// All ChanConn objects are ChannleInterfaces
type ChannelInterface interface {
	Send(*packet.PacketStruct) error
	RecvChan() <-chan *packet.PacketStruct
	Reset() error
	MaxLocalMtu() packet.PacketSizeType
	//Stats() TrackerData
}
