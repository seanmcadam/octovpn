package packet

import (
	"encoding/binary"

	"github.com/seanmcadam/octovpn/octolib/log"
)

// -
//
// -
func calculateChecksumUint32(data []byte) uint32 {
	var sum uint32 = 0
	for _, b := range data {
		sum += uint32(b)
	}

	return sum
}
func calculateChecksumByte(data []byte) (b []byte) {
	b = make([]byte,4)
	if binary.BigEndian.PutUint32(b, calculateChecksumUint32(data)); len(b) != 4 {
		log.FatalStack("bad value returned")
	}
	return b
}

// -
//
// -
func validChecksum(sum1, sum2 uint32) bool {
	return sum1 == sum2
}
