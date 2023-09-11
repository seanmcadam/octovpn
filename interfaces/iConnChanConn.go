package interfaces

import (
	"github.com/seanmcadam/octovpn/internal/link"
	"github.com/seanmcadam/octovpn/internal/packet"
)

// Interface between the Conn and ChanConn Layers
// All Conn objects are ConnInterfaces
//
// Conn Packet implements this for ChanConn
//
type ConnInterface interface {
	Send(*packet.PacketStruct) error
	RecvChan() <-chan *packet.PacketStruct
	Reset() error
	Link() *link.LinkStateStruct
	GetLinkNoticeStateCh() link.LinkNoticeStateCh
	GetLinkStateCh() link.LinkNoticeStateCh
	GetUpCh() link.LinkNoticeStateCh
	GetDownCh() link.LinkNoticeStateCh
	GetLinkCh() link.LinkNoticeStateCh
	GetState() link.LinkStateType
}
