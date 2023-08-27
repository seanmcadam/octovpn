package interfaces

import "github.com/seanmcadam/octovpn/octolib/packet/packetchan"

// Interface between the Channel and ChanConn Layers
// All ChanConn objects are ChannleInterfaces
type PacketInterface interface {
	Active() bool
	Send([]byte) error
	RecvChan() <-chan *packetchan.ChanPacket
	Reset() error
}
