package packet

import (
	"encoding/binary"
	"hash/crc32"
)

// -
//
// -
func calculateCRC32Checksum(data []byte) []byte {
	checksum := crc32.ChecksumIEEE(data)
	checksumBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(checksumBytes, checksum)
	return checksumBytes
}

// -
//
// -
func equalCRC2Checksums(sum1, sum2 []byte) bool {
	if int(PacketChecksumSize) != len(sum1) && len(sum1) != len(sum2) {
		return false
	}
	return string(sum1) == string(sum2)
}
