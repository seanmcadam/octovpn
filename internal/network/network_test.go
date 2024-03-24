package network

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/seanmcadam/bufferpool"
	"github.com/seanmcadam/ctx"
	"github.com/seanmcadam/octovpn/interfaces/layers"
	"github.com/seanmcadam/testlib/netlib"
)

// Used for debugging, set high for manual, low for automated
var establishTimeout time.Duration = time.Second

func TestCompile(t *testing.T) {}

func TestLayersConnection(t *testing.T) {
	var server, client layers.LayerInterface

	server, client = layers.CreateLayerPair()
	connectionValidate(t, server, client)
}

func TestTCPConnection(t *testing.T) {
	var server, client layers.LayerInterface
	var err error

	server, client, err = createTCPPair(t)
	if err != nil {
		t.Errorf("Error:%s", err)
		return
	}
	connectionValidate(t, server, client)

}

func TestUDPConnection(t *testing.T) {
	var server, client layers.LayerInterface
	var err error

	server, client, err = createUDPPair(t)
	if err != nil {
		t.Errorf("Error:%s", err)
		return
	}
	connectionValidate(t, server, client)

}

//func Test10SecPacketSend(t *testing.T, server layers.LayerInterface, client layers.LayerInterface) {
//
//	establishTimeout = time.Minute
//
//	var pool *bufferpool.Pool
//	pool = bufferpool.New()
//	senddone := make(chan bool)
//	recvdone := make(chan bool)
//
//	b1data := []byte("Hi There Cli")
//	b2data := []byte("Hi There Srv")
//
//	server, srvconn, cliconn := MakeConnection(t)
//
//	bufcli := pool.Get()
//	bufcli.Append(b1data)
//	bufsrv := pool.Get()
//	bufsrv.Append(b2data)
//
//	go func() {
//		time.Sleep(10 * time.Second)
//		senddone <- true
//		time.Sleep(1 * time.Second)
//		recvdone <- true
//	}()
//
//	go func() {
//		for {
//			select {
//			case <-senddone:
//				break
//			default:
//			}
//			cliconn.Send(bufcli.Copy())
//			srvconn.Send(bufsrv.Copy())
//			time.Sleep(time.Millisecond * 10)
//		}
//	}()
//
//FORREAD:
//	for {
//		select {
//		case clidata := <-cliconn.RecvCh():
//			if clidata == nil {
//				t.Errorf("Cli RecvCh return nil")
//				return
//			}
//			loggy.Debug("Select Cli")
//			b1 := clidata.Data()
//			if string(b1) != string(b2data) {
//				t.Errorf("Cli Read not equal")
//			}
//			clidata.ReturnToPool()
//		case srvdata := <-srvconn.RecvCh():
//			if srvdata == nil {
//				t.Errorf("Srv RecvCh return nil")
//				return
//			}
//			loggy.Debug("Select Srv")
//			b2 := srvdata.Data()
//			if string(b2) != string(b1data) {
//				t.Errorf("Srv Read not equal")
//			}
//			srvdata.ReturnToPool()
//		case <-recvdone:
//			break FORREAD
//		}
//	}
//
//	server.Close()
//}

func createTCPPair(t *testing.T) (server layers.LayerInterface, client layers.LayerInterface, err error) {
	cx := ctx.New()
	ip := net.IPv4(127, 0, 0, 1)
	port := int(netlib.GetRandomNetworkPort())
	addr := &net.TCPAddr{IP: ip, Port: port}

	return createPair(t, cx, addr)
}

func createUDPPair(t *testing.T) (server layers.LayerInterface, client layers.LayerInterface, err error) {
	cx := ctx.New()
	ip := net.IPv4(127, 0, 0, 1)
	port := int(netlib.GetRandomNetworkPort())
	addr := &net.UDPAddr{IP: ip, Port: port}

	return createPair(t, cx, addr)
}

func createPair(t *testing.T, cx *ctx.Ctx, addr net.Addr) (server layers.LayerInterface, client layers.LayerInterface, err error) {
	serverch, err := Server(cx, addr)
	if err != nil {
		t.Errorf("Server Err:%s", err)
		return nil, nil, fmt.Errorf("Server Err:%s", err)
	}
	clientch, err := Client(cx, addr)
	if err != nil {
		t.Errorf("Client Err:%s", err)
		return nil, nil, fmt.Errorf("Client Err:%s", err)
	}

	select {
	case server = <-serverch:
		if server == nil {
			t.Errorf("Bad Server")
			return nil, nil, fmt.Errorf("No Server")
		}
	case <-time.After(time.Second * 5):
		t.Errorf("Server timeout")
		return nil, nil, fmt.Errorf("Server Timeout")
	}

	select {
	case client = <-clientch:
		if client == nil {
			t.Errorf("Bad Client")
			return nil, nil, fmt.Errorf("No Client")
		}
	case <-time.After(time.Second):
		t.Errorf("Client timeout")
		return nil, nil, fmt.Errorf("Client Timeout")
	}

	return server, client, err
}

func connectionValidate(t *testing.T, server layers.LayerInterface, client layers.LayerInterface) {

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
