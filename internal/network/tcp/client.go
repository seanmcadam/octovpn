package tcp

import (
	"net"

	"github.com/seanmcadam/ctx"
	"github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/network/connection"
)

// Client()
// Will connect to a target host:port
// If the connection closes it will reconnect
func Client(cx *ctx.Ctx, addr net.Addr) (ch chan interfaces.LayerInterface, err error) {

	ch = make(chan interfaces.LayerInterface, 1)

	go func(cx *ctx.Ctx) {
		defer func() {
			loggy.Debugf("Net Client() Defer Close %s", addr.String())
			close(ch)
		}()
		for {
			loggy.Debugf("Net Client() %s", addr.String())

			conn, err := net.Dial("tcp", addr.String())
			if err != nil {
				return
			}

			clientcx := cx.WithCancel()
			ch <- connection.Connection(clientcx, conn)

			select {
			case <-clientcx.DoneChan():
				loggy.Debugf("Net Client() Closed %s", addr.String())
			}

			if cx.Done() {
				return
			}

		}
	}(cx)

	return ch, err
}
