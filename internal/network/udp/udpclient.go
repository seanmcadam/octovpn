package udp

import (
	"net"

	"github.com/seanmcadam/bufferpool"
	"github.com/seanmcadam/ctx"
	"github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/interfaces/layers"
)

type UDPClientStruct struct {
	cx    *ctx.Ctx
	ch    chan layers.LayerInterface
	conn  *net.UDPConn
	raddr *net.UDPAddr
}

// Client()
// Will connect to a target host:port
// If the connection closes it will reconnect
func Client(cx *ctx.Ctx, raddr *net.UDPAddr) (ch chan layers.LayerInterface, err error) {

	ch = make(chan layers.LayerInterface, 1)
	client := &UDPClientStruct{
		cx:    cx,
		ch:    ch,
		raddr: raddr,
	}
	go client.goRecv()

	return ch, nil
}

// goRecv()
// Dial the remote side
// Create the connection()
// Loop on receive
func (uc *UDPClientStruct) goRecv() {

	pool := bufferpool.New()

	var raddrstr = uc.raddr.AddrPort().String()

	loggy.Debugf("Running(%s)", raddrstr)

	defer func() {
		loggy.Debugf(" Defer() Close %s", raddrstr)
		close(uc.ch)
	}()

	for {
		var buffer = make([]byte, 2048)
		var err error

		loggy.Debugf("Dialing() %s", raddrstr)

		uc.conn, err = net.DialUDP(uc.raddr.Network(), nil, uc.raddr)
		if err != nil {
			loggy.Debugf("DialUDP() Err:%s", err)
			return
		}

		loggy.Debugf("Dialed() %s", raddrstr)

		connection := NewConnection(uc.cx.WithCancel(), uc.conn, uc.raddr, nil)
		uc.ch <- layers.LayerInterface(connection)

		for !uc.cx.Done() {
			n, err := uc.conn.Read(buffer)
			if err != nil {
				loggy.Fatalf("Err:%s", err)
			}
			b := pool.Get().Append(buffer[:n])
			connection.PushRecv(b)
			//buffer = buffer[:0]
		}
	}
}
