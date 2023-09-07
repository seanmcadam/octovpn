package interfaces

import (
	"github.com/seanmcadam/octovpn/internal/link"
	"github.com/seanmcadam/octovpn/internal/packet"
)

// Interface between the Conn and ChanConn Layers
// All Conn objects are ConnInterfaces
type ConnInterface interface {
	Send(*packet.PacketStruct) error
	RecvChan() <-chan *packet.PacketStruct
	Reset() error
	StateToggleCh() <-chan link.LinkStateType
	GetState() link.LinkStateType
}
