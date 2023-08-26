package packetchan

import (
	"encoding/binary"
	"fmt"

	"github.com/seanmcadam/octovpn/octolib/counter"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

const ChanOverhead int = 12
const sigStart int = 0    // +1
const typeStart int = 1   // +1
const lengthStart int = 2 // +2
const countStart int = 4  // +8
const payloadStart int = ChanOverhead

type ChanSig uint8               // 1
type ChanType uint8              // 1
type ChanLength uint16           // 2
type ChanCount counter.Counter64 // 8

const ChanSigVal ChanSig = 0xEE

const (
	CHAN_TYPE_DATA  ChanType = 0x11 // TBD
	CHAN_TYPE_ACK   ChanType = 0x21 // []byte
	CHAN_TYPE_NAK   ChanType = 0x22 // []byte
	CHAN_TYPE_PING  ChanType = 0xE1 // []byte
	CHAN_TYPE_PONG  ChanType = 0xE2 // []byte
	CHAN_TYPE_ERROR ChanType = 0xFF // []byte
)

type ChanPacket struct {
	cSig     ChanSig
	cType    ChanType
	cLength  ChanLength
	cCounter ChanCount
	cPayload interface{}
}

// NewPacket()
func NewPacket(t ChanType, payload interface{}) (cp *ChanPacket, err error) {
	var plen int

	switch payload.(type) {
	case nil:
		plen = 0
	case []byte:
		plen = len(payload.([]byte))
	default:
		log.Errorf("Bad Payload Type:%t", payload)
		return nil, errors.ErrChanPayloadType
	}

	cp = &ChanPacket{
		cSig:     ChanSigVal,
		cType:    t,
		cLength:  ChanLength(plen),
		cCounter: 0,
		cPayload: payload,
	}

	return cp, err
}

func (cp *ChanPacket) GetType() (t ChanType) {
	return cp.cType
}

func (cp *ChanPacket) GetSize() (l int) {
	return int(cp.cLength) + ChanOverhead
}

func (cp *ChanPacket) GetPayloadLength() (l ChanLength) {
	return cp.cLength
}

func (cp *ChanPacket) GetCounter() (l ChanCount) {
	return cp.cCounter
}

func (cp *ChanPacket) GetPayload() interface{} {
	log.Debug("TODO Need to make this a copy of the data")
	return cp.cPayload
}

func (cp *ChanPacket) Copy() (copy *ChanPacket) {
	copy = &ChanPacket{
		cSig:     ChanSigVal,
		cType:    cp.cType,
		cLength:  cp.cLength,
		cCounter: cp.cCounter,
		cPayload: cp.GetPayload(),
	}
	return copy
}
func (cp *ChanPacket) CopyDataToAck() (ack *ChanPacket) {
	ack = &ChanPacket{
		cSig:     ChanSigVal,
		cType:    CHAN_TYPE_ACK,
		cLength:  0,
		cCounter: cp.GetCounter(),
		cPayload: nil,
	}
	return ack
}

func MakePacket(data []byte) (cp *ChanPacket, err error) {

	if len(data) < (ChanOverhead) {
		log.Debugf("Short Packet data:%d < %d", len(data), ChanOverhead)
		return nil, errors.ErrChanShortPacket
	}

	if ChanSig(data[sigStart]) != ChanSigVal {
		return nil, errors.ErrChanBadSig
	}

	cp = &ChanPacket{
		cSig:     ChanSigVal,
		cType:    ChanType(data[typeStart]),
		cLength:  ChanLength(binary.LittleEndian.Uint16(data[lengthStart : lengthStart+2])),
		cCounter: ChanCount(binary.LittleEndian.Uint64(data[countStart : countStart+8])),
	}

	var payloadlen ChanLength

	switch cp.cType {
	case CHAN_TYPE_DATA:
		fallthrough
	case CHAN_TYPE_ACK:
		fallthrough
	case CHAN_TYPE_NAK:
		fallthrough
	case CHAN_TYPE_PING:
		fallthrough
	case CHAN_TYPE_PONG:

		payloadlen = ChanLength(len(data) - ChanOverhead)
		cp.cPayload = data[ChanOverhead:]
	default:
		return nil, errors.ErrChanBadPacket
	}

	if cp.cLength != payloadlen {
		log.Debugf("Bad Packet length:%d != %d", payloadlen, uint16(cp.cLength))
		return nil, errors.ErrChanPayloadLength

	}

	return cp, nil
}

func (cp *ChanPacket) ToByte() (b []byte) {
	// Signature
	b = append(b, byte(ChanSigVal))
	// Type
	b = append(b, byte(cp.cType))
	// Length
	len := make([]byte, 2)
	binary.LittleEndian.PutUint16(len, uint16(cp.cLength))
	b = append(b, len...)
	// Count
	count := make([]byte, 8)
	binary.LittleEndian.PutUint64(count, uint64(cp.cCounter))
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

func (p ChanType) String() string {
	switch p {
	case CHAN_TYPE_DATA:
		return "DATA"
	case CHAN_TYPE_ACK:
		return "ACK"
	case CHAN_TYPE_NAK:
		return "NAK"
	case CHAN_TYPE_PING:
		return "PING"
	case CHAN_TYPE_PONG:
		return "PONG"
	case CHAN_TYPE_ERROR:
		return "ERROR"
	default:
		return fmt.Sprintf("UNKNOWN CHAN TYPE:%x", p)
	}
}
