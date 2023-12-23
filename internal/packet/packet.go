package packet

import (
	"github.com/seanmcadam/counter"
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/internal/pinger"
	"github.com/seanmcadam/octovpn/octolib/errors"
)

const DefaultChannelDepth = 16

// Sig Type CounterType Overhead
// Type
//
//

type PacketStruct struct {
	pSig   PacketSigType
	pSize  PacketSizeType
	pWidth PacketWidth
	count  counter.Count
	ping   pinger.Ping
	pong   pinger.Pong
	router *RouterPacket
	ipv6   *IPv6Packet
	ipv4   *IPv4Packet
	eth    *EthPacket
	auth   *AuthPacket
	id     *IDPacket
	packet *PacketStruct
	raw    []byte
}

func NewPacket(sig PacketSigType, v ...interface{}) (ps *PacketStruct, err error) {
	var width PacketWidth

	// log.Debugf("Sig:%s", sig.String())

	if !sig.V1() {
		return nil, log.Errf("Bad Packet Version1: 0x%02X", uint16(sig&MASK_VERSION))
	}

	if sig.Width0() {
		width = PacketWidth0
	} else if sig.Width32() {
		width = PacketWidth32
	} else if sig.Width64() {
		width = PacketWidth64
	}

	ps = &PacketStruct{
		pSig:   sig,
		pSize:  PacketSigSize,
		pWidth: width,
		count:  nil,
		ping:   nil,
		pong:   nil,
		router: nil,
		ipv6:   nil,
		ipv4:   nil,
		eth:    nil,
		auth:   nil,
		id:     nil,
		packet: nil,
		raw:    nil,
	}

	if sig.Close() || sig.Start() {
		return ps, nil
	}

	// If it is a router or close there is no width needed
	if sig.Width0() != sig.RouterLayer() {
		err = log.Err("Width0 not with Router or Close")
		return nil, err
	}

	if sig.RouterLayer() && !sig.Data() {
		err = log.Errf("Router:%v Data:%v", sig.Router(), sig.Data())
		return nil, err
	}

	if sig.Size16() {
		ps.pSize += PacketSize16Size
	} else if sig.Size8() {
		ps.pSize += PacketSize8Size
	}

	for _, i := range v {
		switch u := i.(type) {
		case []byte:
			if !ps.pSig.Raw() {
				return nil, errors.ErrPacketBadParameter(log.Errf("Got RAW type for %s", ps.pSig))
			}
			if !ps.pSig.Data() {
				return nil, errors.ErrPacketBadParameter(log.Errf("Missing FIELD_DATA for %s", ps.pSig))
			}
			ps.pSize += PacketSizeType(len(u))
			ps.raw = u
		case *PacketStruct:
			if !ps.pSig.Packet() {
				return nil, errors.ErrPacketBadParameter(log.Errf("Got Parent type for %s", ps.pSig))
			}
			if !ps.pSig.Data() {
				return nil, errors.ErrPacketBadParameter(log.Errf("Missing FIELD_DATA for %s", ps.pSig))
			}
			ps.pSize += u.Size()
			ps.packet = u
		case *AuthPacket:
			if !ps.pSig.Auth() {
				return nil, errors.ErrPacketBadParameter(log.Errf("Got Auth type for %s", ps.pSig))
			}
			if !ps.pSig.Data() {
				return nil, errors.ErrPacketBadParameter(log.Errf("Missing FIELD_DATA for %s", ps.pSig))
			}
			ps.pSize += u.Size()
			ps.auth = u
		case *IDPacket:
			if !ps.pSig.ID() {
				return nil, errors.ErrPacketBadParameter(log.Errf("Got ID type for %s", ps.pSig))
			}
			if !ps.pSig.Data() {
				return nil, errors.ErrPacketBadParameter(log.Errf("Missing FIELD_DATA for %s", ps.pSig))
			}
			ps.pSize += u.Size()
			ps.id = u
		case *RouterPacket:
			if !ps.pSig.Router() {
				return nil, errors.ErrPacketBadParameter(log.Errf("Got Router type for %s", ps.pSig))
			}
			if !ps.pSig.Data() {
				return nil, errors.ErrPacketBadParameter(log.Errf("Missing FIELD_DATA for %s", ps.pSig))
			}
			ps.pSize += u.Size()
			ps.router = u
		case *IPv6Packet:
			if !ps.pSig.IPV6() {
				return nil, errors.ErrPacketBadParameter(log.Errf("Got IPv6 type for %s", ps.pSig))
			}
			if !ps.pSig.Data() {
				return nil, errors.ErrPacketBadParameter(log.Errf("Missing FIELD_DATA for %s", ps.pSig))
			}
			ps.pSize += u.Size()
			ps.ipv6 = u
		case *IPv4Packet:
			if !ps.pSig.IPV4() {
				return nil, errors.ErrPacketBadParameter(log.Errf("Got IPv4 type for %s", ps.pSig))
			}
			if !ps.pSig.Data() {
				return nil, errors.ErrPacketBadParameter(log.Errf("Missing FIELD_DATA for %s", ps.pSig))
			}
			ps.pSize += u.Size()
			ps.ipv4 = u
		case *EthPacket:
			if !ps.pSig.Eth() {
				return nil, errors.ErrPacketBadParameter(log.Errf("Got Eth type for %s", ps.pSig))
			}
			if !ps.pSig.Data() {
				return nil, errors.ErrPacketBadParameter(log.Errf("Missing FIELD_DATA for %s", ps.pSig))
			}
			ps.pSize += u.Size()
			ps.eth = u
		case counter.Count:
			if !ps.pSig.Count() {
				return nil, errors.ErrPacketBadParameter(log.Errf("Got Count type for %s", ps.pSig))
			}
			ps.count = u
			if ps.pSig.Width32() {
				ps.pSize += 4
			} else if ps.pSig.Width64() {
				ps.pSize += 8
			} else {
				log.FatalStack("No size")
			}
		case pinger.Ping:
			if !ps.pSig.Ping() {
				return nil, errors.ErrPacketBadParameter(log.Errf("Got Ping type for %s", ps.pSig))
			}
			ps.ping = u
			if ps.pSig.Width32() {
				ps.pSize += 4
			} else if ps.pSig.Width64() {
				ps.pSize += 8
			} else {
				log.FatalStack("No size")
			}
		case pinger.Pong:
			if !ps.pSig.Pong() {
				return nil, errors.ErrPacketBadParameter(log.Errf("Got Pong type for %s", ps.pSig))
			}
			ps.pong = u
			if ps.pSig.Width32() {
				ps.pSize += 4
			} else if ps.pSig.Width64() {
				ps.pSize += 8
			} else {
				log.FatalStack("No size")
			}
		default:
			log.FatalfStack("Default Reached Type:%t", v)
		}
	}

	if sig.Count() {
		if ps.count == nil {
			return nil, errors.ErrPacketNoCounterParameter(log.Errf("No Counter"))
		}
	} else if sig.Ack() {
		if ps.count == nil {
			return nil, errors.ErrPacketNoCounterParameter(log.Errf("No Counter"))
		}
	} else if sig.Nak() {
		if ps.count == nil {
			return nil, errors.ErrPacketNoCounterParameter(log.Errf("No Counter"))
		}
	}

	if sig.Ping() {
		if ps.count == nil {
			return nil, errors.ErrPacketNoPingParameter(log.Errf("No Ping"))
		}
	} else if sig.Pong() {
		if ps.count == nil {
			return nil, errors.ErrPacketNoPongParameter(log.Errf("No Pong"))
		}
	}

	if sig.Packet() {
		if ps.packet == nil {
			return nil, errors.ErrPacketNoPacketParameter(log.Errf(""))
		}
	} else if sig.IPV4() {
		if ps.ipv4 == nil {
			return nil, errors.ErrPacketNoIPv4Parameter(log.Errf(""))
		}
	} else if sig.IPV6() {
		if ps.ipv6 == nil {
			return nil, errors.ErrPacketNoIPv6Parameter(log.Errf(""))
		}
	} else if sig.Eth() {
		if ps.eth == nil {
			return nil, errors.ErrPacketNoEthParameter(log.Errf(""))
		}
	} else if sig.Router() {
		if ps.router == nil {
			return nil, errors.ErrPacketNoRouterParameter(log.Errf(""))
		}
	} else if sig.Raw() {
		if ps.raw == nil {
			return nil, errors.ErrPacketNoRawParameter(log.Errf(""))
		}
	}

	return ps, err
}

// MakePacket()
// Read raw byte in, parse, load data into appropriate data structures
// func MakePacket(raw []byte) (p *PacketStruct, err error) {
func ReadPacketBuffer(buf []byte) (sig PacketSigType, length PacketSizeType, err error) {

	// Not enough data to check if there is a full packet
	if len(buf) < 4 {
		log.ErrorStack("buf to short")
		return SIG_ERROR, 0, errors.ErrPacketBadParameter(log.Errf("Buffer is too short:%d", len(buf)))
	}

	sig = PacketSigType(BtoU32(buf[:4]))
	if !sig.V1() {
		return SIG_ERROR, 0, errors.ErrPacketBadParameter(log.Errf("Bad Sig Version:%s", sig))
	}
	buf = buf[4:]
	//
	// pSize
	//
	if sig.Close() || sig.Start() {
		length = 4
	} else if sig.Size8() {
		length = PacketSizeType(BtoU8(buf))
	} else if sig.Size16() {
		length = PacketSizeType(BtoU16(buf))
	} else {
		log.FatalStack("No size")
	}

	return sig, length, nil
}

// -
// MakePacket()
// Expects that raw is the correct size to create the packet
// Validates the checksum at the tail end of the packet
// -
func MakePacket(raw []byte) (p *PacketStruct, err error) {
	if raw == nil {
		return nil, errors.ErrNetNilMethodPointer(log.Errf(""))
	}

	var rawsize uint16 = uint16(len(raw))
	var calcsize PacketSizeType = 0

	if rawsize < uint16(PacketSigSize) {
		log.ErrorStack("size is too small, thats what she said...")
		return nil, errors.ErrPacketBadParameter(log.Errf("buffer is too small"))
	}

	p = &PacketStruct{}

	p.pSig = PacketSigType(BtoU32(raw[:4]))
	if !p.pSig.V1() {
		return nil, errors.ErrPacketBadVersion(log.Errf("Bad Sig Version:%s", p.pSig))
	}

	//
	// No checksum for Close() packets
	//
	if p.pSig.Close() || p.pSig.Start() {
		return p, nil
	}

	//
	// PacketSig  PacketSize(8) PacketChecksum
	//
	minsize := uint16(1) + uint16(PacketSigSize) + uint16(PacketSize8Size) + uint16(PacketChecksumSize)
	if rawsize < minsize {
		log.ErrorStack("size is too small, thats what she said...")
		return nil, errors.ErrPacketBadParameter(log.Errf("buffer is too small < minsize:%d", minsize))
	}

	checksumindex := len(raw) - int(PacketChecksumSize)
	//checksum_calc := calculateChecksumUint32(raw[:checksumindex])
	//checksum_actual := BtoU32(raw[checksumindex:])

	if !equalCRC2Checksums(calculateCRC32Checksum(raw[:checksumindex]), raw[checksumindex:]) {
		return nil, errors.ErrPacketBadChecksum(log.Errf("BadCheckSum"))
	}

	//
	// Remove Sig and Checksum from the raw packet
	//
	raw = raw[4 : len(raw)-int(PacketChecksumSize)]
	calcsize += PacketSigSize

	//
	// pSize (includes checksum at the end, so remove that part here)
	//
	if p.pSig.Size8() {
		p.pSize = PacketSizeType(BtoU8(raw)) - PacketChecksumSize
		raw = raw[1:]
		calcsize += 1

	} else if p.pSig.Size16() {
		p.pSize = PacketSizeType(BtoU16(raw)) - PacketChecksumSize
		raw = raw[2:]
		calcsize += 2
	} else {
		return nil, errors.ErrPacketNoSizeParameter(log.Errf(""))
	}

	// If the buffer does not contain enough data, Error out now
	if rawsize != (uint16(p.pSize) + uint16(PacketChecksumSize)) {
		return nil, log.Errf("Buffer size has %d, needs:%d", rawsize, p.pSize+PacketChecksumSize)
	}

	if p.pSig.Width32() {
		p.pWidth = PacketWidth32
	} else if p.pSig.Width64() {
		p.pWidth = PacketWidth64
	}

	//
	// Counters
	//
	if p.pSig.Count() || p.pSig.Ack() || p.pSig.Nak() {
		if p.pSig.Width32() {
			if p.count, err = counter.ByteToCount(raw[:4]); err != nil {
				return nil, err
			}
			raw = raw[4:]
			calcsize += 4
		} else if p.pSig.Width64() {
			if p.count, err = counter.ByteToCount(raw[:8]); err != nil {
				return nil, err
			}
			raw = raw[8:]
			calcsize += 8
		} else {
			return nil, errors.ErrPacketNoCounterParameter(log.Errf("Counter"))
		}
	}

	//
	// Ping/Pong
	//
	if p.pSig.Ping() {
		if p.pSig.Width32() {
			if p.ping, err = pinger.NewBytePing32(raw[:4]); err != nil {
				return nil, err
			}
			raw = raw[4:]
			calcsize += 4
		} else if p.pSig.Width64() {
			if p.ping, err = pinger.NewBytePing64(raw[:8]); err != nil {
				return nil, err
			}
			raw = raw[8:]
			calcsize += 8
		} else {
			return nil, errors.ErrPacketNoCounterParameter(log.Errf("Ping"))
		}
	}

	if p.pSig.Pong() {
		if p.pSig.Width32() {
			if p.ping, err = pinger.NewBytePong32(raw[:4]); err != nil {
				return nil, err
			}
			raw = raw[4:]
			calcsize += 4
		} else if p.pSig.Width64() {
			if p.ping, err = pinger.NewBytePong64(raw[:8]); err != nil {
				return nil, err
			}
			raw = raw[8:]
			calcsize += 8
		} else {
			return nil, errors.ErrPacketNoCounterParameter(log.Errf("Pong"))
		}
	}

	if p.pSig.Data() {
		if p.pSig.Packet() {
			if p.packet, err = MakePacket(raw); err != nil {
				return nil, err
			}
			calcsize += p.packet.Size()
		} else if p.pSig.IPV4() {
			if p.ipv4, err = MakeIPv4(raw); err != nil {
				return nil, err
			}
			calcsize += p.ipv4.Size()
		} else if p.pSig.IPV6() {
			if p.ipv6, err = MakeIPv6(raw); err != nil {
				return nil, err
			}
			calcsize += p.ipv6.Size()
		} else if p.pSig.Eth() {
			if p.eth, err = MakeEth(raw); err != nil {
				return nil, err
			}
			calcsize += p.eth.Size()
		} else if p.pSig.Router() {
			if p.router, err = MakeRouter(raw); err != nil {
				return nil, err
			}
			calcsize += p.router.Size()
		} else if p.pSig.ID() {
			if p.id, err = MakeID(raw); err != nil {
				return nil, err
			}
			calcsize += p.id.Size()
		} else if p.pSig.Auth() {
			if p.auth, err = MakeAuth(raw); err != nil {
				return nil, err
			}
			calcsize += p.auth.Size()
		} else if p.pSig.Raw() {
			p.raw = raw
			calcsize += PacketSizeType(len(raw))
		}

	}

	if calcsize != p.pSize || PacketSizeType(rawsize) != p.pSize+PacketChecksumSize {
		p.DebugPacket("PACKET SIZE ISSUE")
		return nil, AuthErrPacketSizeMismatch(log.Errf("Packet Sizes dont match psize:%d rawsize:%d calcsize:%d", p.pSize, rawsize, calcsize))
	}

	//checksum := calculateChecksumByte(raw)
	checksum := calculateCRC32Checksum(raw)

	raw = append(raw, checksum...)

	return p, err
}

// -
// ToByte()
//
// Adds the checksum on the end and adjusts the sending size to account for it (but not the original object - dont touch the original)
//
// -
func (p *PacketStruct) ToByte() (raw []byte, err error) {
	if p == nil {
		return nil, errors.ErrPacketNilMethodPointer(log.Errf(""))
	}

	//
	// pSig
	//
	raw = append(raw, U32toB(uint32(p.pSig))...)

	//
	// This is a Close packet (4 bytes)
	if p.pSig.Close() || p.pSig.Start() {
		return raw, nil
	}

	//
	// pSize (account for added checksum)
	//
	var packetsendsize PacketSizeType = p.pSize + PacketChecksumSize

	if p.pSig.Size8() {
		if packetsendsize > 255 {
			return nil, errors.ErrPacketBadParameter(log.Errf("size > 255"))
		}
		raw = append(raw, U8toB(uint8(packetsendsize))...)
	} else if p.pSig.Size16() {
		if packetsendsize > 65534 {
			return nil, errors.ErrPacketBadParameter(log.Errf("size > 65534"))
		}
		raw = append(raw, U16toB(uint16(packetsendsize))...)
	} else {
		return nil, errors.ErrPacketBadParameter(log.Errf("No size"))
	}

	//
	// Counter/Ack/Nak
	//
	if p.pSig.Count() || p.pSig.Ack() || p.pSig.Nak() {
		raw = append(raw, p.count.ToByte()...)
	}

	//
	// Ping/Pong
	//
	if p.pSig.Ping() {
		raw = append(raw, p.ping.ToByte()...)
	} else if p.pSig.Pong() {
		raw = append(raw, p.pong.ToByte()...)
	}

	if p.pSig.Packet() {
		if rawpacket, err := p.packet.ToByte(); err != nil {
			return nil, errors.ErrPacketBadParameter(log.Errf("Sig:%s Err:%s", p.pSig, err))
		} else {
			raw = append(raw, rawpacket...)
		}
	} else if p.pSig.IPV4() {
		raw = append(raw, p.ipv4.ToByte()...)
	} else if p.pSig.IPV6() {
		raw = append(raw, p.ipv6.ToByte()...)
	} else if p.pSig.Eth() {
		raw = append(raw, p.eth.ToByte()...)
	} else if p.pSig.Router() {
		raw = append(raw, p.router.ToByte()...)
	} else if p.pSig.Auth() {
		if rawpacket, err := p.auth.ToByte(); err != nil {
			return nil, errors.ErrPacketBadParameter(log.Errf("Sig:%s Err:%s", p.pSig, err))
		} else {
			raw = append(raw, rawpacket...)
		}
	} else if p.pSig.ID() {
		raw = append(raw, p.id.ToByte()...)
	} else if p.pSig.Raw() {
		raw = append(raw, p.raw...)
	}

	raw = append(raw, calculateCRC32Checksum(raw)...)
	return raw, nil
}

func (p *PacketStruct) Sig() PacketSigType {
	if p == nil {
		return SIG_ERROR
	}

	return p.pSig
}

func (p *PacketStruct) Size() PacketSizeType {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return PacketSizeTypeERROR
	}

	return p.pSize
}

func (p *PacketStruct) Width() PacketWidth {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return 0
	}

	return p.pWidth
}

func (p *PacketStruct) Count() counter.Count {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return nil
	}

	return p.count
}

func (p *PacketStruct) Ping() pinger.Ping {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return nil
	}

	return p.ping
}

func (p *PacketStruct) Pong() pinger.Pong {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return nil
	}

	return p.pong
}

func (p *PacketStruct) Router() *RouterPacket {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return nil
	}

	return p.router
}

func (p *PacketStruct) IPv4() *IPv4Packet {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return nil
	}

	return p.ipv4
}

func (p *PacketStruct) IPv6() *IPv6Packet {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return nil
	}

	return p.ipv6
}

func (p *PacketStruct) Eth() *EthPacket {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return nil
	}

	return p.eth
}

func (p *PacketStruct) Auth() *AuthPacket {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return nil
	}

	return p.auth
}

func (p *PacketStruct) ID() *IDPacket {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return nil
	}

	return p.id
}

func (p *PacketStruct) Packet() *PacketStruct {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return nil
	}

	return p.packet
}

func (p *PacketStruct) Raw() []byte {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return nil
	}

	return p.raw
}
