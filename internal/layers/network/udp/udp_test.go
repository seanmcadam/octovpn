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
	if err == nil {
		t.Error("Send() returned nil")
	}
	us.Cancel()
	_ = us.doneChan()
	_ = us.closed()
	us.RecvChan()
	us.goRecv()
	us.goRecvTimeout()
	us.emptyrecv()
	us.Link()
	us.run()
	us.sendclose()

	srv, cli, err := connection(cx)
	if err != nil {
		t.Error(err)
	}

	cli.endpoint()
	srv.endpoint()
	cli.emptyrecv()
	srv.emptysend()
	srv.Link()
	srv.doneChan()
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

	p, err := packet.TestConn32Packet()
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
	case <-time.After(UDPRecvTimeout):
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

	p, err := packet.TestConn32Packet()
	if err != nil {
		t.Error(err)
	}

	cli.Send(p)

	<-time.After(time.Millisecond)

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
		t.Error("Srv Recieve timeout")
		return
	}

	srv.Send(p)

	<-time.After(time.Millisecond)

	select {
	case r := <-cli.RecvChan():
		if r == nil {
			t.Error("Recieved nil")
			return
		}
		err = packet.Validatepackets(p, r)
	case <-time.After(2 * time.Second):
		t.Error("Cli Recieve timeout")
	}
}

// Validate the link system
func TestNewUDP_link_validation(t *testing.T) {
	cx := ctx.NewContext()
	defer cx.Cancel()

	srv, cli, err := connection(cx)
	if err != nil {
		t.Error(err)
	}

	srvUpCh := srv.link.LinkUpCh()
	cliUpCh := srv.link.LinkUpCh()
	srvDnCh := srv.link.LinkDownCh()
	cliDnCh := srv.link.LinkDownCh()
	srvCloseCh := srv.link.LinkCloseCh()
	cliCloseCh := srv.link.LinkCloseCh()

	select {
	case <-cliUpCh:
	case <-time.After(time.Millisecond):
		if cli.link.IsDown() {
			t.Error("Cli Up Timeout")
		}
	}

	p, err := packet.TestConn32Packet()
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
	case <-time.After(10 * time.Second):
		t.Error("Srv Recieve timeout")
		return
	}

	select {
	case <-srvUpCh:
	case <-time.After(time.Millisecond):
		if srv.link.IsDown() {
			t.Error("Srv Up Timeout")
		}
	}

	srv.Send(p)
	select {
	case r := <-cli.RecvChan():
		if r == nil {
			t.Error("Recieved nil")
			return
		}
		err = packet.Validatepackets(p, r)
	case <-time.After(time.Second):
		t.Error("Cli Recieve timeout")
	}

	srv.Cancel()

	select {
	case <-srvDnCh:
	case <-time.After(time.Millisecond):
		t.Error("Srv Dn Timeout")
	}

	select {
	case <-cliDnCh:
	case <-time.After(time.Millisecond):
		t.Error("Cli Dn Timeout")
	}

	select {
	case <-srvCloseCh:
	case <-time.After(time.Millisecond):
		t.Error("Srv Close Timeout")
	}

	select {
	case <-cliCloseCh:
	case <-time.After(time.Millisecond):
		t.Error("Cli Close Timeout")
	}

}

//-
// The Server should send a close... and go back to listening
// The Client should get the close, and close
// 
//-
func TestNewUDP_cli_send_bad_sig(t *testing.T) {
	cx := ctx.NewContext()
	defer cx.Cancel()

	srv, cli, err := connection(cx)
	if err != nil {
		t.Error(err)
	}

	srvCloseCh := srv.link.LinkCloseCh()
	cliCloseCh := cli.link.LinkCloseCh()

	raw := []byte{0, 0, 0, 0}
	cli.sendtestpacket(raw)

	<-time.After(time.Millisecond)

	select {
	case <-srvCloseCh:
	case <-time.After(100 * time.Millisecond):
		t.Error("Srv Close Timeout")
	}

	<-time.After(time.Millisecond)

	select {
	case <-cliCloseCh:
	case <-time.After(100 * time.Millisecond):
		t.Error("Cli Close Timeout")
	}
}

func TestNewUDP_srv_send_bad_sig(t *testing.T) {
	cx := ctx.NewContext()
	defer cx.Cancel()

	srv, cli, err := connection(cx)
	if err != nil {
		t.Error(err)
	}

	srvCloseCh := srv.link.LinkCloseCh()
	cliCloseCh := cli.link.LinkCloseCh()

	raw := []byte{0, 0, 0, 0}
	srv.sendtestpacket(raw)

	select {
	case <-cliCloseCh:
	case <-time.After(100 * time.Millisecond):
		t.Error("Cli Close Timeout")
	}

	select {
	case <-srvCloseCh:
	case <-time.After(100 * time.Millisecond):
		t.Error("Srv Close Timeout")
	}
}

func TestNewUDP_srv_send_short_packet(t *testing.T) {
	cx := ctx.NewContext()
	defer cx.Cancel()

	srv, cli, err := connection(cx)
	if err != nil {
		t.Error(err)
	}

	srvCloseCh := srv.link.LinkCloseCh()
	cliCloseCh := cli.link.LinkCloseCh()

	p, err := packet.TestConn32Packet()
	if err != nil {
		t.Error(err)
	}

	var raw1, raw2 []byte
	if raw, err := p.ToByte(); err != nil {
		t.Error("ToByte() Err:", err)
	} else {
		raw1 = raw[:len(raw)-3]
		raw2 = raw[len(raw)-3:]
	}

	srv.sendtestpacket(raw1)
	<-time.After(time.Millisecond)
	srv.sendtestpacket(raw2)

	select {
	case r := <-cli.RecvChan():
		if r == nil {
			//t.Error("Recieved nil")
			//return
		}
		err = packet.Validatepackets(p, r)
	case <-time.After(100 * time.Millisecond):
		t.Error("Cli Recieve timeout")
	}

	srv.Cancel()

	select {
	case <-srvCloseCh:
	case <-time.After(100 * time.Millisecond):
		t.Error("Srv Close Timeout")
	}

	select {
	case <-cliCloseCh:
	case <-time.After(10 * time.Second):
		t.Error("Cli Close Timeout")
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
