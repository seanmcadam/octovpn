package packet

import (
	"fmt"

	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/internal/pinger"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Sig Type CounterType Overhead
// Type
//
//

type PacketStruct struct {
	pSig    PacketSigType
	pSize   PacketSizeType
	pWidth  PacketWidth
	counter counter.Counter
	ping    pinger.Ping
	pong    pinger.Pong
	router  *RouterStruct
	ipv6    *IPv6Struct
	ipv4    *IPv4Struct
	eth     *EthStruct
	auth    *AuthStruct
	id      *IDStruct
	packet  *PacketStruct
	raw     []byte
}

func NewPacket(sig PacketSigType, v ...interface{}) (ps *PacketStruct, err error) {
	var width PacketWidth
	if !sig.V1() {
		return nil, fmt.Errorf("Bad Packet Version1: 0x%02X", uint16(sig&MASK_VERSION))
	}

	if sig.Width0() {
		width = PacketWidth0
	} else if sig.Width32() {
		width = PacketWidth32
	} else if sig.Width64() {
		width = PacketWidth64
	}

	ps = &PacketStruct{
		pSig:    sig,
		pSize:   PacketSigSize,
		pWidth:  width,
		counter: nil,
		ping:    nil,
		pong:    nil,
		router:  nil,
		ipv6:    nil,
		ipv4:    nil,
		eth:     nil,
		auth:    nil,
		id:      nil,
		packet:  nil,
		raw:     nil,
	}

	if sig.Router() != sig.Width0() {
		err = fmt.Errorf("Router:%v Width0:%v", sig.Router(), sig.Width0())
		return nil, err
	}

	if sig.Size16() {
		ps.pSize += PacketSize16Size
	} else if sig.Size8() {
		ps.pSize += PacketSize8Size
	} else {
		log.Error("NewPacket() no 8/16 size")
	}

	for _, i := range v {
		switch u := i.(type) {
		case []byte:
			if !ps.pSig.Raw() {
				log.Errorf("Got RAW type for %s", ps.pSig)
			}
			ps.pSize += PacketSizeType(len(u))
			ps.raw = u
		case *PacketStruct:
			if !ps.pSig.Packet() {
				log.Errorf("Got Parent type for %s", ps.pSig)
			}
			ps.pSize += u.Size()
			ps.packet = u
		case *AuthStruct:
			if !ps.pSig.Auth() {
				log.Errorf("Got Auth type for %s", ps.pSig)
			}
			ps.pSize += u.Size()
			ps.auth = u
		case *IDStruct:
			if !ps.pSig.ID() {
				log.Errorf("Got ID type for %s", ps.pSig)
			}
			ps.pSize += u.Size()
			ps.id = u
		case *RouterStruct:
			if !ps.pSig.Router() {
				log.Errorf("Got Router type for %s", ps.pSig)
			}
			ps.pSize += u.Size()
			ps.router = u
		case *IPv6Struct:
			if !ps.pSig.IPV6() {
				log.Errorf("Got IPv6 type for %s", ps.pSig)
			}
			ps.pSize += u.Size()
			ps.ipv6 = u
		case *IPv4Struct:
			if !ps.pSig.IPV4() {
				log.Errorf("Got IPv4 type for %s", ps.pSig)
			}
			ps.pSize += u.Size()
			ps.ipv4 = u
		case *EthStruct:
			if !ps.pSig.Eth() {
				log.Errorf("Got Eth type for %s", ps.pSig)
			}
			ps.pSize += u.Size()
			ps.eth = u
		case counter.Counter:
			if !ps.pSig.Counter() {
				log.Errorf("Got Counter type for %s", ps.pSig)
			}
			ps.counter = u
			if ps.pSig.Width32() {
				ps.pSize += 4
			} else if ps.pSig.Width64() {
				ps.pSize += 8
			} else {
				log.FatalStack("No size")
			}
		case pinger.Ping:
			if !ps.pSig.Ping() {
				log.Errorf("Got Ping type for %s", ps.pSig)
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
				log.Errorf("Got Pong type for %s", ps.pSig)
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

	if sig.Counter() {
		if ps.counter == nil {
			log.FatalStack("no counter")
		}
	} else if sig.Ack() {
		if ps.counter == nil {
			log.FatalStack("no ack - counter")
		}
	} else if sig.Nak() {
		if ps.counter == nil {
			log.FatalStack("no nak - counter")
		}
	}

	if sig.Ping() {
		if ps.counter == nil {
			log.FatalStack("no ping")
		}
	} else if sig.Pong() {
		if ps.counter == nil {
			log.FatalStack("no pong")
		}
	}

	if sig.Packet() {
		if ps.packet == nil {
			log.FatalStack("no packet")
		}
	} else if sig.IPV4() {
		if ps.ipv4 == nil {
			log.FatalStack("no ipv4")
		}
	} else if sig.IPV6() {
		if ps.ipv6 == nil {
			log.FatalStack("no ipv6")
		}
	} else if sig.Eth() {
		if ps.eth == nil {
			log.FatalStack("no eth")
		}
	} else if sig.Router() {
		if ps.router == nil {
			log.FatalStack("no router")
		}
	} else if sig.Raw() {
		if ps.raw == nil {
			log.FatalStack("no raw")
		}
	}

	return ps, err
}

// MakePacket()
// Read raw byte in, parse, load data into appropriate data structures
func MakePacket(raw []byte) (p *PacketStruct, err error) {

	var rawsize uint16 = uint16(len(raw))
	var calcsize uint16 = 0

	p = &PacketStruct{}

	p.pSig = PacketSigType(BtoU32(raw[:4]))
	if !p.pSig.V1() {
		log.Fatalf("Bad Sig Version:%04X", p.pSig)
	}
	raw = raw[4:]
	calcsize += 4

	if !p.pSig.V1() {
		return nil, fmt.Errorf("Bad Version %s", (p.pSig & MASK_VERSION))
	}

	//
	// Counters
	//
	if p.pSig.Counter() || p.pSig.Ack() || p.pSig.Nak() {
		if p.pSig.Width32() {
			p.counter = counter.NewByteCounter32(raw[:4])
			raw = raw[4:]
			calcsize += 4
		} else if p.pSig.Width64() {
			p.counter = counter.NewByteCounter64(raw[:8])
			raw = raw[8:]
			calcsize += 8
		} else {
			panic("no counter input")
		}
	}

	//
	// Ping/Pong
	//
	if p.pSig.Ping() {
		if p.pSig.Width32() {
			p.ping = pinger.NewBytePing32(raw[:4])
			raw = raw[4:]
			calcsize += 4
		} else if p.pSig.Width64() {
			p.ping = pinger.NewBytePing64(raw[:8])
			raw = raw[8:]
			calcsize += 8
		} else {
			panic("no counter input")
		}
	}

	if p.pSig.Pong() {
		if p.pSig.Width32() {
			p.ping = pinger.NewBytePong32(raw[:4])
			raw = raw[4:]
			calcsize += 4
		} else if p.pSig.Width64() {
			p.ping = pinger.NewBytePong64(raw[:8])
			raw = raw[8:]
			calcsize += 8
		} else {
			panic("no counter input")
		}
	}

	if p.pSig.Packet() {
		p.packet, err = MakePacket(raw)
		p.pSize += p.packet.Size()
	} else if p.pSig.IPV4() {
		p.ipv4, err = MakeIPv4(raw)
		p.pSize += p.ipv4.Size()
	} else if p.pSig.IPV6() {
		p.ipv6, err = MakeIPv6(raw)
		p.pSize += p.ipv6.Size()
	} else if p.pSig.Eth() {
		p.eth, err = MakeEth(raw)
		p.pSize += p.eth.Size()
	} else if p.pSig.Router() {
		p.router, err = MakeRouter(raw)
		p.pSize += p.router.Size()
	} else if p.pSig.ID() {
		p.id, err = MakeID(raw)
		p.pSize += p.id.Size()
	} else if p.pSig.Auth() {
		p.auth, err = MakeAuth(raw)
		p.pSize += p.auth.Size()
	} else if p.pSig.Raw() {
		p.raw = raw
		p.pSize += PacketSizeType(len(raw))
	}

	if p.pSig.Size8() {
		p.pSize += 1
	} else if p.pSig.Size16() {
		p.pSize += 2
	}

	if PacketSizeType(rawsize) != p.pSize {
		log.FatalStack("sizes dont match")
	}

	return p, err
}

// -
// ToByte()
// -
func (p *PacketStruct) ToByte() (raw []byte) {

	//
	// pSig
	//
	raw = append(raw, U32toB(uint32(p.pSig))...)

	//
	// pSize
	//
	if p.pSig.Size8() {
		raw = append(raw, U8toB(uint8(p.pSize))...)
	} else if p.pSig.Size16() {
		raw = append(raw, U16toB(uint16(p.pSize))...)
	} else {
		log.FatalStack("No size")
	}

	//
	// Counter/Ack/Nak
	//
	if p.pSig.Counter() || p.pSig.Ack() || p.pSig.Nak() {
		raw = append(raw, p.counter.ToByte()...)
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
		raw = append(raw, p.packet.ToByte()...)
	} else if p.pSig.IPV4() {
		raw = append(raw, p.ipv4.ToByte()...)
	} else if p.pSig.IPV6() {
		raw = append(raw, p.ipv6.ToByte()...)
	} else if p.pSig.Eth() {
		raw = append(raw, p.eth.ToByte()...)
	} else if p.pSig.Router() {
		raw = append(raw, p.router.ToByte()...)
	} else if p.pSig.Auth() {
		raw = append(raw, p.auth.ToByte()...)
	} else if p.pSig.ID() {
		raw = append(raw, p.id.ToByte()...)
	} else if p.pSig.Raw() {
		raw = append(raw, p.raw...)
	}

	return raw
}

func (p *PacketStruct) Sig() PacketSigType {
	return p.pSig
}

func (p *PacketStruct) Size() PacketSizeType {
	return p.pSize
}

func (p *PacketStruct) Width() PacketWidth {
	return p.pWidth
}

func (p *PacketStruct) Counter() counter.Counter {
	return p.counter
}

func (p *PacketStruct) Ping() pinger.Ping {
	return p.ping
}

func (p *PacketStruct) Pong() pinger.Pong {
	return p.pong
}

func (p *PacketStruct) Router() *RouterStruct {
	return p.router
}

func (p *PacketStruct) IPv4() *IPv4Struct {
	return p.ipv4
}

func (p *PacketStruct) IPv6() *IPv6Struct {
	return p.ipv6
}

func (p *PacketStruct) Eth() *EthStruct {
	return p.eth
}

func (p *PacketStruct) Auth() *AuthStruct {
	return p.auth
}

func (p *PacketStruct) ID() *IDStruct {
	return p.id
}

func (p *PacketStruct) Packet() *PacketStruct {
	return p.packet
}

func (p *PacketStruct) Raw() []byte {
	return p.raw
}
