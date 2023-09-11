package interfaces

import (
	"github.com/seanmcadam/octovpn/internal/link"
	"github.com/seanmcadam/octovpn/internal/packet"
)

//
// Site Packet implements this for Router
//
type SiteRouterInterface interface {
	Send(*packet.PacketStruct) error
	RecvChan() <-chan *packet.PacketStruct
	Link() *link.LinkStateStruct
	Reset() error
	MaxLocalMtu() packet.PacketSizeType
}
