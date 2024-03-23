package network

import (
	"net"
	"testing"
	"time"

	"github.com/seanmcadam/bufferpool"
	"github.com/seanmcadam/ctx"
	"github.com/seanmcadam/octovpn/interfaces/layers"
	"github.com/seanmcadam/testlib/netlib"
)

func TestCompile(t *testing.T) {}

func TestLayersConnection(t *testing.T) {
	var server, client layers.LayerInterface

	server, client = layers.CreateLayerPair()
	connection(t, server, client)
}

func TestTCPConnection(t *testing.T) {
	var server, client layers.LayerInterface
	var err error

	cx := ctx.New()
	ip := net.IPv4(127, 0, 0, 1)
	port := int(netlib.GetRandomNetworkPort())
	addr := &net.TCPAddr{IP: ip, Port: port}

	server, client, err = createPair(t, cx, addr)
	if err != nil {
		t.Errorf("Error:%s", err)
		return
	}
	connection(t, server, client)

}

func TestUDPConnection(t *testing.T) {
	var server, client layers.LayerInterface
	var err error

	cx := ctx.New()
	ip := net.IPv4(127, 0, 0, 1)
	port := int(netlib.GetRandomNetworkPort())
	addr := &net.UDPAddr{IP: ip, Port: port}

	server, client, err = createPair(t, cx, addr)
	if err != nil {
		t.Errorf("Error:%s", err)
		return
	}
	connection(t, server, client)

}

func createPair(t *testing.T, cx *ctx.Ctx, addr net.Addr) (server layers.LayerInterface, client layers.LayerInterface, err error) {
	serverch, err := Server(cx, addr)
	if err != nil {
		t.Errorf("Server Err:%s", err)
	}
	clientch, err := Client(cx, addr)
	if err != nil {
		t.Errorf("Client Err:%s", err)
	}

	select {
	case server = <-serverch:
		if server == nil {
			t.Errorf("Bad Server")
			return
		}
	case <-time.After(time.Second * 5):
		t.Errorf("Server timeout")
		return
	}

	select {
	case client = <-clientch:
		if client == nil {
			t.Errorf("Bad Client")
			return
		}
	case <-time.After(time.Second):
		t.Errorf("Client timeout")
		return
	}
	return server, client, err
}

func connection(t *testing.T, server layers.LayerInterface, client layers.LayerInterface) {

	p := bufferpool.New()
	b1data := []byte("Hi There Cli")
	b2data := []byte("Hi There Srv")
	b1 := p.Get().Append(b1data)
	b2 := p.Get().Append(b2data)

	server.Send(b1)
	client.Send(b2)

	select {
	case <-server.RecvCh():
	case <-client.RecvCh():
	case <-time.After(time.Second):
		t.Errorf("Timeout on connection")
		return
	}

	t.Logf("connection test completed")
}
