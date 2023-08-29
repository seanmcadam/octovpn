package interfaces

import "github.com/seanmcadam/octovpn/internal/packet"

// Interface between the Channel and ChanConn Layers
// All ChanConn objects are ChannleInterfaces
type PacketInterface interface {
	Type() packet.PacketType
	Size() packet.PacketSize
	PayloadSize() packet.PacketPayloadSize
	Counter32() packet.PacketCounter32
	Payload() interface{}
	Copy() PacketInterface
	CopyAck() PacketInterface
	CopyPong64() PacketInterface
	ToByte() []byte
}
