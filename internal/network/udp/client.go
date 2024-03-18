package udp

import (
	"net"

	"github.com/seanmcadam/ctx"
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/interfaces"
)

func Client(cx *ctx.Ctx, addr net.Addr) (chan interfaces.LayerInterface, error) {
	log.Printf("udp Client")
	return nil, nil
}
