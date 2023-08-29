package packetchan

import (
	"encoding/binary"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

const sigStart int = 0
const typeStart int = packet.SizePacketSig
const payloadLengthStart int = typeStart + packet.SizePacketType
const countStart int = payloadLengthStart + packet.SizePacketPayloadSize

// const Overhead int = 12
const Overhead int = countStart + packet.SizePacketCounter32
const payloadStart int = Overhead

type ChanPacket struct {
	cSig         packet.PacketSig
	cType        packet.PacketType
	cPayloadSize packet.PacketPayloadSize
	cCounter     packet.PacketCounter32
	cPayload     interface{}
}

// NewPacket()
func NewPacket(t packet.PacketType, payload interface{}) (pi interfaces.PacketInterface, err error) {
	var plen int

	//
	// TODO validate Type with Payload
	//

	switch payload.(type) {
	case nil:
		plen = 0
	case []byte:
		plen = len(payload.([]byte))
	case interfaces.PacketInterface:
		plen = int(payload.(interfaces.PacketInterface).PayloadSize())
	default:
		log.Errorf("Bad Payload Type:%t", payload)
		return nil, errors.ErrChanPayloadType
	}

	pi = &ChanPacket{
		cSig:         packet.CHAN_SIGV1,
		cType:        t,
		cPayloadSize: packet.PacketPayloadSize(plen),
		cCounter:     0,
		cPayload:     payload,
	}

	return pi, err
}

func (cp *ChanPacket) Type() (t packet.PacketType) {
	return cp.cType
}

func (cp *ChanPacket) Size() (l packet.PacketSize) {
	return packet.PacketSize(int(cp.cPayloadSize) + Overhead)
}

func (cp *ChanPacket) PayloadSize() (l packet.PacketPayloadSize) {
	return cp.cPayloadSize
}

func (cp *ChanPacket) Counter32() (l packet.PacketCounter32) {
	return cp.cCounter
}

func (cp *ChanPacket) Payload() interface{} {
	log.Debug("Need to copy payload here")
	return cp.cPayload
}

func (cp *ChanPacket) Copy() (copy interfaces.PacketInterface) {
	copy = &ChanPacket{
		cSig:         packet.CHAN_SIGV1,
		cType:        cp.cType,
		cPayloadSize: cp.cPayloadSize,
		cCounter:     cp.cCounter,
		cPayload:     cp.Payload(),
	}
	return copy
}

func (cp *ChanPacket) CopyAck() (ack interfaces.PacketInterface) {
	ack = &ChanPacket{
		cSig:         packet.CHAN_SIGV1,
		cType:        packet.CHAN_TYPE_ACK,
		cPayloadSize: 0,
		cCounter:     cp.Counter32(),
		cPayload:     nil,
	}
	return ack
}

func (cp *ChanPacket) CopyPong64() (ack interfaces.PacketInterface) {
	ack = &ChanPacket{
		cSig:         packet.CHAN_SIGV1,
		cType:        packet.CHAN_TYPE_PONG64,
		cPayloadSize: 0,
		cCounter:     cp.Counter32(),
		cPayload:     nil,
	}
	return ack
}

func MakePacket(data []byte) (cp *ChanPacket, err error) {

	if len(data) < (Overhead) {
		log.Debugf("Short Packet data:%d < %d", len(data), Overhead)
		return nil, errors.ErrChanShortPacket
	}

	if packet.PacketSig(data[sigStart]) != packet.CHAN_SIGV1 {
		return nil, errors.ErrChanBadSig
	}

	cp = &ChanPacket{
		cSig:         packet.CHAN_SIGV1,
		cType:        packet.PacketType(data[typeStart]),
		cPayloadSize: packet.PacketPayloadSize(binary.LittleEndian.Uint16(data[payloadLengthStart:countStart])),
		cCounter:     packet.PacketCounter32(binary.LittleEndian.Uint32(data[countStart:Overhead])),
		cPayload:     data[Overhead:],
	}

	return cp, nil
}

func (cp *ChanPacket) ToByte() (b []byte) {
	// Signature
	b = append(b, byte(packet.CHAN_SIGV1))
	// Type
	b = append(b, byte(cp.cType))
	// Length
	len := make([]byte, 2)
	binary.LittleEndian.PutUint16(len, uint16(cp.cPayloadSize))
	b = append(b, len...)
	// Count
	count := make([]byte, 8)
	binary.LittleEndian.PutUint32(count, uint32(cp.cCounter))
	b = append(b, count...)
	// Payload
	switch cp.cPayload.(type) {
	case nil:
	case []byte:
		b = append(b, cp.cPayload.([]byte)...)
	default:
		log.Fatalf("Unhandled Type:%t", cp.cPayload)
	}

	return b
}
