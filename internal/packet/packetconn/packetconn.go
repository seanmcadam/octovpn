package packetconn

import (
	"encoding/binary"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/internal/packet/packetchan"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

const sigStart int = 0
const typeStart int = packet.SizePacketSig
const payloadSizeStart int = typeStart + packet.SizePacketType

// const Overhead int = 12
const Overhead int = payloadSizeStart + packet.SizePacketPayloadSize
const payloadStart int = Overhead

type ConnPacket struct {
	cSig         packet.PacketSig
	cType        packet.PacketType
	cPayloadSize packet.PacketPayloadSize
	cPayload     interface{}
}

// NewConn()
// Conns coming from the low level connections
func NewPacket(t packet.PacketType, payload interface{}) (cp *ConnPacket, err error) {
	var plen int

	log.Debug("TODO validate Type with Payload")

	switch payload.(type) {
	case nil:
		plen = 0
	case []byte:
		plen = len(payload.([]byte))
		//	case *packetchan.ChanPacket:
		//		plen = payload.(*packetchan.ChanPacket).PayloadSize()
	default:
		log.Errorf("Bad Payload Type:%t", payload)
		return nil, errors.ErrConnPayloadType
	}

	cp = &ConnPacket{
		cSig:         packet.CONN_SIGV1,
		cType:        t,
		cPayloadSize: packet.PacketPayloadSize(plen),
		cPayload:     payload,
	}

	return cp, nil
}

func (cp *ConnPacket) Type() (t packet.PacketType) {
	return cp.cType
}

func (cp *ConnPacket) Size() packet.PacketSize {
	return packet.PacketSize(int(cp.cPayloadSize) + Overhead)
}

func (cp *ConnPacket) PayloadSize() (l packet.PacketPayloadSize) {
	return cp.cPayloadSize
}

func (cp *ConnPacket) Counter32() packet.PacketCounter32 {
	return 0
}

func (cp *ConnPacket) Payload() (payload interface{}) {

	// Make copy of Payload

	switch cp.cPayload.(type) {
	case nil:
		return nil

	case []byte:
		payload = make([]byte, len(cp.cPayload.([]byte)))
		copy(payload.([]byte), cp.cPayload.([]byte))

	case *packetchan.ChanPacket:
		payload = cp.cPayload.(*packetchan.ChanPacket).Copy()

	default:
		log.Fatalf("Unhandled Type:%t", cp.cPayload)
	}
	return payload
}

func (cp *ConnPacket) Copy() (copy interfaces.PacketInterface) {
	copy = &ConnPacket{
		cSig:         packet.CONN_SIGV1,
		cType:        cp.cType,
		cPayloadSize: cp.cPayloadSize,
		cPayload:     cp.Payload(),
	}
	return copy
}

func (cp *ConnPacket) CopyPong64() (copy interfaces.PacketInterface) {
	copy = &ConnPacket{
		cSig:         packet.CONN_SIGV1,
		cType:        packet.CONN_TYPE_PONG64,
		cPayloadSize: cp.cPayloadSize,
		cPayload:     cp.Payload(),
	}
	return copy
}

func (cp *ConnPacket) CopyAck() (copy interfaces.PacketInterface) {
	copy = &ConnPacket{
		cSig:         packet.CONN_SIGV1,
		cType:        packet.CONN_TYPE_ACK,
		cPayloadSize: cp.cPayloadSize,
		cPayload:     cp.Payload(),
	}
	return copy
}

func MakePacket(data []byte) (cp *ConnPacket, err error) {

	if len(data) < (Overhead) {
		log.Debugf("Short Conn data:%d < %d", len(data), Overhead)
		return nil, errors.ErrConnShortPacket
	}

	if packet.PacketSig(data[sigStart]) != packet.CONN_SIGV1 {
		return nil, errors.ErrChanBadSig
	}

	cp = &ConnPacket{
		cSig:         packet.CONN_SIGV1,
		cType:        packet.PacketType(data[typeStart]),
		cPayloadSize: packet.PacketPayloadSize(binary.LittleEndian.Uint16(data[payloadSizeStart:payloadStart])),
		cPayload:     data[Overhead:],
	}

	return cp, nil
}

func (p *ConnPacket) ToByte() (b []byte) {
	// Signature
	b = append(b, byte(packet.CONN_SIGV1))
	// Type
	b = append(b, byte(p.cType))
	// Length
	len := make([]byte, 2)
	binary.LittleEndian.PutUint16(len, uint16(p.cPayloadSize))
	b = append(b, len...)
	// Payload
	switch p.cPayload.(type) {
	case nil:
	case []byte:
		b = append(b, p.cPayload.([]byte)...)
	case *packetchan.ChanPacket:
		b = append(b, p.cPayload.(*packetchan.ChanPacket).ToByte()...)
	default:
		log.Fatalf("Unhandled Type:%t", p.cPayload)
	}

	return b
}
