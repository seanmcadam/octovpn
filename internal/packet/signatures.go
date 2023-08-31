package packet

import "github.com/seanmcadam/octovpn/octolib/log"

type PacketSigType uint16
type PacketFields uint32

//
// Layout of the PacketSigType
// This provides the ability to & values from certain columns
// 0xVLDP
// //||||
// //|||l__ Packet info (counts, pings, acks, higher layer data)
// //||l___ Data packets (IP,Eth,Routing, app messaging)
// //|l____ Which Layer did the packet originate
// //l_____ What Version is this packet (V1 for now)
//

// Layout of the Field Bits
// Predefined fiels that should be present in a packet
// 0xABCDEFGH
// //|\||\|||
// //| \| \||
// //|  |  \| Data Structure pointers (ie *RouterStruct)
// //|  | Single Value object (ping,pong,counter)
// //| Single Value Modifier or raw []byte values in stead of objects
const (
	Field_Sig           PacketFields  = 0x00000001
	Field_Payload       PacketFields  = 0x00000010
	Field_Router        PacketFields  = 0x00000020
	Field_IPV6          PacketFields  = 0x00000040
	Field_IPV4          PacketFields  = 0x00000080
	Field_Eth           PacketFields  = 0x00000100
	Field_Auth          PacketFields  = 0x00000200
	Field_ID            PacketFields  = 0x00000400
	field_Counter       PacketFields  = 0x00010000
	field_Ping          PacketFields  = 0x00020000
	field_Pong          PacketFields  = 0x00040000
	field_32            PacketFields  = 0x10000000
	field_64            PacketFields  = 0x20000000
	Field_Raw           PacketFields  = 0x80000000
	Field_Counter32     PacketFields  = field_Counter | field_32
	Field_Counter64     PacketFields  = field_Counter | field_64
	Field_Ping32        PacketFields  = field_Ping | field_32
	Field_Ping64        PacketFields  = field_Ping | field_64
	Field_Pong32        PacketFields  = field_Pong | field_32
	Field_Pong64        PacketFields  = field_Pong | field_64
	FieldSize_Sig       uint8         = 2
	FieldSize_Counter32 uint8         = 4
	FieldSize_Counter64 uint8         = 8
	FieldSize_Ping32    uint8         = 4
	FieldSize_Pong32    uint8         = 4
	FieldSize_Ping64    uint8         = 8
	FieldSize_Pong64    uint8         = 8
	Packet_VERSION      PacketSigType = 0xF000
	Packet_LAYER        PacketSigType = 0x0F00
	Packet_TYPE         PacketSigType = 0x00FF
	Packet_VER1         PacketSigType = 0x1000
	Packet_VER2         PacketSigType = 0x2000
	Packet_VER3         PacketSigType = 0x4000
	Packet_VER4         PacketSigType = 0x8000
	// Add more Layers here
	Packet_ROUT PacketSigType = 0x0800
	Packet_SITE PacketSigType = 0x0400
	Packet_CHAN PacketSigType = 0x0200
	Packet_CONN PacketSigType = 0x0100
	// Add more packet types here
	Packet_PARENT    PacketSigType = 0x0000
	Packet_RAW       PacketSigType = 0x0001
	Packet_AUTH      PacketSigType = 0x0002
	Packet_ACK       PacketSigType = 0x0003
	Packet_NAK       PacketSigType = 0x0004
	Packet_PING32    PacketSigType = 0x0005
	Packet_PONG32    PacketSigType = 0x0006
	Packet_PING64    PacketSigType = 0x0007
	Packet_PONG64    PacketSigType = 0x0008
	Packet_COUNTER32 PacketSigType = 0x0009
	Packet_COUNTER64 PacketSigType = 0x000A
	Packet_ID        PacketSigType = 0x000B
	Packet_ERROR     PacketSigType = 0x000F
	// Add more data types here
	Packet_IPV4   PacketSigType = 0x0010
	Packet_IPV6   PacketSigType = 0x0020
	Packet_ETH    PacketSigType = 0x0030
	Packet_ROUTER PacketSigType = 0x0040

	Packet_ROUTE_RAW       = Packet_VER1 | Packet_ROUT | Packet_RAW
	Packet_ROUTE_RAW_SIZE  = FieldSize_Sig
	Packet_ROUTE_RAW_FIELD = Field_Sig | Field_Payload

	Packet_ROUTE_AUTH       = Packet_VER1 | Packet_ROUT | Packet_AUTH
	Packet_ROUTE_AUTH_SIZE  = FieldSize_Sig
	Packet_ROUTE_AUTH_FIELD = Field_Sig | Field_Payload

	Packet_ROUTE_ROUTER       = Packet_VER1 | Packet_ROUT | Packet_ROUTER
	Packet_ROUTE_ROUTER_SIZE  = FieldSize_Sig
	Packet_ROUTE_ROUTER_FIELD = Field_Sig | Field_Router

	Packet_ROUTE_ETH       = Packet_VER1 | Packet_ROUT | Packet_ETH
	Packet_ROUTE_ETH_SIZE  = FieldSize_Sig
	Packet_ROUTE_ETH_FIELD = Field_Sig | Field_Eth

	Packet_ROUTE_IPV4       = Packet_VER1 | Packet_ROUT | Packet_IPV4
	Packet_ROUTE_IPV4_SIZE  = FieldSize_Sig
	Packet_ROUTE_IPV4_FIELD = Field_Sig | Field_IPV4

	Packet_ROUTE_IPV6       = Packet_VER1 | Packet_ROUT | Packet_IPV6
	Packet_ROUTE_IPV6_SIZE  = FieldSize_Sig
	Packet_ROUTE_IPV6_FIELD = Field_Sig | Field_IPV6

	Packet_SITE_RAW       = Packet_VER1 | Packet_SITE | Packet_RAW
	Packet_SITE_RAW_SIZE  = FieldSize_Sig
	Packet_SITE_RAW_FIELD = Field_Sig | Field_Payload

	Packet_SITE_PARENT       = Packet_VER1 | Packet_SITE | Packet_PARENT
	Packet_SITE_PARENT_SIZE  = FieldSize_Sig
	Packet_SITE_PARENT_FIELD = Field_Sig | Field_Payload

	Packet_SITE_AUTH       = Packet_VER1 | Packet_SITE | Packet_AUTH
	Packet_SITE_AUTH_SIZE  = FieldSize_Sig
	Packet_SITE_AUTH_FIELD = Field_Sig | Field_Payload

	Packet_SITE_ID       = Packet_VER1 | Packet_SITE | Packet_ID
	Packet_SITE_ID_SIZE  = FieldSize_Sig
	Packet_SITE_ID_FIELD = Field_Sig | Field_Payload

	Packet_CHAN_RAW       = Packet_VER1 | Packet_CHAN | Packet_RAW
	Packet_CHAN_RAW_SIZE  = FieldSize_Sig + FieldSize_Counter64
	Packet_CHAN_RAW_FIELD = Field_Sig | Field_Payload | Field_Counter64

	Packet_CHAN_PARENT       = Packet_VER1 | Packet_CHAN | Packet_PARENT
	Packet_CHAN_PARENT_SIZE  = FieldSize_Sig + FieldSize_Counter64
	Packet_CHAN_PARENT_FIELD = Field_Sig | Field_Payload | Field_Counter64

	Packet_CHAN_ACK       = Packet_VER1 | Packet_CHAN | Packet_ACK
	Packet_CHAN_ACK_SIZE  = FieldSize_Sig + FieldSize_Counter64
	Packet_CHAN_ACK_FIELD = Field_Sig | Field_Counter64

	Packet_CHAN_NAK       = Packet_VER1 | Packet_CHAN | Packet_NAK
	Packet_CHAN_NAK_SIZE  = FieldSize_Sig + FieldSize_Counter64
	Packet_CHAN_NAK_FIELD = Field_Sig | Field_Counter64

	Packet_CONN_RAW       = Packet_VER1 | Packet_CONN | Packet_RAW
	Packet_CONN_RAW_SIZE  = FieldSize_Sig
	Packet_CONN_RAW_FIELD = Field_Sig | Field_Payload

	Packet_CONN_PARENT       = Packet_VER1 | Packet_CONN | Packet_PARENT
	Packet_CONN_PARENT_SIZE  = FieldSize_Sig
	Packet_CONN_PARENT_FIELD = Field_Sig | Field_Payload

	Packet_CONN_AUTH       = Packet_VER1 | Packet_CONN | Packet_AUTH
	Packet_CONN_AUTH_SIZE  = FieldSize_Sig
	Packet_CONN_AUTH_FIELD = Field_Sig | Field_Payload

	Packet_CONN_PING64       = Packet_VER1 | Packet_CONN | Packet_PING64
	Packet_CONN_PING64_SIZE  = FieldSize_Sig + FieldSize_Ping64
	Packet_CONN_PING64_FIELD = Field_Sig | Field_Ping64

	Packet_CONN_PONG64       = Packet_VER1 | Packet_CONN | Packet_PONG64
	Packet_CONN_PONG64_SIZE  = FieldSize_Sig + FieldSize_Pong64
	Packet_CONN_PONG64_FIELD = Field_Sig | Field_Pong64
)

var PacketFieldTypes map[PacketFields]string

func init() {
	PacketFieldTypes = make(map[PacketFields]string, 12)
	PacketFieldTypes[Field_Sig] = "PacketSigType"
	PacketFieldTypes[Field_Payload] = "*PacketStruct"
	PacketFieldTypes[Field_Router] = "[]byte"
	PacketFieldTypes[Field_IPV6] = "[]byte"
	PacketFieldTypes[Field_IPV4] = "[]byte"
	PacketFieldTypes[Field_Eth] = "[]byte"
	PacketFieldTypes[Field_Counter32] = "*counter.Counter32"
	PacketFieldTypes[Field_Counter64] = "*counter.Counter64"
	PacketFieldTypes[Field_Ping32] = "*pinger.Ping32"
	PacketFieldTypes[Field_Ping64] = "*pinger.Ping64"
	PacketFieldTypes[Field_Pong32] = "*pinger.Pong32"
	PacketFieldTypes[Field_Pong64] = "*pinger.Pong64"
	PacketFieldTypes[Field_Raw] = "[]byte"

}
func (p PacketSigType) V1() bool {
	return ((p & Packet_VERSION) ^ Packet_VER1) == 0
}
func (p PacketSigType) GetLayer() PacketSigType {
	return (p & Packet_LAYER)
}

func (p PacketSigType) RouterLayer() bool {
	return (p&Packet_LAYER)^Packet_ROUT == 0
}
func (p PacketSigType) SiteLayer() bool {
	return (p&Packet_LAYER)^Packet_SITE == 0
}
func (p PacketSigType) ChanLayer() bool {
	return (p&Packet_LAYER)^Packet_CHAN == 0
}
func (p PacketSigType) ConnLayer() bool {
	return (p&Packet_LAYER)^Packet_CONN == 0
}

func (p PacketSigType) Parent() bool {
	return p&Packet_TYPE == 0
}
func (p PacketSigType) Raw() bool {
	return (p&Packet_TYPE)^Packet_RAW == 0
}
func (p PacketSigType) Auth() bool {
	return (p&Packet_TYPE)^Packet_AUTH == 0
}
func (p PacketSigType) Ack() bool {
	return (p&Packet_TYPE)^Packet_ACK == 0
}
func (p PacketSigType) Nak() bool {
	return (p&Packet_TYPE)^Packet_NAK == 0
}
func (p PacketSigType) Counter32() bool {
	return (p&Packet_TYPE)^Packet_COUNTER32 == 0
}
func (p PacketSigType) Counter64() bool {
	return (p&Packet_TYPE)^Packet_COUNTER64 == 0
}
func (p PacketSigType) Ping32() bool {
	return (p&Packet_TYPE)^Packet_PING32 == 0
}
func (p PacketSigType) Pong32() bool {
	return (p&Packet_TYPE)^Packet_PONG32 == 0
}
func (p PacketSigType) Ping64() bool {
	return (p&Packet_TYPE)^Packet_PING64 == 0
}
func (p PacketSigType) Pong64() bool {
	return (p&Packet_TYPE)^Packet_PONG64 == 0
}
func (p PacketSigType) Error() bool {
	return (p&Packet_TYPE)^Packet_ERROR == 0
}
func (p PacketSigType) ID() bool {
	return (p&Packet_TYPE)^Packet_ID == 0
}
func (p PacketSigType) IPV4() bool {
	return (p&Packet_TYPE)^Packet_IPV4 == 0
}
func (p PacketSigType) IPV6() bool {
	return (p&Packet_TYPE)^Packet_IPV6 == 0
}
func (p PacketSigType) Eth() bool {
	return (p&Packet_TYPE)^Packet_ETH == 0
}
func (p PacketSigType) Router() bool {
	return (p&Packet_TYPE)^Packet_ROUTER == 0
}

func (f PacketFields) PayloadSize() bool {
	return uint32(f&Field_Payload) > 0
}
func (f PacketFields) RouterSize() bool {
	return uint32(f&Field_Router) > 0
}
func (f PacketFields) IPV6Size() bool {
	return uint32(f&Field_IPV6) > 0
}
func (f PacketFields) IPV4Size() bool {
	return uint32(f&Field_IPV4) > 0
}
func (f PacketFields) EthSize() bool {
	return uint32(f&Field_Eth) > 0
}
func (f PacketFields) Counter32() bool {
	return uint32(f&Field_Counter32) > 0
}
func (f PacketFields) Counter64() bool {
	return uint32(f&Field_Counter64) > 0
}
func (f PacketFields) Ping32() bool {
	return uint32(f&Field_Ping32) > 0
}
func (f PacketFields) Ping64() bool {
	return uint32(f&Field_Ping64) > 0
}
func (f PacketFields) Pong32() bool {
	return uint32(f&Field_Pong32) > 0
}
func (f PacketFields) Pong64() bool {
	return uint32(f&Field_Pong64) > 0
}

func (f PacketSigType) String() (name string) {
	switch f {
	case Packet_VERSION:
		name = "VERSION"
	case Packet_LAYER:
		name = "LAYER"
	case Packet_VER1:
		name = "VER1"
	case Packet_VER2:
		name = "VER2"
	case Packet_VER3:
		name = "VER3"
	case Packet_VER4:
		name = "VER4"
	case Packet_ROUT:
		name = "ROUT"
	case Packet_SITE:
		name = "SITE"
	case Packet_CHAN:
		name = "CHAN"
	case Packet_CONN:
		name = "CONN"
	case Packet_PARENT:
		name = "PARENT"
	case Packet_RAW:
		name = "RAW"
	case Packet_AUTH:
		name = "AUTH"
	case Packet_ACK:
		name = "ACK"
	case Packet_NAK:
		name = "NAK"
	case Packet_COUNTER32:
		name = "COUNTER32"
	case Packet_COUNTER64:
		name = "COUNTER64"
	case Packet_PING32:
		name = "PING32"
	case Packet_PONG32:
		name = "PONG32"
	case Packet_PING64:
		name = "PING64"
	case Packet_PONG64:
		name = "PONG64"
	case Packet_ID:
		name = "ID"
	case Packet_ERROR:
		name = "ERROR"
	case Packet_IPV4:
		name = "IPV4"
	case Packet_IPV6:
		name = "IPV6"
	case Packet_ETH:
		name = "ETH"
	case Packet_ROUTER:
		name = "ROUTER"
	case Packet_ROUTE_RAW:
		name = "ROUTE_RAW"
	case Packet_ROUTE_AUTH:
		name = "ROUTE_AUTH"
	case Packet_ROUTE_ROUTER:
		name = "ROUTE_ROUTER"
	case Packet_ROUTE_ETH:
		name = "ROUTE_ETH"
	case Packet_ROUTE_IPV4:
		name = "ROUTE_IPV4"
	case Packet_ROUTE_IPV6:
		name = "ROUTE_IPV6"
	case Packet_SITE_RAW:
		name = "SITE_RAW"
	case Packet_SITE_PARENT:
		name = "SITE_PARENT"
	case Packet_SITE_AUTH:
		name = "SITE_AUTH"
	case Packet_SITE_ID:
		name = "SITE_ID"
	case Packet_CHAN_RAW:
		name = "CHAN_RAW"
	case Packet_CHAN_PARENT:
		name = "CHAN_PARENT"
	case Packet_CHAN_ACK:
		name = "CHAN_ACK"
	case Packet_CHAN_NAK:
		name = "CHAN_NAK"
	case Packet_CONN_RAW:
		name = "CHAN_RAW"
	case Packet_CONN_PARENT:
		name = "CHAN_PARENT"
	case Packet_CONN_AUTH:
		name = "CHAN_AUTH"
	case Packet_CONN_PING64:
		name = "CHAN_PING64"
	case Packet_CONN_PONG64:
		name = "CHAN_PONG64"
	default:
		log.FatalfStack("Unhandled Type:%04X", uint16(f))
	}
	return name
}
