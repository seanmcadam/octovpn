package interfaces

import "github.com/seanmcadam/octovpn/internal/packet"

// Interface between the Conn and ChanConn Layers
// All Conn objects are ConnInterfaces
type ConnInterface interface {
	Active() bool
	Send(*packet.PacketStruct) error
	RecvChan() <-chan *packet.PacketStruct
	Reset() error
	//Stats() TrackerData
}
