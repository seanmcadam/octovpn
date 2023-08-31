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

// type PacketSize uint16
type PacketCounter32 counter.Counter32
type PacketCounter64 counter.Counter64
type PacketPing32 pinger.Ping32
type PacketPong32 pinger.Pong32
type PacketPing64 pinger.Ping64
type PacketPong64 pinger.Pong64

type PacketStruct struct {
	pSig      PacketSigType
	counter32 *PacketCounter32
	counter64 *PacketCounter64
	ping32    *PacketPing32
	ping64    *PacketPing64
	pong32    *PacketPong32
	pong64    *PacketPong64
	router    *RouterStruct
	ipv6      *IPv6Struct
	ipv4      *IPv4Struct
	eth       *EthStruct
	auth      *AuthStruct
	id        *IDStruct
	raw       []byte
	payload   *PacketStruct
}

func NewPacket(sig PacketSigType, v ...interface{}) (ps *PacketStruct, err error) {

	if !sig.V1() {
		return nil, fmt.Errorf("Bad Packet Version1: 0x%02X", uint16(sig&Packet_VERSION))
	}

	ps = &PacketStruct{
		pSig: sig,
	}

	for _, i := range v {
		switch u := i.(type) {
		case []byte:
			if !ps.pSig.Raw() {
				log.Errorf("Got RAW type for %s", ps.pSig)
			}
			ps.raw = u
		case *PacketStruct:
			if !ps.pSig.Parent() {
				log.Errorf("Got Parent type for %s", ps.pSig)
			}
			ps.payload = u
		case *AuthStruct:
			if !ps.pSig.Auth() {
				log.Errorf("Got Auth type for %s", ps.pSig)
			}
			ps.auth = u
		case *IDStruct:
			if !ps.pSig.ID() {
				log.Errorf("Got ID type for %s", ps.pSig)
			}
			ps.id = u
		case *RouterStruct:
			if !ps.pSig.Router() {
				log.Errorf("Got Router type for %s", ps.pSig)
			}
			ps.router = u
		case *IPv6Struct:
			if !ps.pSig.IPV6() {
				log.Errorf("Got IPv6 type for %s", ps.pSig)
			}
			ps.ipv6 = u
		case *IPv4Struct:
			if !ps.pSig.IPV4() {
				log.Errorf("Got IPv4 type for %s", ps.pSig)
			}
			ps.ipv4 = u
		case *EthStruct:
			if !ps.pSig.Eth() {
				log.Errorf("Got Eth type for %s", ps.pSig)
			}
			ps.eth = u
		case *PacketCounter32:
			if !ps.pSig.Counter32() {
				log.Errorf("Got Counter32 type for %s", ps.pSig)
			}
			ps.counter32 = u
		case *PacketCounter64:
			if !ps.pSig.Counter64() {
				log.Errorf("Got Counter64 type for %s", ps.pSig)
			}
			ps.counter64 = u
		case *PacketPing32:
			if !ps.pSig.Ping32() {
				log.Errorf("Got Ping32 type for %s", ps.pSig)
			}
			ps.ping32 = u
		case *PacketPing64:
			if !ps.pSig.Ping64() {
				log.Errorf("Got Ping64 type for %s", ps.pSig)
			}
			ps.ping64 = u
		case *PacketPong32:
			if !ps.pSig.Pong32() {
				log.Errorf("Got Pong32 type for %s", ps.pSig)
			}
			ps.pong32 = u
		case *PacketPong64:
			if !ps.pSig.Pong64() {
				return nil, fmt.Errorf("Got Pong64 type for %s", ps.pSig)
			}
			ps.pong64 = u
		default:
			log.FatalfStack("Default Reached Type:%t", v)
		}
	}
	return ps, err
}
