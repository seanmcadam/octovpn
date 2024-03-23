package tcp

import (
	"net"
	"testing"
	"time"

	"github.com/seanmcadam/bufferpool"
	"github.com/seanmcadam/ctx"
	"github.com/seanmcadam/octovpn/interfaces/layers"
	"github.com/seanmcadam/testlib/netlib"
)

func TestCompile(t *testing.T) {
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

	srvconn.Send(bufsrv)
	cliconn.Send(bufcli)

	clidata := <-cliconn.RecvCh()
	if clidata == nil {
		t.Errorf("Cli RecvCh return nil")
		return
	}

	srvdata := <-srvconn.RecvCh()
	if srvdata == nil {
		t.Errorf("Srv RecvCh return nil")
		return
	}

	b1 := clidata.Data()
	if string(b1) != string(b2data) {
		t.Errorf("Cli Read not equal")
	}
	b2 := srvdata.Data()
	if string(b2) != string(b1data) {
		t.Errorf("Srv Read not equal")
	}


	server.Close()
}

func MakeConnection(t *testing.T) (server *TCPServer, srvcomm layers.LayerInterface, clicomm layers.LayerInterface) {

	cx := ctx.New()
	ip := net.IPv4(127, 0, 0, 1)
	port := int(netlib.GetRandomNetworkPort())
	addr := &net.TCPAddr{IP: ip, Port: port}

	srvch, server, srverr := Server(cx, addr)
	if srverr != nil {
		t.Errorf("Server Error:%s", srverr)
	}

	time.Sleep(time.Millisecond * 100)

	clich, clierr := Client(cx, addr)
	if clierr != nil {
		t.Errorf("Server Error:%s", clierr)
	}


	srvcomm = <-srvch
	clicomm = <-clich

	//server.Close()
	return
}
