package interfaces

import (
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

// Interface between the Channel and ChanConn Layers
// All ChanConn objects are ChannleInterfaces
type ConnInterface interface {
	Active() bool
	Send(*packetconn.ConnPacket) error
	RecvChan() <-chan *packetconn.ConnPacket
	Reset() error
	//Stats() TrackerData
}
