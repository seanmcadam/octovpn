package packet

import (
	"github.com/seanmcadam/octovpn/octolib/counter"
)

const SizePacketSig int = 1
const SizePacketType int = 1
const SizePacketSize int = 2
const SizePacketPayloadSize int = 2
const SizePacketCounter32 int = 4
const SizePacketCounter64 int = 8

type PacketSig uint8                   // 1
type PacketType uint8                  // 1
type PacketSize uint16                 // 2
type PacketPayloadSize uint16          // 2
type PacketCounter32 counter.Counter32 // 4
type PacketCounter64 counter.Counter64 // 8
