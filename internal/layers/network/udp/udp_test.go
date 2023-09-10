package udp

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestNewUdp_basic(t *testing.T) {
	cx := ctx.NewContext()
	defer cx.Cancel()

	_, _, err := connection(cx)
	if err != nil {
		t.Errorf("Err:%s", err)
		return
	}
}

func TestNewUDP_test_nil_returns(t *testing.T) {
	cx := ctx.NewContext()

	var us *UdpStruct
	NewUDPCli(nil, nil)
	NewUDPSrv(nil, nil)
	us.goSend()
	us.sendpacket(nil)
	us.emptysend()
	us.endpoint()
	err := us.Send(nil)
	if err == nil{
		t.Error("Send() returned nil")
	}
	us.Cancel()
	_ = us.DoneChan()
	_ = us.closed()
	us.RecvChan()
	us.goRecv()
	us.emptyrecv()
	us.Link()
	us.run()

	srv, cli, err := connection(cx)
	if err != nil {
		t.Error(err)
	}

	cli.endpoint()
	srv.endpoint()
	cli.emptyrecv()
	srv.emptysend()
	srv.Link()
	srv.DoneChan()
	srv.Cancel()

	defer cx.Cancel()
}

func TestNewUdp_send_over_cli(t *testing.T) {
	cx := ctx.NewContext()
	defer cx.Cancel()
	srv, cli, err := connection(cx)
	if err != nil {
		t.Error(err)
		return
	}

	p, err := packet.Testpacket()
	if err != nil {
		t.Error(err)
		return
	}

	cli.Send(p)

	select {
	case r := <-srv.RecvChan():
		if r == nil {
			t.Error("Recieved nil")
			return
		}
	case <-time.After(2 * time.Second):
		t.Error("Recieve timeout")
		return
	}

}

func TestNewUdp_send_over_cli_srv(t *testing.T) {
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
func connection(cx *ctx.Ctx) (srvconn *UdpStruct, cliconn *UdpStruct, err error) {

	srv, cli, err := createPairedConnections()
	if err != nil {
		return nil, nil, err
	}

	//go handleConnection(srv, "Srv Connection")
	//go handleConnection(cli, "Cli Connection")

	cliconn = NewUDPCli(cx, cli)
	srvconn = NewUDPSrv(cx, srv)

	return srvconn, cliconn, err

}

// -
//
// -
func createPairedConnections() (srv *net.UDPConn, cli *net.UDPConn, err error) {
	port := getRandomPort()

	addr := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: port}

	srv, err = net.ListenUDP("udp", addr)
	if err != nil {
		return nil, nil, err
	}

	cli, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, nil, err
	}

	fmt.Printf("Connections established on port %d\n", port)
	return srv, cli, nil
}

// -
//
// -
func handleConnection(conn *net.UDPConn, name string) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("[%s] Error reading: %v\n", name, err)
			return
		}

		receivedData := buffer[:n]
		fmt.Printf("[%s] Received from %s: %s\n", name, addr.String(), receivedData)
	}
}

// -
//
// -
func getRandomPort() int {
	return rand.Intn(60000) + 1025
}
