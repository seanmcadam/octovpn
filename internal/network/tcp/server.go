package tcp

import (
	"net"
	"time"

	"github.com/seanmcadam/ctx"
	"github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/network/connection"
)

type TCPServer struct {
	cx       *ctx.Ctx
	listener *net.TCPListener
}

func Server(cx *ctx.Ctx, addr net.Addr) (ch chan interfaces.LayerInterface, server *TCPServer, err error) {

	server = &TCPServer{
		cx: cx,
	}

	ch = make(chan interfaces.LayerInterface, 1)
	laddr, err := net.ResolveTCPAddr(addr.Network(), addr.String())

	go func(cx *ctx.Ctx) {

		server.listener, err = net.ListenTCP(addr.Network(), laddr)
		if err != nil {
			loggy.FatalfStack("net.ListenTCP(%s) Err:%s", addr.String(), err)
			return
		}
		defer func() {
			loggy.Debugf("Net Server() Defer Close %s", addr.String())
			close(ch)
			server.listener.Close()
		}()

		loggy.Debugf("Net Listener(%s) Running", addr.String())

		for {
			select {
			case <-cx.DoneChan():
				return
			default:
			}

			var conn net.Conn
			conn, err = server.listener.Accept()

			if err == nil {
				select {
				case <-cx.DoneChan():
					return
				case ch <- connection.Connection(cx.WithCancel(), conn):
				case <-time.After(time.Second):
					loggy.FatalfStack("Server Accept(%s) Timeout ch %s", addr.String(), conn.RemoteAddr().String())
					continue
				}
				loggy.Debugf("Server Accept(%s) %s", addr.String(), conn.RemoteAddr().String())
			} else {
				loggy.Debugf("ServerNet Accept(%s) Error %s", addr.String(), err)
			}
		}
	}(cx)

	return ch, server, err
}

func (t *TCPServer) Close() {
	t.cx.Cancel()
}
