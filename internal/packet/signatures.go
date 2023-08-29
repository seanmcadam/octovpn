package packet

// Signatures uint8 0xVS: V=Version S=Packet Signature
//
// Version = 1
const (
	ROUTE_SIGV1 PacketSig = 0x10
	SITE_SIGV1  PacketSig = 0x11
	CHAN_SIGV1  PacketSig = 0x12
	CONN_SIGV1  PacketSig = 0x13
)
