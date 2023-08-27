package packetconn

import (
	"encoding/binary"
	"fmt"

	"github.com/seanmcadam/octovpn/octolib/counter"
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
	PACKET_TYPE_RAW    PacketType = 0x00 // []byte
	PACKET_TYPE_AUTH   PacketType = 0x01 // []byte  (payload conversion)
	PACKET_TYPE_CHAN   PacketType = 0x02 // ChanPacket
	PACKET_TYPE_PING64 PacketType = 0xE1 // []byte
	PACKET_TYPE_PONG64 PacketType = 0xE2 // []byte
	PACKET_TYPE_ERROR  PacketType = 0xFF // []byte
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
	case counter.Counter64:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(payload.(counter.Counter64)))
		payload = b
		plen = 8
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

func (cp *ConnPacket) GetPayload() (payload interface{}) {
	switch cp.payload.(type) {
	case nil:
		return nil

	case []byte:
		payload = make([]byte, len(cp.payload.([]byte)))
		copy(payload.([]byte), cp.payload.([]byte))

	case *packetchan.ChanPacket:
		payload = cp.payload.(*packetchan.ChanPacket).Copy()

	default:
		log.Fatalf("Unhandled Type:%t", cp.payload)
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
	case PACKET_TYPE_RAW:
		cp.payload = data[ConnOverhead:]

	case PACKET_TYPE_CHAN:

		ch, err := packetchan.MakePacket(data[ConnOverhead:])
		if err != nil {
			return nil, err
		}

		payloadlen = PacketLength(ch.GetSize())
		cp.payload = ch

	case PACKET_TYPE_PING64:
		fallthrough
	case PACKET_TYPE_PONG64:
		if payloadlen != 8 {
			log.Fatalf("Bad PING-PONG payload len:%d", payloadlen)
		}
		cp.payload = counter.Counter64(binary.LittleEndian.Uint64(data[ConnOverhead:]))

	case PACKET_TYPE_AUTH:
		fallthrough
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
	case PACKET_TYPE_CHAN:
		return "CHAN"
	case PACKET_TYPE_AUTH:
		return "AUTH"
	case PACKET_TYPE_PING64:
		return "PING64"
	case PACKET_TYPE_PONG64:
		return "PONG64"
	case PACKET_TYPE_RAW:
		return "RAW"
	case PACKET_TYPE_ERROR:
		return "ERROR"
	default:
		return fmt.Sprintf("UNKNOWN TYPE:%s", p.String())

	}
}
