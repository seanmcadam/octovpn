package packet

import (
	"github.com/seanmcadam/octovpn/octolib/log"
)

//
// PING64 the layer sens a PING packet to the other side expecting a PONG response
//

// ROUTE 8 - Packets originating from the routing layer
// SITE  C - Packets coming form the site layer
// CHAN  E - Packets from the channel layer
// CONN  F -

// RAW 0
// AUTH 1
// PARENT 2
// ACK 3
// NAK 4
// PING64 A
// PONG64 B
// ERROR

const (
	ROUTE_TYPE_RAW   PacketType = 0x80
	ROUTE_TYPE_ROUTE PacketType = 0x81
	ROUTE_TYPE_ETH   PacketType = 0x82
	ROUTE_TYPE_IP4   PacketType = 0x83
	ROUTE_TYPE_IP6   PacketType = 0x84

	SITE_TYPE_RAW    PacketType = 0xC0
	SITE_TYPE_AUTH   PacketType = 0xC1
	SITE_TYPE_PARENT PacketType = 0xC2
	SITE_TYPE_ACK    PacketType = 0xC3
	SITE_TYPE_NAK    PacketType = 0xC4
	SITE_TYPE_PING64 PacketType = 0xCA
	SITE_TYPE_PONG64 PacketType = 0xCB
	SITE_TYPE_ERROR  PacketType = 0xCF

	CHAN_TYPE_RAW    PacketType = 0xE0
	CHAN_TYPE_AUTH   PacketType = 0xE1
	CHAN_TYPE_PARENT PacketType = 0xE2
	CHAN_TYPE_ACK    PacketType = 0xE3
	CHAN_TYPE_NAK    PacketType = 0xE4
	CHAN_TYPE_PING64 PacketType = 0xEA
	CHAN_TYPE_PONG64 PacketType = 0xEB
	CHAN_TYPE_ERROR  PacketType = 0xEF

	CONN_TYPE_RAW    PacketType = 0xF0
	CONN_TYPE_AUTH   PacketType = 0xF1
	CONN_TYPE_PARENT PacketType = 0xF2
	CONN_TYPE_ACK    PacketType = 0xF3
	CONN_TYPE_NAK    PacketType = 0xF4
	CONN_TYPE_PING64 PacketType = 0xFA
	CONN_TYPE_PONG64 PacketType = 0xFB
	CONN_TYPE_ERROR  PacketType = 0xFF
)

func (p PacketType) String() string {
	switch p {
	case ROUTE_TYPE_RAW:
		return "ROUTE RAW"
	case ROUTE_TYPE_ROUTE:
		return "ROUTE ROUTE"
	case ROUTE_TYPE_ETH:
		return "ROUTE ETH"
	case ROUTE_TYPE_IP4:
		return "ROUTE IPV4"
	case ROUTE_TYPE_IP6:
		return "ROUTE IPV6"

	case SITE_TYPE_RAW:
		return "SITE RAW"
	case SITE_TYPE_AUTH:
		return "SITE AUTH"
	case SITE_TYPE_PARENT:
		return "SITE PARENT"
	case SITE_TYPE_ACK:
		return "SITE ACK"
	case SITE_TYPE_NAK:
		return "SITE NAK"
	case SITE_TYPE_PING64:
		return "SITE PONG64"
	case SITE_TYPE_PONG64:
		return "SITE PONG64"
	case SITE_TYPE_ERROR:
		return "SITE ERROR"

	case CHAN_TYPE_RAW:
		return "CHAN RAW"
	case CHAN_TYPE_AUTH:
		return "CHAN AUTH"
	case CHAN_TYPE_PARENT:
		return "CHAN PARENT"
	case CHAN_TYPE_ACK:
		return "CHAN ACK"
	case CHAN_TYPE_NAK:
		return "CHAN NAK"
	case CHAN_TYPE_PING64:
		return "CHAN PING64"
	case CHAN_TYPE_PONG64:
		return "CHAN PONG64"
	case CHAN_TYPE_ERROR:
		return "CHAN ERROR"

	case CONN_TYPE_RAW:
		return "CONN RAW"
	case CONN_TYPE_AUTH:
		return "CONN AUTH"
	case CONN_TYPE_PARENT:
		return "CONN PARENT"
	case CONN_TYPE_ACK:
		return "CONN ACK"
	case CONN_TYPE_NAK:
		return "CONN NAK"
	case CONN_TYPE_PING64:
		return "CONN PING64"
	case CONN_TYPE_PONG64:
		return "CONN PONG64"
	case CONN_TYPE_ERROR:
		return "CONN ERROR"

	default:
		log.Fatalf("UNKNOWN CHAN TYPE:%d", p)
	}
	return ""
}
