package tcp

import (
	"net"

	"github.com/seanmcadam/ctx"
	"github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/network/connection"
)

func Server(cx *ctx.Ctx, addr net.Addr) (ch chan interfaces.LayerInterface, err error) {

	ch = make(chan interfaces.LayerInterface, 1)
	laddr, err := net.ResolveTCPAddr(addr.Network(), addr.String())

	go func(cx *ctx.Ctx) {
		defer func() {
			close(ch)
		}()

		listener, err := net.ListenTCP(addr.Network(), laddr)
		if err != nil {
			return
		}
		defer func() {
			listener.Close()
		}()
		loggy.Debugf("Net Listener() %s", addr.String())

		for {
			var conn net.Conn
			conn, err = listener.Accept()
			if err != nil {
				loggy.Debugf("Net Accept() Error %s:%s", addr.String(), err)
				continue
			}

			loggy.Debugf("Net Accept() %s:%s", addr.String(), conn.RemoteAddr().String())

			clientcx := cx.WithCancel()
			ch <- connection.Connection(clientcx, conn)

			if cx.Done() {
				return
			}
		}
	}(cx)

	return ch, err
}
