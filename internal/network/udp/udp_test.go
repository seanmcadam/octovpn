package udp

import (
	"net"
	"testing"
	"time"

	"github.com/seanmcadam/bufferpool"
	"github.com/seanmcadam/ctx"
	"github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/testlib/netlib"
)

// Used for debugging, set high for manual, low for automated
var establishTimeout time.Duration = time.Second

func TestCompile(t *testing.T) {
}

func TestMakeConnection(t *testing.T) {
	//establishTimeout = time.Minute
	_, _, _ = MakeConnection(t)
}

// Setup a listen socket
// Create a client to connect to listen socket
// Write and read data over the sockets
func TestConnection(t *testing.T) {

	var pool *bufferpool.Pool
	pool = bufferpool.New()

	b1data := []byte("Hi There Cli")
	b2data := []byte("Hi There Srv")

	server, srvconn, cliconn := MakeConnection(t)

	bufcli := pool.Get()
	bufcli.Append(b1data)
	bufsrv := pool.Get()
	bufsrv.Append(b2data)

	cliconn.Send(bufcli)
	srvconn.Send(bufsrv)

	select {
	case clidata := <-cliconn.RecvCh():
		if clidata == nil {
			t.Errorf("Cli RecvCh return nil")
			return
		}
		loggy.Debug("Select Cli")
		b1 := clidata.Data()
		if string(b1) != string(b2data) {
			t.Errorf("Cli Read not equal")
		}
	case srvdata := <-srvconn.RecvCh():
		if srvdata == nil {
			t.Errorf("Srv RecvCh return nil")
			return
		}
		loggy.Debug("Select Srv")
		b2 := srvdata.Data()
		if string(b2) != string(b1data) {
			t.Errorf("Srv Read not equal")
		}
	case <-time.After(establishTimeout):
		t.Errorf("Select Time out")
	}

	server.Close()
}

func MakeConnection(t *testing.T) (server *UDPServerStruct, srvlayer interfaces.LayerInterface, clilayer interfaces.LayerInterface) {

	pool := bufferpool.New()

	cx := ctx.New()
	ip := net.IPv4(127, 0, 0, 1)
	port := int(netlib.GetRandomNetworkPort())
	addr := &net.UDPAddr{IP: ip, Port: port}

	srvch, server, srverr := Server(cx, addr)
	if srverr != nil {
		t.Errorf("Server Error:%s", srverr)
	}

	time.Sleep(time.Millisecond * 100)

	clich, clierr := Client(cx, addr)
	if clierr != nil {
		t.Errorf("Server Error:%s", clierr)
	}

	select {
	case clilayer = <-clich:
	case <-time.After(establishTimeout):
		t.Fatalf("SRVCH timeout")
	}

	// Send Empty buffer to jump start the connection
	b := pool.Get()
	clilayer.Send(b)
	select {
	case srvlayer = <-srvch:
	case <-time.After(establishTimeout):
		t.Fatalf("SRVCH timeout")
	}

	bsrv := <-srvlayer.RecvCh()
	_ = bsrv

	return server, srvlayer, clilayer
}
