package network

import (
	"net"

	"github.com/seanmcadam/ctx"
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/network/tcp"
	"github.com/seanmcadam/octovpn/internal/network/udp"
)

func Client(cx *ctx.Ctx, addr net.Addr) (chan interfaces.LayerInterface, error) {
	switch addr.Network() {
	case "tcp":
		return tcp.Client(cx, addr)
	case "udp":
		return udp.Client(cx, addr)
	default:
		log.Panic()
	}

	return nil, nil
}

func Server(cx *ctx.Ctx, addr net.Addr) (ch chan interfaces.LayerInterface, err error) {
	switch addr.Network() {
	case "tcp":
		ch, _, err = tcp.Server(cx, addr)
	case "udp":
		ch, err = udp.Server(cx, addr)
	default:
		log.Panic()
	}

	return
}
