package packetsite

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

type SitePacket struct {
	sSig         packet.PacketSig
	sType        packet.PacketType
	sPayloadSize packet.PacketPayloadSize
	sPayload     interface{}
}

// NewPacket()
// Packets coming from the low level connections
func NewPacket(t packet.PacketType, payload interface{}) (cp *SitePacket, err error) {
	var plen int

	switch payload.(type) {
	case nil:
		plen = 0
	case []byte:
		plen = len(payload.([]byte))
		//	case *packetchan.ChanPacket:
		//		plen = payload.(*packetchan.ChanPacket).GetSize()
	default:
		log.Errorf("Bad Payload Type:%t", payload)
		return nil, errors.ErrConnPayloadType
	}

	cp = &SitePacket{
		sSig:         packet.SITE_SIGV1,
		sType:        t,
		sPayloadSize: packet.PacketPayloadSize(plen),
		sPayload:     payload,
	}

	return cp, nil
}

func (cp *SitePacket) Type() (t packet.PacketType) {
	return cp.sType
}

func (cp *SitePacket) Size() packet.PacketSize {
	return packet.PacketSize(int(cp.sPayloadSize) + Overhead)
}

func (cp *SitePacket) PayloadSize() (l packet.PacketPayloadSize) {
	return cp.sPayloadSize
}

func (cp *SitePacket) Counter32() packet.PacketCounter32 {
	return 0
}

func (cp *SitePacket) Payload() (payload interface{}) {
	switch cp.sPayload.(type) {
	case nil:
		return nil

	case []byte:
		payload = make([]byte, len(cp.sPayload.([]byte)))
		copy(payload.([]byte), cp.sPayload.([]byte))

	case *packetchan.ChanPacket:
		payload = cp.sPayload.(*packetchan.ChanPacket).Copy()

	default:
		log.Fatalf("Unhandled Type:%t", cp.sPayload)
	}
	return payload
}

func (cp *SitePacket) Copy() (copy interfaces.PacketInterface) {
	copy = &SitePacket{
		sSig:         packet.SITE_SIGV1,
		sType:        cp.sType,
		sPayloadSize: cp.sPayloadSize,
		sPayload:     cp.Payload(),
	}
	return copy
}

func (cp *SitePacket) CopyAck() (copy interfaces.PacketInterface) {
	copy = &SitePacket{
		sSig:         packet.SITE_SIGV1,
		sType:        packet.SITE_TYPE_ACK,
		sPayloadSize: cp.sPayloadSize,
		sPayload:     cp.Payload(),
	}
	return copy
}

func (cp *SitePacket) CopyPong64() (copy interfaces.PacketInterface) {
	copy = &SitePacket{
		sSig:         packet.SITE_SIGV1,
		sType:        packet.SITE_TYPE_PONG64,
		sPayloadSize: cp.sPayloadSize,
		sPayload:     cp.Payload(),
	}
	return copy
}

func MakePacket(data []byte) (cp *SitePacket, err error) {

	if len(data) < (Overhead) {
		log.Debugf("Short Packet data:%d < %d", len(data), Overhead)
		return nil, errors.ErrConnShortPacket
	}

	if packet.PacketSig(data[sigStart]) != packet.SITE_SIGV1 {
		return nil, errors.ErrChanBadSig
	}

	cp = &SitePacket{
		sSig:         packet.SITE_SIGV1,
		sType:        packet.PacketType(data[typeStart]),
		sPayloadSize: packet.PacketPayloadSize(binary.LittleEndian.Uint16(data[payloadStart:Overhead])),
		sPayload:     data[Overhead:],
	}

	return cp, nil
}

func (s *SitePacket) ToByte() (b []byte) {
	// Signature
	b = append(b, byte(packet.SITE_SIGV1))
	// Type
	b = append(b, byte(s.sType))
	// Length
	len := make([]byte, 2)
	binary.LittleEndian.PutUint16(len, uint16(s.sPayloadSize))
	b = append(b, len...)
	// Payload
	switch s.sPayload.(type) {
	case nil:
	case []byte:
		b = append(b, s.sPayload.([]byte)...)
	case *packetchan.ChanPacket:
		b = append(b, s.sPayload.(*packetchan.ChanPacket).ToByte()...)
	default:
		log.Fatalf("Unhandled Type:%t", s.sPayload)
	}

	return b
}
