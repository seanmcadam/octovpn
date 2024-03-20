package network

import (
	"net"
	"testing"
	"time"

	"github.com/seanmcadam/bufferpool"
	"github.com/seanmcadam/ctx"
	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/testlib/netlib"
)

func TestCompile(t *testing.T) {}

func TestConnection(t *testing.T) {
	var server, client interfaces.LayerInterface

	b1data := []byte("Hi There Cli")
	b2data := []byte("Hi There Srv")
	cx := ctx.New()
	ip := net.IPv4(127, 0, 0, 1)
	port := int(netlib.GetRandomNetworkPort())
	addr := &net.TCPAddr{IP: ip, Port: port}
	p := bufferpool.New()
	b1 := p.Get().Append(b1data)
	b2 := p.Get().Append(b2data)

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

	server.Send(b1)
	client.Send(b2)

}
