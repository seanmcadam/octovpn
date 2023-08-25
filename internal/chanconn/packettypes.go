package chanconn

import "encoding/binary"

const PacketOverhead int = 3

type PacketType uint8
type PacketLength uint16
type PacketPayload []byte

const (
	PACKET_TYPE_UDPAUTH PacketType = 0x41
	PACKET_TYPE_UDP     PacketType = 0x42
	PACKET_TYPE_TCPAUTH PacketType = 0x84
	PACKET_TYPE_TCP     PacketType = 0x88
	PACKET_TYPE_LOOP    PacketType = 0xF2
	PACKET_TYPE_ERROR   PacketType = 0xFF
)

type ConnPacket struct {
	pType   PacketType
	pLength PacketLength
	payload PacketPayload
}

func NewPacket(t PacketType, payload []byte) (cp *ConnPacket, err error) {

	if len(payload) == 0 {
		return nil, ErrChanConnPayloadLength
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

	if len(data) < 4 {
		return nil, ErrChanConnShortPacket
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
	case PACKET_TYPE_LOOP:
		t = PACKET_TYPE_LOOP
	default:
		return nil, ErrChanConnBadPacket
	}

	payloadlen := binary.LittleEndian.Uint16(data[1:3])

	if payloadlen != uint16(len(data)-3) {
		return nil, ErrChanConnPayloadLength
	}

	cp = &ConnPacket{
		pType:   t,
		pLength: PacketLength(payloadlen),
		payload: data[3:],
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
