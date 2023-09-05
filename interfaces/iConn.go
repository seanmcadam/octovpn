package interfaces

import "github.com/seanmcadam/octovpn/internal/packet"

// Interface between the Channel and ChanConn Layers
// All ChanConn objects are ChannleInterfaces
type ConnInterface interface {
	Active() bool
	Send(*packet.PacketStruct) error
	RecvChan() <-chan *packet.PacketStruct
	Reset() error
	//Stats() TrackerData
}
