package packetconn

import (
	"encoding/binary"
	"fmt"

	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/packet/packetchan"
)

const ConnOverhead int = 4
const sigStart int = 0    // +1
const typeStart int = 1   // +1
const lengthStart int = 2 // +2
const payloadStart int = ConnOverhead

type PacketSig uint8     // 1
type PacketType uint8    // 1
type PacketLength uint16 // 2

const ConnSigVal PacketSig = 0xAA

const (
	PACKET_TYPE_UDPAUTH PacketType = 0x41 // []byte  (payload conversion)
	PACKET_TYPE_UDP     PacketType = 0x42 // ChanPacket
	PACKET_TYPE_TCPAUTH PacketType = 0x84 // []byte
	PACKET_TYPE_TCP     PacketType = 0x88 // ChanPacket
	PACKET_TYPE_LOOP    PacketType = 0xA2 // []byte
	PACKET_TYPE_PING    PacketType = 0xE1 // []byte
	PACKET_TYPE_PONG    PacketType = 0xE2 // []byte
	PACKET_TYPE_ERROR   PacketType = 0xFF // []byte
)

type ConnPacket struct {
	pSig    PacketSig
	pType   PacketType
	pLength PacketLength
	payload interface{}
}

// NewPacket()
// Packets coming from the low level connections
func NewPacket(t PacketType, payload interface{}) (cp *ConnPacket, err error) {
	var plen int

	switch payload.(type) {
	case nil:
		plen = 0
	case []byte:
		plen = len(payload.([]byte))
	case *packetchan.ChanPacket:
		plen = payload.(*packetchan.ChanPacket).GetSize()
	default:
		log.Errorf("Bad Payload Type:%t", payload)
		return nil, errors.ErrConnPayloadType
	}

	cp = &ConnPacket{
		pSig:    ConnSigVal,
		pType:   t,
		pLength: PacketLength(plen),
		payload: payload,
	}

	return cp, nil
}

func (cp *ConnPacket) GetType() (t PacketType) {
	return cp.pType
}

func (cp *ConnPacket) GetSize() (l int) {
	return int(cp.pLength) + ConnOverhead
}

func (cp *ConnPacket) GetPayloadLength() (l PacketLength) {
	return cp.pLength
}

func (cp *ConnPacket) GetPayload()( payload interface{} ) {
	switch cp.payload.(type) {
	case nil:
		return nil
	case []byte:
		payload = make([]byte, len(cp.payload.([]byte)))
		copy(payload.([]byte), cp.payload.([]byte))
	case *packetchan.ChanPacket:
		payload = cp.payload.(*packetchan.ChanPacket).Copy()
	default:
		log.Fatalf("Unhandled Type:%t",cp.payload)
	}
	return payload
}

func (cp *ConnPacket) Copy() (copy *ConnPacket) {
	copy = &ConnPacket{
		pSig:    ConnSigVal,
		pType:   cp.pType,
		pLength: cp.pLength,
		payload: cp.GetPayload(),
	}
	return copy
}

func MakePacket(data []byte) (cp *ConnPacket, err error) {

	if len(data) < (ConnOverhead) {
		log.Debugf("Short Packet data:%d < %d", len(data), ConnOverhead)
		return nil, errors.ErrConnShortPacket
	}

	if PacketSig(data[sigStart]) != ConnSigVal {
		return nil, errors.ErrChanBadSig
	}

	cp = &ConnPacket{
		pSig:    ConnSigVal,
		pType:   PacketType(data[typeStart]),
		pLength: PacketLength(binary.LittleEndian.Uint16(data[lengthStart : lengthStart+2])),
	}

	var payloadlen PacketLength

	switch cp.pType {
	case PACKET_TYPE_UDP:
		fallthrough
	case PACKET_TYPE_TCP:

		ch, err := packetchan.MakePacket(data[ConnOverhead:])
		if err != nil {
			return nil, err
		}

		payloadlen = PacketLength(ch.GetSize())
		cp.payload = ch

	case PACKET_TYPE_UDPAUTH:
		fallthrough
	case PACKET_TYPE_TCPAUTH:
		fallthrough
	case PACKET_TYPE_PING:
		fallthrough
	case PACKET_TYPE_PONG:
		fallthrough
	case PACKET_TYPE_LOOP:
		payloadlen = PacketLength(len(data) - ConnOverhead)
		cp.payload = data[ConnOverhead:]
	default:
		log.Debugf("Bad Packet type:%d", cp.pType)
		return nil, errors.ErrConnBadPacket
	}

	if payloadlen != cp.pLength {
		log.Debugf("Bad Packet length:%d != %d", payloadlen, uint16(cp.pLength))
		return nil, errors.ErrConnPayloadLength
	}

	return cp, nil
}

func (p *ConnPacket) ToByte() (b []byte) {
	// Signature
	b = append(b, byte(ConnSigVal))
	// Type
	b = append(b, byte(p.pType))
	// Length
	len := make([]byte, 2)
	binary.LittleEndian.PutUint16(len, uint16(p.pLength))
	b = append(b, len...)
	// Payload
	switch p.payload.(type) {
	case nil:
	case []byte:
		b = append(b, p.payload.([]byte)...)
	case *packetchan.ChanPacket:
		b = append(b, p.payload.(*packetchan.ChanPacket).ToByte()...)
	default:
		log.Fatalf("Unhandled Type:%t", p.payload)
	}

	return b
}

func (p PacketType) String() string {
	switch p {
	case PACKET_TYPE_UDPAUTH:
		return "UDPAUTH"
	case PACKET_TYPE_UDP:
		return "UDP"
	case PACKET_TYPE_TCPAUTH:
		return "TCPAUTH"
	case PACKET_TYPE_TCP:
		return "TCP"
	case PACKET_TYPE_LOOP:
		return "LOOP"
	case PACKET_TYPE_PING:
		return "PING"
	case PACKET_TYPE_PONG:
		return "PONG"
	case PACKET_TYPE_ERROR:
		return "ERROR"
	default:
		return fmt.Sprintf("UNKNOWN TYPE:%x", p)

	}
}
