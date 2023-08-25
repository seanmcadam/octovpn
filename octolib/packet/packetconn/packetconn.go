package packetconn

import (
	"encoding/binary"

	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

const PacketOverhead int = 3

type PacketType uint8
type PacketLength uint16
type PacketPayload []byte

const (
	PACKET_TYPE_UDPAUTH PacketType = 0x41
	PACKET_TYPE_UDP     PacketType = 0x42
	PACKET_TYPE_TCPAUTH PacketType = 0x84
	PACKET_TYPE_TCP     PacketType = 0x88
	PACKET_TYPE_LOOP    PacketType = 0xA2
	PACKET_TYPE_PING    PacketType = 0xE1
	PACKET_TYPE_PONG    PacketType = 0xE2
	PACKET_TYPE_ERROR   PacketType = 0xFF
)

type ConnPacket struct {
	pType   PacketType
	pLength PacketLength
	payload PacketPayload
}

func NewPacket(t PacketType, payload []byte) (cp *ConnPacket, err error) {

	if len(payload) == 0 {
		return nil, errors.ErrChanConnPayloadLength
	}
	cp = &ConnPacket{
		pType:   t,
		pLength: PacketLength(len(payload)),
		payload: payload,
	}

	return cp, nil
}

func (cp *ConnPacket) GetType() (t PacketType) {
	return cp.pType
}

func (cp *ConnPacket) GetLength() (l PacketLength) {
	return cp.pLength
}

func (cp *ConnPacket) GetPayload() (b []byte) {
	return cp.payload
}

func MakePacket(data []byte) (cp *ConnPacket, err error) {

	if len(data) < (PacketOverhead + 1) {
		log.Debugf("Short Packet data:%d < %d", len(data), PacketOverhead+1)
		return nil, errors.ErrChanConnShortPacket
	}

	var t PacketType = PACKET_TYPE_ERROR

	switch PacketType(data[0]) {
	case PACKET_TYPE_UDPAUTH:
		t = PACKET_TYPE_UDPAUTH
	case PACKET_TYPE_UDP:
		t = PACKET_TYPE_UDP
	case PACKET_TYPE_TCPAUTH:
		t = PACKET_TYPE_TCPAUTH
	case PACKET_TYPE_TCP:
		t = PACKET_TYPE_TCP
	case PACKET_TYPE_PING:
		t = PACKET_TYPE_PING
	case PACKET_TYPE_PONG:
		t = PACKET_TYPE_PONG
	case PACKET_TYPE_LOOP:
		t = PACKET_TYPE_LOOP
	default:
		log.Debugf("Bad Packet type:%d", PacketType(data[0]))
		return nil, errors.ErrChanConnBadPacket
	}

	payloadlen := binary.LittleEndian.Uint16(data[1:PacketOverhead])

	if payloadlen != uint16(len(data)-PacketOverhead) {
		log.Debugf("Bad Packet length:%d != %d", payloadlen, uint16(len(data)-PacketOverhead))
		return nil, errors.ErrChanConnPayloadLength
	}

	cp = &ConnPacket{
		pType:   t,
		pLength: PacketLength(payloadlen),
		payload: data[PacketOverhead:],
	}

	return cp, nil
}

func (p *ConnPacket) ToByte() (b []byte) {
	b = append(b, byte(p.pType))
	len := make([]byte, 2)
	binary.LittleEndian.PutUint16(len, uint16(p.pLength))
	b = append(b, len...)
	b = append(b, p.payload...)
	return b
}
