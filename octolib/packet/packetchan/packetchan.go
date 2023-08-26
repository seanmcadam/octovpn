package packetchan

import (
	"encoding/binary"

	"github.com/seanmcadam/octovpn/octolib/counter"
	"github.com/seanmcadam/octovpn/octolib/errors"
)

const ChanOverhead int = 12
const SigStart int = 0
const TypeStart int = 1
const LengthStart int = 2
const CountStart int = 4
const PayloadStart int = ChanOverhead

type ChanSig uint8               // 1
type ChanType uint8              // 1
type ChanLength uint16           // 2
type ChanCount counter.Counter64 // 8
type ChanPayload []byte          // n

const ChanSigVal ChanSig = 0xee

const (
	CHAN_TYPE_DATA  ChanType = 0x11
	CHAN_TYPE_ACK   ChanType = 0x21
	CHAN_TYPE_NAK   ChanType = 0x22
	CHAN_TYPE_PING  ChanType = 0xE1
	CHAN_TYPE_PONG  ChanType = 0xE2
	CHAN_TYPE_ERROR ChanType = 0xFF
)

type ChanPacket struct {
	cSig     ChanSig
	cType    ChanType
	cLength  ChanLength
	cCounter ChanCount
	cPayload ChanPayload
}

func NewPacket(t ChanType, payload []byte) (cp *ChanPacket, err error) {

	cp = &ChanPacket{
		cSig:     ChanSigVal,
		cType:    t,
		cLength:  ChanLength(len(payload)),
		cCounter: 0,
		cPayload: payload,
	}

	return cp, err
}

func (cp *ChanPacket) GetType() (t ChanType) {
	return cp.cType
}

func (cp *ChanPacket) GetLength() (l ChanLength) {
	return cp.cLength
}

func (cp *ChanPacket) GetCounter() (l ChanCount) {
	return cp.cCounter
}

func (cp *ChanPacket) GetPayload() (b []byte) {
	return cp.cPayload
}

func (cp *ChanPacket) Copy() (copy *ChanPacket) {
	copy = &ChanPacket{
		cSig:     ChanSigVal,
		cType:    cp.GetType(),
		cLength:  cp.GetLength(),
		cCounter: cp.GetCounter(),
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
		return nil, errors.ErrChanShortPacket
	}

	if ChanSig(data[SigStart]) != ChanSigVal {
		return nil, errors.ErrChanBadSig
	}

	var t ChanType = CHAN_TYPE_ERROR
	switch ChanType(data[TypeStart]) {
	case CHAN_TYPE_DATA:
		t = CHAN_TYPE_DATA
	case CHAN_TYPE_ACK:
		t = CHAN_TYPE_ACK
	case CHAN_TYPE_NAK:
		t = CHAN_TYPE_NAK
	case CHAN_TYPE_PING:
		t = CHAN_TYPE_PING
	case CHAN_TYPE_PONG:
		t = CHAN_TYPE_PONG
	default:
		return nil, errors.ErrChanBadPacket
	}

	payloadlen := binary.LittleEndian.Uint16(data[LengthStart:CountStart])

	if payloadlen != uint16(len(data)-ChanOverhead) {
		return nil, errors.ErrChanPayloadLength
	}

	counter := binary.LittleEndian.Uint64(data[CountStart:PayloadStart])

	cp = &ChanPacket{
		cSig:     ChanSigVal,
		cType:    t,
		cLength:  ChanLength(payloadlen),
		cCounter: ChanCount(counter),
		cPayload: nil,
	}

	if len(data) > (ChanOverhead) {
		cp.cPayload = data[PayloadStart:]
	}

	return cp, nil
}

func (cp *ChanPacket) ToByte() (b []byte) {
	b = append(b, byte(ChanSigVal))
	b = append(b, byte(cp.cType))
	len := make([]byte, 2)
	binary.LittleEndian.PutUint16(len, uint16(cp.cLength))
	b = append(b, len...)
	count := make([]byte, 8)
	binary.LittleEndian.PutUint64(count, uint64(cp.cCounter))
	b = append(b, count...)
	b = append(b, cp.cPayload...)
	return b
}
