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

// -
// Will the module compile
// -
func TestNewTCP_basic(t *testing.T) {
	cx := ctx.NewContext()
	defer cx.Cancel()
}

// -
// Verify the functions handle nil input with out panic
// -
func TestNewTCP_test_nil_returns(t *testing.T) {
	cx := ctx.NewContext()
	defer cx.Cancel()

	var ts *TcpStruct
	NewTCP(nil, nil)
	ts.goSend()
	ts.sendclose()
	ts.sendpacket(nil)
	ts.sendtestpacket(nil)
	ts.emptysend()
	err := ts.Send(nil)
	if err == nil {
		t.Error("Send() returned nil")
	}
	ts.Cancel()
	_ = ts.doneChan()
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

// -
// Send a packet, and recieve it
// -
func TestNewTCP_send(t *testing.T) {
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

// -
// Validate the link system
// -
func TestNewTCP_link_validation(t *testing.T) {
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
	case <-srvUpCh:
	case <-time.After(time.Millisecond):
		if srv.link.IsDown() {
			t.Error("Srv Up Timeout")
		}
	}

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
	case <-time.After(time.Second):
		t.Error("Srv Recieve timeout")
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

// -
//
// -
func TestNewTCP_cli_send_bad_sig(t *testing.T) {
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

	select {
	case <-srvCloseCh:
	case <-time.After(100 * time.Millisecond):
		t.Error("Srv Close Timeout")
	}

	select {
	case <-cliCloseCh:
	case <-time.After(100 * time.Millisecond):
		t.Error("Cli Close Timeout")
	}
}

// -
//
// -
func TestNewTCP_srv_send_bad_sig(t *testing.T) {
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

// -
//
// -
func TestNewTCP_srv_send_short_packet(t *testing.T) {
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
			t.Error("Recieved nil")
			return
		}
		err = packet.Validatepackets(p, r)
	case <-time.After(100 * time.Millisecond):
		t.Error("Cli Recieve timeout")
	}

	srv.Cancel()

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

// -
// Server Send a bad Sig packet
// -
func TestNewTCP_cli_recv_bad_sig(t *testing.T) {
	cx := ctx.NewContext()
	defer cx.Cancel()

	srv, cli, err := connection(cx)
	if err != nil {
		t.Error(err)
	}

	srvCloseCh := srv.link.LinkCloseCh()
	cliCloseCh := cli.link.LinkCloseCh()

	p, err := packet.TestChan32Packet()
	if err != nil {
		t.Error(err)
	}

	srv.Send(p)

	select {
	case <-srvCloseCh:
	case <-time.After(100 * time.Millisecond):
		t.Error("Srv Close Timeout")
	}

	select {
	case <-cliCloseCh:
	case <-time.After(100 * time.Millisecond):
		t.Error("Cli Close Timeout")
	}

}

// -
//
// -
func TestNewTCP_cli_recv_short_packet(t *testing.T) {
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
			t.Error("Recieved nil")
			return
		}
		err = packet.Validatepackets(p, r)
	case <-time.After(100 * time.Millisecond):
		t.Error("Cli Recieve timeout")
	}

	srv.Cancel()

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
