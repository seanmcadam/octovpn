package msgbus

import (
	"fmt"

	"github.com/seanmcadam/octovpn/internal/packet"
)

type PacketMsg uint8

// -
// Packet()
// -
func (mb *MsgBus) Packet(target MsgTarget, packet *packet.PacketStruct) {
	topic := fmt.Sprintf("%s:%s", string(target), string(MsgPacket))
	mb.eventbus.Publish(topic, packet)
}

// -
// Register Packet Handler
// -
func (mb *MsgBus) PacketHandler(source MsgTarget, handler func(...interface{})) (err error) {
	topic := fmt.Sprintf("%s:%s", string(source), string(MsgPacket))
	return mb.eventbus.Subscribe(topic, handler)
}
