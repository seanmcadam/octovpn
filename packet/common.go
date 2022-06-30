package packet

import (
	"encoding/binary"
	"fmt"
)

// PacketBlock Unique Value #
type PacketBlock uint32

// PacketType defined the payload type
type PacketType uint32

// PacketVersion future proof the system
type PacketVersion uint32

// PacketPayloadLen indicated the lenth of the payload section
type PacketPayloadLen uint32

// PacketPayloadLen indicated the lenth of the payload section
type PacketID uint64

const PacketHeaderSize = 16 // Bytes
const ValidHeaderV1 PacketBlock = 0xf1e2
const VersionV1 PacketVersion = 1

const NoPacket PacketType = 0
const EthPacket PacketType = 2
const CommPacket PacketType = 3
const PingPacket PacketType = 4
const PongPacket PacketType = 5

type PacketInterface interface {
	ToByte() []byte
	Length() PacketPayloadLen
	Type() PacketType
	PacketID() PacketID
}

type CommonHeader struct {
	packet0Block      PacketBlock
	packet1Type       PacketType
	packet2Version    PacketVersion
	packet3PayloadLen PacketPayloadLen
	payload           []byte
}

func PacketBlockBuf(buf []byte) PacketBlock {
	if PacketHeaderSize != len(buf) {
		panic("")
	}
	return PacketBlock(convertTouint32(buf[0:4]))
}

func PacketVersionBuf(buf []byte) PacketVersion {
	if PacketHeaderSize != len(buf) {
		panic("")
	}
	return PacketVersion(convertTouint32(buf[8:12]))
}

//
// NewHeaderV1()
// Creates a CommonHeader struct used to transmit data over the network
//
func NewHeaderV1Payload(payload interface{}) (ch *CommonHeader) {

	var packet PacketInterface
	var packettype PacketType

	switch payload := payload.(type) {
	case *EthFrame:
		packet = payload
		packettype = EthPacket

	case *Ping:
		packet = payload
		packettype = PingPacket

	case *Pong:
		packet = payload
		packettype = PongPacket

	default:
		panic(fmt.Sprintf("bad type:%t", payload))
	}

	ch = &CommonHeader{
		packet0Block:      ValidHeaderV1,
		packet1Type:       packettype,
		packet2Version:    VersionV1,
		packet3PayloadLen: packet.Length(),
		payload:           packet.ToByte(),
	}

	return ch
}

func NewHeaderRead(packetheader []byte) (ch *CommonHeader, size PacketPayloadLen) {
	if len(packetheader) != PacketHeaderSize {
		panic("")
	}

	header := convertTouint32(packetheader[0:4])
	packettype := convertTouint32(packetheader[4:8])
	version := convertTouint32(packetheader[8:12])
	size = PacketPayloadLen(convertTouint32(packetheader[12:16]))

	if header != uint32(ValidHeaderV1) {
		panic("")
	}
	if version != 1 {
		panic("")
	}

	ch = &CommonHeader{
		packet0Block:      ValidHeaderV1,
		packet1Type:       PacketType(packettype),
		packet2Version:    PacketVersion(1),
		packet3PayloadLen: 0,
		payload:           nil,
	}

	return ch, size
}

//
//
//
func (c *CommonHeader) AddPayload(payload []byte) {
	length := len(payload)
	c.packet3PayloadLen = PacketPayloadLen(length)
	c.payload = payload
}

//
//
//
func (ch *CommonHeader) ToByte() (buf []byte) {
	totallen := ch.packet3PayloadLen + PacketHeaderSize
	buf = make([]byte, 0, totallen)
	buf = append(buf, convertToByte(ch.packet0Block)...)
	buf = append(buf, convertToByte(ch.packet1Type)...)
	buf = append(buf, convertToByte(ch.packet2Version)...)
	buf = append(buf, convertToByte(ch.packet3PayloadLen)...)
	buf = append(buf, ch.payload...)

	if totallen != PacketPayloadLen(len(buf)) {
		panic(fmt.Sprintf("Totallen:%d BufLen:%d", totallen, len(buf)))
	}

	return buf
}

func convertTouint32(buf []byte) (val uint32) {
	return binary.BigEndian.Uint32(buf)
}

func convertTouint64(buf []byte) (val uint64) {
	return binary.BigEndian.Uint64(buf)
}

func convertToByte(v interface{}) (buf []byte) {

	switch v := v.(type) {
	case PacketBlock:
		buf = make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(v))
	case PacketType:
		buf = make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(v))
	case PacketVersion:
		buf = make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(v))
	case PacketPayloadLen:
		buf = make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(v))
	case PacketID:
		buf = make([]byte, 4)
		binary.BigEndian.PutUint64(buf, uint64(v))
	default:
		panic(fmt.Sprintf("type:%t", v))
	}

	return buf
}

//
//
//
func (c CommonHeader) Valid() bool {
	if c.packet0Block == ValidHeaderV1 {
		return true
	}
	return false
}

//
//
//
func (c CommonHeader) GetVersion() PacketVersion {
	return c.packet2Version
}

//
//
//
func (c CommonHeader) GetType() PacketType {
	return c.packet1Type
}

//
// Do we need this?
//
func (c CommonHeader) GetPayloadLen() PacketPayloadLen {
	return PacketPayloadLen(len(c.payload))
}

//
//
//
func (c CommonHeader) GetPayload() (payload interface{}) {
	switch c.GetType() {

	case NoPacket:
		return nil

	case EthPacket:
		var eth EthFrame
		eth = append([]byte(eth), c.payload...)
		payload = eth

	case CommPacket:
		var eth EthFrame
		id := ConnFrameID(convertTouint32(c.payload[0:4]))
		eth = append([]byte(eth), c.payload[4:]...)
		return &ConnFrame{
			ID:    id,
			Frame: &eth,
		}

	case PingPacket:
		id := PingID(convertTouint32(c.payload[0:4]))
		sendunixmicro := int64(convertTouint64(c.payload[4:12]))
		_ = sendunixmicro
		//send := time.UnixMicro(sendunixmicro)
		return &Ping{
			id: id,
			//Send: send,
		}

	case PongPacket:
		pingid := PingID(convertTouint32(c.payload[0:4]))
		sendunixmicro := int64(convertTouint64(c.payload[4:12]))
		recvunixmicro := int64(convertTouint64(c.payload[12:20]))
		_ = sendunixmicro
		_ = recvunixmicro
		//send := time.UnixMicro(sendunixmicro)
		//recv := time.UnixMicro(recvunixmicro)
		return &Pong{
			id: pingid,
			//Send: send,
			//Recv: recv,
		}

	default:
		panic("")
	}

	return PacketPayloadLen(len(c.payload))
}
