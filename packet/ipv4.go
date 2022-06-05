package packet

import (
	"fmt"
	"net"
)

type IPv4Frame []byte

type IPProtocol uint8

const IP_ICMP IPProtocol = 0x01
const IP_IGMP IPProtocol = 0x02
const IP_TCP IPProtocol = 0x06
const IP_UDP IPProtocol = 0x11
const IP_GRE IPProtocol = 0x2F
const IP_ESP IPProtocol = 0x32
const IP_AH IPProtocol = 0x33

//
// Version bits 0-3 (O-0) (Always equals 4)
// IHL      Bits 4-7 (O-0)(Always equals 4)
//          The IPv4 header is variable in size due to the options.
//          The IHL field contains the size of the IPv4 header, made up of 4 bits that specify the number of 32-bit words in the header.
//          The MIN value for this field is 5 - 5 × 32 bits = 160 bits = 20 bytes.
//          The MAX value is 15 - 15 × 32 bits = 480 bits = 60 bytes
// DSCP     Bits 8-13 (O-1)
// ECN      Bits 14-15 (O-1)
// Length   Octet 2-3
// Ident    Octet 4-5
// Flags    Bits 32-34(0-2) (O-6)
// Fragment Bits 35-63(3-7,0-7) (O-6-7)
// TTL      Octet 8
// Protocol Octet 9
// Checksum Octet 10-11
// Source   Octet 12-15
// Dest     Octet 16-19
// Options  Octet 20-59
//

//
// ipv4IHL()
//
func (f IPv4Frame) IHL() (ihl uint8) {

	var b byte

	b = f[0]
	b = b & 0xf0
	b = b >> 4

	if b < 5 {
		panic("IHL < 5")
	}

	if b > 15 {
		panic("IHL > 15")
	}

	return uint8(b)
}

//
// ipv4DSCP()
//

//
// ipv4ECN()
//

//
// ipv4TotalLenth()
//
func (f IPv4Frame) TotalLenth() (len uint16) {

	len = uint16(f[2])
	len = len << 8
	len = len + uint16(f[3])

	if len < 20 {
		panic("len < 20")
	}

	return len
}

//
// IP4Identification()
//

//
// IPv4DF()
//
func (f IPv4Frame) DF() (df bool) {
	var b byte
	df = false
	b = f[6]
	b = b & 0x40
	if b > 0 {
		df = true
	}
	return df
}

//
// IPv4MF()
//
func (f IPv4Frame) MF() (mf bool) {
	var b byte
	mf = false
	b = f[6]
	b = b & 0x20
	if b > 0 {
		mf = true
	}
	return mf
}

//
// IPv4Fragment()
//
func (f IPv4Frame) Fragment() (fragment uint16) {
	var b [2]byte
	b[0] = f[6]
	b[1] = f[7]

	b[0] = b[0] & 0x31

	fragment = uint16(b[0])
	fragment = fragment << 8
	fragment = fragment + uint16(b[1])

	return fragment
}

//
// TTL()
//
func (f IPv4Frame) TTL() (ttl uint8) {
	ttl = uint8(f[8])
	return ttl
}

//
// Protocol()
//
func (f IPv4Frame) Protocol() IPProtocol {
	switch IPProtocol(f[9]) {
	case IP_ICMP:
		return IP_ICMP
	case IP_IGMP:
		return IP_IGMP
	case IP_TCP:
		return IP_TCP
	case IP_UDP:
		return IP_UDP
	case IP_GRE:
		return IP_GRE
	case IP_ESP:
		return IP_ESP
	case IP_AH:
		return IP_AH
	default:
		panic(fmt.Sprintf("IP Protocol not supported: 0x%x", f[9]))
	}
}

//
// Checksum()
//

//
// PayloadOffset()
//
func (f IPv4Frame) PayloadOffset() (offset uint8) {
	offset = f.IHL() * 4
	return offset
}

//
// ipv4SourceIP()
// octets 12-15
func (f IPv4Frame) Source() (ip net.IP) {
	ip = make(net.IP, 4)
	ip[0] = f[12]
	ip[1] = f[13]
	ip[2] = f[14]
	ip[3] = f[15]
	return ip
}

//
// ipv4DestIP()
// octets 16-19
func (f IPv4Frame) Dest() (ip net.IP) {
	ip = make(net.IP, 4)
	ip[0] = f[16]
	ip[1] = f[17]
	ip[2] = f[18]
	ip[3] = f[19]
	return ip
}

//
//
//
func (p IPProtocol) String() (proocol string) {
	switch p {
	case IP_ICMP:
		return "ICMP"
	case IP_IGMP:
		return "IGMP"
	case IP_TCP:
		return "TCP"
	case IP_UDP:
		return "UDP"
	case IP_GRE:
		return "GRE"
	case IP_ESP:
		return "ESP"
	case IP_AH:
		return "AH"
	default:
		panic(fmt.Sprintf("IP Protocol not supported: 0x%d", p))
	}
}
