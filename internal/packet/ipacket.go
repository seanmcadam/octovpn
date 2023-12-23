package packet

// Interface between the Channel and ChanConn Layers
// All ChanConn objects are ChannleInterfaces
type PacketInterface interface {
	Type() PacketSigType
	Size() PacketSizeType
	ToByte() []byte
	Payload() PacketInterface
	Copy() PacketInterface
}
