package tcp

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestNewTCP_basic(t *testing.T) {
	cx := ctx.NewContext()
	defer cx.Cancel()
}

func TestNewTCP_test_nil_returns(t *testing.T) {
	cx := ctx.NewContext()
	defer cx.Cancel()

	var ts *TcpStruct
	NewTCP(nil, nil)
	ts.goSend()
	ts.sendpacket(nil)
	ts.emptysend()
	err := ts.Send(nil)
	if err == nil {
		t.Error("Send() returned nil")
	}
	ts.Cancel()
	_ = ts.doneChan()
	_ = ts.closed()
	ts.RecvChan()
	ts.goRecv()
	ts.emptyrecv()
	ts.Link()
	ts.run()

	srv, cli, err := connection(cx)
	if err != nil {
		t.Error(err)
	}

	cli.emptyrecv()
	srv.emptysend()

	srv.Link()
	srv.doneChan()
	srv.Cancel()

}

func TestNewTCP_send(t *testing.T) {
	cx := ctx.NewContext()
	defer cx.Cancel()

	srv, cli, err := connection(cx)
	if err != nil {
		t.Error(err)
	}

	p, err := packet.Testpacket()
	if err != nil {
		t.Error(err)
	}

	cli.Send(p)
	select {
	case r := <-srv.RecvChan():
		if r == nil {
			t.Error("Recieved nil")
			return
		}
		err = packet.Validatepackets(p, r)
		if r == nil {
			t.Error("Recieved nil")
			return
		}
	case <-time.After(2 * time.Second):
		t.Error("Recieve timeout")
		return
	}

	srv.Send(p)
	select {
	case r := <-cli.RecvChan():
		if r == nil {
			t.Error("Recieved nil")
			return
		}
		err = packet.Validatepackets(p, r)
	case <-time.After(2 * time.Second):
		t.Error("Recieve timeout")
	}

}

//-----------------------------------------------------------------------------

// -
//
// -
func connection(cx *ctx.Ctx) (srvconn *TcpStruct, cliconn *TcpStruct, err error) {

	srv, cli, err := createPairedConnections()
	if err != nil {
		return nil, nil, err
	}

	srvconn = NewTCP(cx, srv)
	cliconn = NewTCP(cx, cli)

	return srvconn, cliconn, err

}

func createPairedConnections() (*net.TCPConn, *net.TCPConn, error) {
	port := getRandomPort()

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: port})
	if err != nil {
		return nil, nil, err
	}

	conn1, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: port})
	if err != nil {
		listener.Close()
		return nil, nil, err
	}

	conn2, err := listener.AcceptTCP()
	if err != nil {
		listener.Close()
		conn1.Close()
		return nil, nil, err
	}

	fmt.Printf("Connections established on port %d\n", port)
	return conn1, conn2, nil
}

func handleConnection(conn *net.TCPConn, name string) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("[%s] Error reading: %v\n", name, err)
			return
		}

		receivedData := buffer[:n]
		fmt.Printf("[%s] Received: %s\n", name, receivedData)
	}
}

func getRandomPort() int {
	return rand.Intn(60000) + 1025
}
