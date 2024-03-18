package tcp

import (
	"net"
	"testing"

	"github.com/seanmcadam/ctx"
	"github.com/seanmcadam/testlib/netlib"
)

func TestCompile(t *testing.T) {
}

// Setup a listen socket
// Create a client to connect to listen socket
// Write and read data over the sockets
func TestConnection(t *testing.T) {

	//var srvcomm interfaces.LayerInterface
	//var clicomm interfaces.LayerInterface

	cx := ctx.New()
	ip := net.IPv4(127, 0, 0, 1)
	port := int(netlib.GetRandomNetworkPort())

	addr := &net.TCPAddr{IP: ip, Port: port}

	//n := addr.Network()
	//ap := addr.AddrPort()

	srvch, srverr := Server(cx, addr)
	if srverr != nil {
		t.Errorf("Server Error:%s", srverr)
	}

	clich, clierr := Client(cx, addr)
	if clierr != nil {
		t.Errorf("Server Error:%s", clierr)
	}

	srvcomm := <-srvch
	clicomm := <-clich

	_ = srvcomm
	clicomm.Reset()
}
