package packet

import (
	"fmt"
	"net"
)

//
// 0-5 [6]: Dest MAC
// 6-11 [6]: Src MAC
// 12-13 [2]: EtherType
// 14-1514 [46-1500] Payload  (min len 46, max 1500)
//
// 0-5 [6]: Dest MAC
// 6-11 [6]: Src MAC
// 12-15 [4]: 802.1Q 0x8100 0xnnnn
// 16-17 [2]: EtherType
// 14-1514 [46-1500]: Payload
//

type EthFrame []byte
type TaggingSize byte
type Ethertype uint16

const (
	NotTaggedSize    TaggingSize = 0
	TaggedSize       TaggingSize = 4
	DoubleTaggedSize TaggingSize = 8
)
const SourceSize = 6
const DestinationSize = 6
const SourceDestinationSize = SourceSize + DestinationSize
const EthertypeSize = 2

// EtherTypes  http://en.wikipedia.org/wiki/Ethertype
var ET_IPv4 Ethertype = 0x0800
var ET_ARP Ethertype = 0x0806
var ET_IPv6 Ethertype = 0x86DD

func (et Ethertype) String() string {
	switch et {
	case ET_IPv4:
		return "IPv4"
	case ET_IPv6:
		return "IPv6"
	case ET_ARP:
		return "ARP"
	default:
		return fmt.Sprintf("%04x", et)
	}
}

//
//
//
func (f EthFrame) Length() PacketPayloadLen {
	return PacketPayloadLen(len(f))
}

//
//
//
func (f EthFrame) ToByte() []byte {
	return []byte(f)
}

//
//
//
func (f EthFrame) Type() PacketType {
	return EthPacket
}

//
// ID()
// EthFrames have no ID, so return 0
func (f EthFrame) PacketID() PacketID {
	return 0
}

//
//
//
func (f EthFrame) Destination() net.HardwareAddr {
	return net.HardwareAddr(f[:6:6])
}

//
//
//
func (f EthFrame) Source() net.HardwareAddr {
	return net.HardwareAddr(f[6:12:12])
}

//
//
//
func (f EthFrame) Tagging() TaggingSize {
	if f[12] == 0x81 && f[13] == 0x00 {
		return TaggedSize
	} else if f[12] == 0x88 && f[13] == 0xa8 {
		return DoubleTaggedSize
	}
	return NotTaggedSize
}

//
//
//
func (f EthFrame) Tags() []byte {
	min := SourceDestinationSize
	tagSize := f.Tagging()
	max := min + int(tagSize)
	limit := max
	return f[min:max:limit]
}

func (f EthFrame) Ethertype() Ethertype {
	ethertypePos := SourceDestinationSize + f.Tagging()
	return Ethertype(uint16(f[ethertypePos])<<8 | uint16(f[ethertypePos+1]))
}

//
//
//
func (f EthFrame) Payload() []byte {
	tagSize := f.Tagging()
	min := SourceDestinationSize + int(tagSize) + EthertypeSize
	return f[min:]
}

//
func (f *EthFrame) ResizePayload(payloadSize int) {
	tagging := NotTaggedSize
	if len(*f) > SourceDestinationSize+EthertypeSize {
		tagging = f.Tagging()
	}
	f.resize(SourceDestinationSize + int(tagging) + EthertypeSize + payloadSize)
}

// Prepare prepares *f to be used, by filling in dst/src address, setting up
// proper tagging and ethertype, and resizing it to proper length.
//
// It is safe to call Prepare on a pointer to a nil Frame or invalid Frame.
//func (f *Frame) Prepare(dst net.HardwareAddr, src net.HardwareAddr, tagging Tagging, ethertype Ethertype, payloadSize int) {
//	f.resize(6 + 6 + int(tagging) + 2 + payloadSize)
//	copy((*f)[0:6:6], dst)
//	copy((*f)[6:12:12], src)
//	if tagging == Tagged {
//		(*f)[12] = 0x81
//		(*f)[13] = 0x00
//	} else if tagging == DoubleTagged {
//		(*f)[12] = 0x88
//		(*f)[13] = 0xa8
//	}
//	(*f)[12+tagging] = ethertype[0]
//	(*f)[12+tagging+1] = ethertype[1]
//	return
//}

func (f *EthFrame) resize(length int) {
	if cap(*f) < length {
		old := *f
		*f = make(EthFrame, length, length)
		copy(*f, old)
	} else {
		*f = (*f)[:length]
	}
}
