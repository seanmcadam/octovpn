package tcp

import (
	"net"

	"github.com/seanmcadam/ctx"
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/interfaces"
)

func Server(cx *ctx.Ctx, addr net.Addr) (ch chan interfaces.LayerInterface, err error) {

	ch = make(chan interfaces.LayerInterface, 1)




	return ch, err
}
