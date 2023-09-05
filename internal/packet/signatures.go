package packet

type PacketSigType uint32
type PacketSize8FieldType uint8
type PacketSize16FieldType uint16
type PacketSizeType uint16
type PacketWidth uint8

const (
	PacketWidth0        PacketWidth    = 0
	PacketWidth32       PacketWidth    = 32
	PacketWidth64       PacketWidth    = 64
	PacketSigSize       PacketSizeType = 4
	PacketSize8Size     PacketSizeType = 1
	PacketSize16Size    PacketSizeType = 2
	PacketCounter32Size PacketSizeType = 4
	PacketCounter64Size PacketSizeType = 8
	PacketAck32Size     PacketSizeType = 4
	PacketAck64Size     PacketSizeType = 8
	PacketNak32Size     PacketSizeType = 4
	PacketNak64Size     PacketSizeType = 8
	PacketPing32Size    PacketSizeType = 4
	PacketPing64Size    PacketSizeType = 8
	PacketPong32Size    PacketSizeType = 4
	PacketPong64Size    PacketSizeType = 8
)

//type PacketFields uint32
//type PacketIDString []byte
//type PacketAuthString []byte

//
// Layout of the PacketSigType
// This provides the ability to & values from certain columns
// 0xA B CD EF GH
// //| |  |  |  |
// //| |  |  |  l__ Packet info (counts, pings, acks, higher layer data)
// //| |  |  l___ Data packets (IP,Eth,Routing, app messaging)
// //| |  l____ Which Layer did the packet originate
// //| l____ What size are the 32/64 bit fields
// //l_____ What Version is this packet (V1 for now)
// //
// //

// Layout of the Field Bits
// Predefined fiels that should be present in a packet
// 0xA BCD EFGH
// //| \|| \|||
// //|  \|  \||
// //|   |   \| Data Structure pointers (ie *RouterStruct, raw []byte)
// //|   | Single Value object (ping,pong,counter,id,auth)
// //| Size Modifier

//
// Order of Fields
// Sig
// Size
// Counters
// Ping/Pongs | ID | Auth
// Payloads (all of them)
// Raw

const (
	MASK_VERSION    PacketSigType = 0xF0000000
	MASK_SIZE       PacketSigType = 0x0F000000
	MASK_WIDTH_64_32 PacketSigType = 0x0C000000
	MASK_SIZE_16_8  PacketSigType = 0x03000000
	MASK_LAYER      PacketSigType = 0x00FF0000
	MASK_FIELDS     PacketSigType = 0x0000FF00
	MASK_DATA       PacketSigType = 0x000000FF
	VERSION_1       PacketSigType = 0x10000000
	VERSION_2       PacketSigType = 0x20000000 // Could double the packet size for ceratin layers
	VERSION_3       PacketSigType = 0x40000000 // Could mix and match versions in differet layers
	SIZE_8          PacketSigType = 0x01000000
	SIZE_16         PacketSigType = 0x02000000
	WIDTH_32         PacketSigType = 0x04000000
	WIDTH_64         PacketSigType = 0x08000000

	LAYER_ROUTER PacketSigType = 0x00080000
	LAYER_SITE   PacketSigType = 0x00040000
	LAYER_CHAN   PacketSigType = 0x00020000
	LAYER_CONN   PacketSigType = 0x00010000

	FIELD_SIG     PacketSigType = 0x00000100 // This data... put in for completeness
	FIELD_SIZE    PacketSigType = 0x00000200 // Requires SIZE_8/16
	FIELD_COUNTER PacketSigType = 0x00000400 // Requires WIDTH_32/64
	FIELD_ACK     PacketSigType = 0x00000800 // Requires WIDTH_32/64
	FIELD_NAK     PacketSigType = 0x00001000 // Requires WIDTH_32/64
	MASK_CAN      PacketSigType = FIELD_COUNTER | FIELD_ACK | FIELD_NAK
	FIELD_PING    PacketSigType = 0x00002000 // Requires WIDTH_32/64
	FIELD_PONG    PacketSigType = 0x00004000 // Requires WIDTH_32/64
	MASK_PINGPONG PacketSigType = FIELD_PING | FIELD_PONG
	FIELD_DATA    PacketSigType = 0x00008000 // data size provided with Size()

	DATA_PACKET PacketSigType = 0x00000001
	DATA_RAW    PacketSigType = 0x00000002
	DATA_AUTH   PacketSigType = 0x00000004
	DATA_ID     PacketSigType = 0x00000008
	DATA_IPV4   PacketSigType = 0x00000010
	DATA_IPV6   PacketSigType = 0x00000020
	DATA_ETH    PacketSigType = 0x00000030
	DATA_ROUTER PacketSigType = 0x00000040
	DATA_ERROR  PacketSigType = 0x000000FF

	SIG_ROUTE = VERSION_1 | LAYER_ROUTER | FIELD_SIG | FIELD_SIZE | FIELD_DATA
	SIG_SITE  = VERSION_1 | LAYER_SITE | FIELD_SIG | FIELD_SIZE
	SIG_CHAN  = VERSION_1 | LAYER_CHAN | FIELD_SIG | FIELD_SIZE
	SIG_CONN  = VERSION_1 | LAYER_CONN | FIELD_SIG | FIELD_SIZE

	//
	// ROUTE
	//

	SIG_ROUTE_AUTH   = SIG_ROUTE | DATA_AUTH | SIZE_8
	SIG_ROUTE_RAW    = SIG_ROUTE | DATA_RAW | SIZE_16
	SIG_ROUTE_ETH    = SIG_ROUTE | DATA_ETH | SIZE_16
	SIG_ROUTE_IPV4   = SIG_ROUTE | DATA_IPV4 | SIZE_16
	SIG_ROUTE_IPV6   = SIG_ROUTE | DATA_IPV6 | SIZE_16
	SIG_ROUTE_ROUTER = SIG_ROUTE | DATA_ROUTER | SIZE_16

	//
	// SITE
	//
	SIG_SITE_32_RAW    = SIG_SITE | FIELD_DATA | FIELD_COUNTER | DATA_RAW | WIDTH_32 | SIZE_16
	SIG_SITE_32_PACKET = SIG_SITE | FIELD_DATA | FIELD_COUNTER | DATA_PACKET | WIDTH_32 | SIZE_16
	SIG_SITE_32_AUTH   = SIG_SITE | FIELD_DATA | FIELD_COUNTER | DATA_AUTH | WIDTH_32 | SIZE_16
	SIG_SITE_32_ID     = SIG_SITE | FIELD_DATA | FIELD_COUNTER | DATA_ID | WIDTH_32 | SIZE_16
	SIG_SITE_32_ACK    = SIG_SITE | FIELD_ACK | WIDTH_32 | SIZE_16
	SIG_SITE_32_NAK    = SIG_SITE | FIELD_NAK | WIDTH_32 | SIZE_16
	SIG_SITE_64_RAW    = SIG_SITE | FIELD_DATA | FIELD_COUNTER | DATA_RAW | WIDTH_64 | SIZE_16
	SIG_SITE_64_PACKET = SIG_SITE | FIELD_DATA | FIELD_COUNTER | DATA_PACKET | WIDTH_64 | SIZE_16
	SIG_SITE_64_AUTH   = SIG_SITE | FIELD_DATA | FIELD_COUNTER | DATA_AUTH | WIDTH_64 | SIZE_16
	SIG_SITE_64_ID     = SIG_SITE | FIELD_DATA | FIELD_COUNTER | DATA_ID | WIDTH_64 | SIZE_16
	SIG_SITE_64_ACK    = SIG_SITE | FIELD_ACK | WIDTH_64 | SIZE_16
	SIG_SITE_64_NAK    = SIG_SITE | FIELD_NAK | WIDTH_64 | SIZE_16
	//
	// CHAN
	//
	SIG_CHAN_32_RAW    = SIG_CHAN | FIELD_DATA | FIELD_COUNTER | DATA_RAW | WIDTH_32 | SIZE_16
	SIG_CHAN_32_PACKET = SIG_CHAN | FIELD_DATA | FIELD_COUNTER | DATA_PACKET | WIDTH_32 | SIZE_16
	SIG_CHAN_32_PING   = SIG_CHAN | FIELD_PING | FIELD_COUNTER | WIDTH_32 | SIZE_8
	SIG_CHAN_32_PONG   = SIG_CHAN | FIELD_PONG | FIELD_COUNTER | WIDTH_32 | SIZE_8
	SIG_CHAN_32_ACK    = SIG_CHAN | FIELD_ACK | WIDTH_32 | SIZE_8
	SIG_CHAN_32_NAK    = SIG_CHAN | FIELD_NAK | WIDTH_32 | SIZE_8
	SIG_CHAN_64_RAW    = SIG_CHAN | FIELD_DATA | FIELD_COUNTER | DATA_RAW | WIDTH_64 | SIZE_16
	SIG_CHAN_64_PACKET = SIG_CHAN | FIELD_DATA | FIELD_COUNTER | DATA_PACKET | WIDTH_64 | SIZE_16
	SIG_CHAN_64_PING   = SIG_CHAN | FIELD_PING | FIELD_COUNTER | WIDTH_64 | SIZE_8
	SIG_CHAN_64_PONG   = SIG_CHAN | FIELD_PONG | FIELD_COUNTER | WIDTH_64 | SIZE_8
	SIG_CHAN_64_ACK    = SIG_CHAN | FIELD_ACK | WIDTH_64 | SIZE_8
	SIG_CHAN_64_NAK    = SIG_CHAN | FIELD_NAK | WIDTH_64 | SIZE_8
	//
	// CHAN
	//
	SIG_CONN_32_RAW    = SIG_CONN | FIELD_DATA | FIELD_COUNTER | DATA_RAW | WIDTH_32 | SIZE_16
	SIG_CONN_32_PACKET = SIG_CONN | FIELD_DATA | FIELD_COUNTER | DATA_PACKET | WIDTH_32 | SIZE_16
	SIG_CONN_32_PING   = SIG_CONN | FIELD_PING | FIELD_COUNTER | WIDTH_32 | SIZE_8
	SIG_CONN_32_PONG   = SIG_CONN | FIELD_PONG | FIELD_COUNTER | WIDTH_32 | SIZE_8
	SIG_CONN_32_AUTH   = SIG_CONN | FIELD_DATA | FIELD_COUNTER | DATA_AUTH | WIDTH_32 | SIZE_16
	SIG_CONN_32_ID     = SIG_CONN | FIELD_DATA | FIELD_COUNTER | DATA_ID | WIDTH_32 | SIZE_16
	SIG_CONN_64_RAW    = SIG_CONN | FIELD_DATA | FIELD_COUNTER | DATA_RAW | WIDTH_64 | SIZE_16
	SIG_CONN_64_PACKET = SIG_CONN | FIELD_DATA | FIELD_COUNTER | DATA_PACKET | WIDTH_64 | SIZE_16
	SIG_CONN_64_PING   = SIG_CONN | FIELD_PING | FIELD_COUNTER | WIDTH_64 | SIZE_8
	SIG_CONN_64_PONG   = SIG_CONN | FIELD_PONG | FIELD_COUNTER | WIDTH_64 | SIZE_8
	SIG_CONN_64_AUTH   = SIG_CONN | FIELD_DATA | FIELD_COUNTER | DATA_AUTH | WIDTH_64 | SIZE_16
	SIG_CONN_64_ID     = SIG_CONN | FIELD_DATA | FIELD_COUNTER | DATA_ID | WIDTH_64 | SIZE_16
)

func (p PacketSigType) Size8() bool {
	return (p & SIZE_8) != 0
}
func (p PacketSigType) Size16() bool {
	return (p & SIZE_16) != 0
}
func (p PacketSigType) Width0() bool {
	return (p & (WIDTH_32|WIDTH_64)) == 0
}
func (p PacketSigType) Width32() bool {
	return (p & WIDTH_32) != 0
}
func (p PacketSigType) Width64() bool {
	return (p & WIDTH_64) != 0
}
func (p PacketSigType) V1() bool {
	return (p & VERSION_1) != 0
}

func (p PacketSigType) GetLayer() PacketSigType {
	return (p & MASK_LAYER)
}

func (p PacketSigType) RouterLayer() bool {
	return (p & LAYER_ROUTER) != 0
}
func (p PacketSigType) SiteLayer() bool {
	return (p & LAYER_SITE) != 0
}
func (p PacketSigType) ChanLayer() bool {
	return (p & LAYER_CHAN) != 0
}
func (p PacketSigType) ConnLayer() bool {
	return (p & LAYER_CONN) != 0
}

func (p PacketSigType) Packet() bool {
	return (p & DATA_PACKET) != 0
}
func (p PacketSigType) Raw() bool {
	return (p & DATA_RAW) != 0
}
func (p PacketSigType) Auth() bool {
	return (p & DATA_AUTH) != 0
}
func (p PacketSigType) Ack() bool {
	return (p & FIELD_ACK) != 0
}
func (p PacketSigType) Nak() bool {
	return (p & FIELD_NAK) != 0
}
func (p PacketSigType) Counter() bool {
	return (p & FIELD_COUNTER) != 0
}
func (p PacketSigType) Ping() bool {
	return (p & FIELD_PING) != 0
}
func (p PacketSigType) Pong() bool {
	return (p & FIELD_PONG) != 0
}
func (p PacketSigType) ID() bool {
	return (p & DATA_ID) != 0
}
func (p PacketSigType) IPV4() bool {
	return (p & DATA_IPV4) != 0
}
func (p PacketSigType) IPV6() bool {
	return (p & DATA_IPV6) != 0
}
func (p PacketSigType) Eth() bool {
	return (p & DATA_ETH) != 0
}
func (p PacketSigType) Router() bool {
	return (p & DATA_ROUTER) != 0
}

func (f PacketSigType) String() string {

	var name = []byte("")

	if (f & MASK_VERSION) == MASK_VERSION {
		name = append(name, []byte("MASK_VERSION ")...)
	} else {
		name = append(name, []byte("ERR-MASK_VERSION ")...)
	}
	if (f & MASK_SIZE) == MASK_SIZE {
		name = append(name, []byte("MASK_SIZE")...)
	}
	if (f & MASK_WIDTH_64_32) == MASK_WIDTH_64_32 {
		name = append(name, []byte("MASK_WIDTH_64_32")...)
	}
	if (f & MASK_SIZE_16_8) == MASK_SIZE_16_8 {
		name = append(name, []byte("MASK_SIZE_16_8")...)
	}
	if (f & MASK_LAYER) == MASK_LAYER {
		name = append(name, []byte("MASK_LAYER")...)
	}
	if (f & MASK_FIELDS) == MASK_FIELDS {
		name = append(name, []byte("MASK_FIELDS")...)
	}
	if (f & MASK_DATA) == MASK_DATA {
		name = append(name, []byte("DATA_ERROR")...)
	}

	if (f & MASK_VERSION) == VERSION_1 {
		name = append(name, []byte("VERSION_1 ")...)
	}

	if (f & MASK_WIDTH_64_32) == WIDTH_64 {
		name = append(name, []byte("WIDTH_64 ")...)
	} else if (f & MASK_WIDTH_64_32) == WIDTH_32 {
		name = append(name, []byte("WIDTH_32 ")...)
	}

	if (f & MASK_SIZE_16_8) == SIZE_16 {
		name = append(name, []byte("SIZE_16 ")...)
	} else if (f & MASK_SIZE_16_8) == SIZE_8 {
		name = append(name, []byte("SIZE_8 ")...)
	}

	if (f & MASK_LAYER) == LAYER_CONN {
		name = append(name, []byte("CONN ")...)
	} else if (f & MASK_LAYER) == LAYER_CHAN {
		name = append(name, []byte("CHAN ")...)
	} else if (f & MASK_LAYER) == LAYER_SITE {
		name = append(name, []byte("SITE ")...)
	} else if (f & MASK_LAYER) == LAYER_ROUTER {
		name = append(name, []byte("ROUTER ")...)
	} else {
		name = append(name, []byte("ERR-LAYER ")...)
	}

	if (f & FIELD_SIG) > 0 {
		name = append(name, []byte("SIG ")...)
	}
	if (f & FIELD_SIZE) > 0 {
		name = append(name, []byte("SIZE ")...)
	}
	if (f & FIELD_COUNTER) > 0 {
		name = append(name, []byte("COUNTER ")...)
	}
	if (f & FIELD_ACK) > 0 {
		name = append(name, []byte("ACK ")...)
	}
	if (f & FIELD_NAK) > 0 {
		name = append(name, []byte("NAK ")...)
	}
	if (f & FIELD_PING) > 0 {
		name = append(name, []byte("PING ")...)
	}
	if (f & FIELD_PONG) > 0 {
		name = append(name, []byte("PONG ")...)
	}
	if (f & FIELD_DATA) > 0 {
		name = append(name, []byte("DATA ")...)
	}

	if (f & DATA_PACKET) > 0 {
		name = append(name, []byte("PACKET ")...)
	} else if (f & DATA_RAW) > 0 {
		name = append(name, []byte("RAW ")...)
	} else if (f & DATA_AUTH) > 0 {
		name = append(name, []byte("AUTH ")...)
	} else if (f & DATA_ID) > 0 {
		name = append(name, []byte("ID ")...)
	} else if (f & DATA_IPV4) > 0 {
		name = append(name, []byte("IPV4 ")...)
	} else if (f & DATA_IPV6) > 0 {
		name = append(name, []byte("IPV6 ")...)
	} else if (f & DATA_ETH) > 0 {
		name = append(name, []byte("ETH ")...)
	} else if (f & DATA_ROUTER) > 0 {
		name = append(name, []byte("ROUTER ")...)
	} else if (f & DATA_ERROR) > 0 {
		name = append(name, []byte("ROUTER ")...)
	}

	return string(name)
}
