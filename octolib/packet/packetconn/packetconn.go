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

type ConnSig uint8     // 1
type ConnType uint8    // 1
type ConnLength uint16 // 2

const ConnSigVal ConnSig = 0xAA

const (
	CONN_TYPE_RAW    ConnType = 0x00 // []byte
	CONN_TYPE_AUTH   ConnType = 0x01 // []byte  (payload conversion)
	CONN_TYPE_CHAN   ConnType = 0x02 // ChanConn
	CONN_TYPE_PING64 ConnType = 0xE1 // []byte
	CONN_TYPE_PONG64 ConnType = 0xE2 // []byte
	CONN_TYPE_ERROR  ConnType = 0xFF // []byte
)

type ConnPacket struct {
	cSig    ConnSig
	cType   ConnType
	cLength ConnLength
	payload interface{}
}

// NewConn()
// Conns coming from the low level connections
func NewConn(t ConnType, payload interface{}) (cp *ConnPacket, err error) {
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
		cSig:    ConnSigVal,
		cType:   t,
		cLength: ConnLength(plen),
		payload: payload,
	}

	return cp, nil
}

func (cp *ConnPacket) GetType() (t ConnType) {
	return cp.cType
}

func (cp *ConnPacket) GetSize() (l int) {
	return int(cp.cLength) + ConnOverhead
}

func (cp *ConnPacket) GetPayloadLength() (l ConnLength) {
	return cp.cLength
}

func (cp *ConnPacket) GetPayload() (payload interface{}) {
	switch cp.payload.(type) {
	case nil:
		return nil

	case []byte:
		payload = make([]byte, len(cp.payload.([]byte)))
		copy(payload.([]byte), cp.payload.([]byte))

	case *packetchan.ChanConn:
		payload = cp.payload.(*packetchan.ChanConn).Copy()

	default:
		log.Fatalf("Unhandled Type:%t", cp.payload)
	}
	return payload
}

func (cp *ConnPacket) Copy() (copy *ConnPacket) {
	copy = &ConnPacket{
		cSig:    ConnSigVal,
		cType:   cp.cType,
		cLength: cp.cLength,
		payload: cp.GetPayload(),
	}
	return copy
}

func MakeConn(data []byte) (cp *ConnPacket, err error) {

	if len(data) < (ConnOverhead) {
		log.Debugf("Short Conn data:%d < %d", len(data), ConnOverhead)
		return nil, errors.ErrConnShortConn
	}

	if ConnSig(data[sigStart]) != ConnSigVal {
		return nil, errors.ErrChanBadSig
	}

	cp = &ConnPacket{
		cSig:    ConnSigVal,
		cType:   ConnType(data[typeStart]),
		cLength: ConnLength(binary.LittleEndian.Uint16(data[lengthStart : lengthStart+2])),
	}

	var payloadlen ConnLength

	switch cp.cType {
	case CONN_TYPE_RAW:
		cp.payload = data[ConnOverhead:]

	case CONN_TYPE_CHAN:

		ch, err := packetchan.MakeConn(data[ConnOverhead:])
		if err != nil {
			return nil, err
		}

		payloadlen = ConnLength(ch.GetSize())
		cp.payload = ch

	case CONN_TYPE_PING64:
		fallthrough
	case CONN_TYPE_PONG64:
		if payloadlen != 8 {
			log.Fatalf("Bad PING-PONG payload len:%d", payloadlen)
		}
		cp.payload = counter.Counter64(binary.LittleEndian.Uint64(data[ConnOverhead:]))

	case CONN_TYPE_AUTH:
		fallthrough
	default:
		log.Debugf("Bad Conn type:%d", cp.cType)
		return nil, errors.ErrConnBadConn
	}

	if payloadlen != cp.cLength {
		log.Debugf("Bad Conn length:%d != %d", payloadlen, uint16(cp.cLength))
		return nil, errors.ErrConnPayloadLength
	}

	return cp, nil
}

func (p *ConnPacket) ToByte() (b []byte) {
	// Signature
	b = append(b, byte(ConnSigVal))
	// Type
	b = append(b, byte(p.cType))
	// Length
	len := make([]byte, 2)
	binary.LittleEndian.PutUint16(len, uint16(p.cLength))
	b = append(b, len...)
	// Payload
	switch p.payload.(type) {
	case nil:
	case []byte:
		b = append(b, p.payload.([]byte)...)
	case *packetchan.ChanConn:
		b = append(b, p.payload.(*packetchan.ChanConn).ToByte()...)
	default:
		log.Fatalf("Unhandled Type:%t", p.payload)
	}

	return b
}

func (p ConnType) String() string {
	switch p {
	case CONN_TYPE_CHAN:
		return "CHAN"
	case CONN_TYPE_AUTH:
		return "AUTH"
	case CONN_TYPE_PING64:
		return "PING64"
	case CONN_TYPE_PONG64:
		return "PONG64"
	case CONN_TYPE_RAW:
		return "RAW"
	case CONN_TYPE_ERROR:
		return "ERROR"
	default:
		return fmt.Sprintf("UNKNOWN TYPE:%s", p.String())

	}
}
