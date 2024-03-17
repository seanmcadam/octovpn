package tcp

import (
	"net"

	"github.com/seanmcadam/ctx"
	"github.com/seanmcadam/network/connection"
	"github.com/seanmcadam/octovpn/interfaces"
)

// Client()
// Will connect to a target host:port
// If the connection closes it will reconnect
func Client(cx *ctx.Ctx, addr net.Addr) (ch chan interfaces.LayerInterface, err error) {

	ch = make(chan interfaces.LayerInterface, 1)

	go func() {
		defer func() {
			close(ch)
		}()
		for {
			conn, err := net.Dial("tcp", addr.String())
			if err != nil {
				return
			}

			clientcx := cx.WithCancel()
			ch <- connection.Connection(clientcx, conn)

			select {
			case <-clientcx.DoneChan():
			}

			if cx.Done() {
				return
			}

		}
	}()

	return ch, err
}
