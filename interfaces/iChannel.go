package interfaces

import "github.com/seanmcadam/octovpn/octolib/packet/packetchan"

// Interface between the Channel and ChanConn Layers
// All ChanConn objects are ChannleInterfaces
type ChannelInterface interface {
	Active() bool
	Send(*packetchan.ChanPacket) error
	RecvChan() <-chan *packetchan.ChanPacket
	Reset() error
}
