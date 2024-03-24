// 
// network package allows for generic access to the underlying protocols (TCP and UDP for now)
// It returns layer interfaces, to seemlessly integrate with the above layers
//
package network

import (
	"net"

	"github.com/seanmcadam/ctx"
	"github.com/seanmcadam/loggy"
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/interfaces/layers"
	"github.com/seanmcadam/octovpn/internal/network/tcp"
	"github.com/seanmcadam/octovpn/internal/network/udp"
)

func Client(cx *ctx.Ctx, addr net.Addr) (chan layers.LayerInterface, error) {
	switch addr.Network() {
	case "tcp":
		fallthrough
	case "tcp4":
		fallthrough
	case "tcp6":
		return tcp.Client(cx, addr)
	case "udp":
		fallthrough
	case "udp4":
		fallthrough
	case "udp6":
		udpAddr, ok := addr.(*net.UDPAddr)
		if !ok {
			return nil, loggy.Err("addr is not a UDP address")
		}
		return udp.Client(cx, udpAddr)
	default:
		log.Panic()
	}

	return nil, nil
}

func Server(cx *ctx.Ctx, addr net.Addr) (ch chan layers.LayerInterface, err error) {
	switch addr.Network() {
	case "tcp":
		fallthrough
	case "tcp4":
		fallthrough
	case "tcp6":
		ch, _, err = tcp.Server(cx, addr)
	case "udp":
		fallthrough
	case "udp4":
		fallthrough
	case "udp6":
		udpAddr, ok := addr.(*net.UDPAddr)
		if !ok {
			return nil, loggy.Err("addr is not a UDP address")
		}
		return udp.Client(cx, udpAddr)
	default:
		log.Panic()
	}

	return
}
