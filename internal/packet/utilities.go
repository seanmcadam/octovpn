package packet

import (
	"encoding/binary"

	"github.com/seanmcadam/octovpn/octolib/log"
)

// -
// BtoU8()
// -
func BtoU8(b []byte) (u8 uint8) {
	if len(b) < 1 {
		log.FatalStack()
	}
	return (uint8(b[0]))
}

// -
// U8toB()
// -
func U8toB(u8 uint8) (b []byte) {
	b = make([]byte, 1)
	b[0] = byte(u8)
	return b
}

// -
// BtoU16()
// -
func BtoU16(b []byte) (u16 uint16) {
	if len(b) < 2 {
		log.FatalStack()
	}
	return uint16(binary.BigEndian.Uint16(b[:2]))
}

// -
// U16toB()
// -
func U16toB(u16 uint16) (b []byte) {
	b = make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(u16))
	return b
}

// -
// BtoU32()
// -
func BtoU32(b []byte) (u32 uint32) {
	if len(b) < 4 {
		log.FatalStack()
	}
	return uint32(binary.BigEndian.Uint32(b[:4]))
}

// -
// U32toB()
// -
func U32toB(u32 uint32) (b []byte) {
	b = make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(u32))
	return b
}

// -
// BtoU64()
// -
func BtoU64(b []byte) (u64 uint64) {
	if len(b) < 2 {
		log.FatalStack()
	}
	return uint64(binary.BigEndian.Uint64(b[:8]))
}

// -
// U64toB()
// -
func U64toB(u64 uint64) (b []byte) {
	b = make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(u64))
	return b
}
