package udp

import (
	"net"
	"net/netip"
	"sync"

	"github.com/seanmcadam/bufferpool"
	"github.com/seanmcadam/ctx"
	"github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/interfaces/layers"
)

type UDPServerStruct struct {
	cx            *ctx.Ctx
	ch            chan layers.LayerInterface
	laddr         *net.UDPAddr
	conn          *net.UDPConn
	connections   map[netip.AddrPort]*UDPConnection
	connectionsmx sync.Mutex
}

//
// Open connection
// create Server Struct
// Accept Connection
// Return connection on chan
//
//

func Server(cx *ctx.Ctx, laddr *net.UDPAddr) (ch chan layers.LayerInterface, server *UDPServerStruct, err error) {

	laddr, err = net.ResolveUDPAddr(laddr.Network(), laddr.String())
	if err != nil {
		loggy.Fatalf("bad addr %s", laddr)
	}

	conn, err := net.ListenUDP(laddr.Network(), laddr)
	if err != nil {
		loggy.Fatalf("Error net.ListenTCP(%s) Err:%s", laddr.String(), err)
		return
	}

	loggy.Debugf("Listening %s", laddr)

	ch = make(chan layers.LayerInterface, 1)
	server = &UDPServerStruct{
		cx:          cx,
		ch:          ch,
		conn:        conn,
		laddr:       laddr,
		connections: make(map[netip.AddrPort]*UDPConnection, 0),
	}
	go server.goRecv()

	return ch, server, nil
}

func (u *UDPServerStruct) goRecv() {
	defer func(u *UDPServerStruct) {
		loggy.Debugf("goRecv() Defer() called %s", u.laddr)
	}(u)

	pool := bufferpool.New()

	for !u.cx.Done() {
		var buffer = make([]byte, 2048)
		n, clientAddr, err := u.conn.ReadFromUDP(buffer)
		if err != nil {
			loggy.Fatalf("Err:%s", err)
		}

		loggy.Debugf("%s->%s ReadFromUDP()", u.laddr.String(), clientAddr.String())

		b := pool.Get().Append(buffer[:n])
		addrPort := clientAddr.AddrPort()

		if _, ok := u.connections[addrPort]; !ok {
			u.connectionsmx.Lock()
			u.connections[addrPort] = NewConnection(u.cx.WithCancel(), u.conn, clientAddr, u.remove)
			u.connectionsmx.Unlock()
			u.ch <- layers.LayerInterface(u.connections[addrPort])
		}

		u.connections[addrPort].PushRecv(b)
	}
}

func (u *UDPServerStruct) remove(addrPort netip.AddrPort) {
	u.connectionsmx.Lock()
	defer u.connectionsmx.Unlock()
	delete(u.connections, addrPort)
	loggy.Printf("UDP Connection %s Removed", addrPort)
}

func (u *UDPServerStruct) Close() {
	u.cx.Cancel()
}
