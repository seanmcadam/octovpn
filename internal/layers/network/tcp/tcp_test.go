package tcp

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"testing"

	"github.com/seanmcadam/ctx"
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/internal/msg"
	"github.com/seanmcadam/octovpn/internal/packet"
)

// -
// Will the module compile
// -
func TestNewTCP_basic(t *testing.T) {
	cx := ctx.New()
	defer cx.Cancel()
}

// -
// Verify the functions handle nil input with out panic
// -
func TestNewTCP_test_nil_returns(t *testing.T) {
	cx := ctx.New()
	defer cx.Cancel()

	var ts *TcpStruct
	ts.goSend()
	ts.goRecv()
	ts.sendclose()
	ts.sendpacket(nil)
	ts.sendtestpacket(nil)
	ts.close()
	_ = ts.doneChan()

	srv, cli, err := connection(cx)
	if err != nil {
		t.Error(err)
	}

	log.Printf("CLI:%s", cli.GetInstanceName())
	log.Printf("SRV:%s", srv.GetInstanceName())

	_ = cli.GetParentMsgHandler()
	_ = cli.GetParentRecvCh()

	cli.doneChan()
	srv.doneChan()
}

// -
// Send a packet, and recieve it
// -
func TestNewTCP_send_and_compare_size(t *testing.T) {
	var wg sync.WaitGroup
	cx := ctx.New()
	defer cx.Cancel()

	testp, err := packet.TestConn32Packet()
	if err != nil {
		t.Error(err)
	}

	srv, cli, err := connection(cx)
	if err != nil {
		t.Error(err)
	}

	parent := func(tcp *TcpStruct) {
		for {
			select {
			case p := <-tcp.parentCh:
				log.Debugf("Parent Recv:%v", p)
				if packet, ok := p.(*msg.PacketStruct); ok {
					if packet.Packet.Size() != testp.Size() {
						t.Errorf("Packet sizes do not match")
					}
					wg.Done()
				}
			case <-tcp.doneChan():
				return
			}
		}
	}

	go parent(srv)
	go parent(cli)

	// Send packet from srv to cli, and have it pop out in the msgnode
	wg.Add(1)
	srv.sendCh <- testp
	wg.Add(1)
	cli.sendCh <- testp

	wg.Wait()
	srv.close()
	cli.close()

}

// -
//
// -
func TestNewTCP_cli_send_bad_sig(t *testing.T) {
	var wg sync.WaitGroup
	cx := ctx.New()
	defer cx.Cancel()

	srv, cli, err := connection(cx)
	if err != nil {
		t.Error(err)
	}

	parent := func(tcp *TcpStruct) {
		for {
			select {
			case p := <-tcp.parentCh:
				log.Debugf("Parent Recv:%v", p)
				if _, ok := p.(*msg.PacketStruct); ok {
					wg.Done()
				}
			case <-tcp.doneChan():
				return
			}
		}
	}

	go parent(srv)
	go parent(cli)

	raw := []byte{0, 0, 0, 0}
	wg.Add(1)
	cli.sendtestpacket(raw)
	wg.Add(1)
	srv.sendtestpacket(raw)

	wg.Wait()
	srv.close()
	cli.close()
}

////
////// -
//////
////// -
////func TestNewTCP_srv_send_bad_sig(t *testing.T) {
////	cx := ctx.New()
////	defer cx.Cancel()
////
////	srv, cli, err := connection(cx)
////	if err != nil {
////		t.Error(err)
////	}
////
////	srvCloseCh := srv.link.LinkCloseCh()
////	cliCloseCh := cli.link.LinkCloseCh()
////
////	raw := []byte{0, 0, 0, 0}
////	srv.sendtestpacket(raw)
////
////	select {
////	case <-cliCloseCh:
////	case <-time.After(100 * time.Millisecond):
////		t.Error("Cli Close Timeout")
////	}
////
////	select {
////	case <-srvCloseCh:
////	case <-time.After(100 * time.Millisecond):
////		t.Error("Srv Close Timeout")
////	}
////}
////
////// -
//////
////// -
////func TestNewTCP_srv_send_short_packet(t *testing.T) {
////	cx := ctx.New()
////	defer cx.Cancel()
////
////	srv, cli, err := connection(cx)
////	if err != nil {
////		t.Error(err)
////	}
////
////	srvCloseCh := srv.link.LinkCloseCh()
////	cliCloseCh := cli.link.LinkCloseCh()
////
////	p, err := packet.TestConn32Packet()
////	if err != nil {
////		t.Error(err)
////	}
////
////	var raw1, raw2 []byte
////	if raw, err := p.ToByte(); err != nil {
////		t.Error("ToByte() Err:", err)
////	} else {
////		raw1 = raw[:len(raw)-3]
////		raw2 = raw[len(raw)-3:]
////	}
////
////	srv.sendtestpacket(raw1)
////	<-time.After(time.Millisecond)
////	srv.sendtestpacket(raw2)
////
////	select {
////	case r := <-cli.RecvChan():
////		if r == nil {
////			t.Error("Recieved nil")
////			return
////		}
////		err = packet.Validatepackets(p, r)
////	case <-time.After(100 * time.Millisecond):
////		t.Error("Cli Recieve timeout")
////	}
////
////	srv.Cancel()
////
////	select {
////	case <-cliCloseCh:
////	case <-time.After(100 * time.Millisecond):
////		t.Error("Cli Close Timeout")
////	}
////
////	select {
////	case <-srvCloseCh:
////	case <-time.After(100 * time.Millisecond):
////		t.Error("Srv Close Timeout")
////	}
////}
////
////// -
////// Server Send a bad Sig packet
////// -
////func TestNewTCP_cli_recv_bad_sig(t *testing.T) {
////	cx := ctx.New()
////	defer cx.Cancel()
////
////	srv, cli, err := connection(cx)
////	if err != nil {
////		t.Error(err)
////	}
////
////	srvCloseCh := srv.link.LinkCloseCh()
////	cliCloseCh := cli.link.LinkCloseCh()
////
////	p, err := packet.TestChan32Packet()
////	if err != nil {
////		t.Error(err)
////	}
////
////	srv.Send(p)
////
////	select {
////	case <-srvCloseCh:
////	case <-time.After(100 * time.Millisecond):
////		t.Error("Srv Close Timeout")
////	}
////
////	select {
////	case <-cliCloseCh:
////	case <-time.After(100 * time.Millisecond):
////		t.Error("Cli Close Timeout")
////	}
////
////}
////
////// -
//////
////// -
////func TestNewTCP_cli_recv_short_packet(t *testing.T) {
////	cx := ctx.New()
////	defer cx.Cancel()
////
////	srv, cli, err := connection(cx)
////	if err != nil {
////		t.Error(err)
////	}
////
////	srvCloseCh := srv.link.LinkCloseCh()
////	cliCloseCh := cli.link.LinkCloseCh()
////
////	p, err := packet.TestConn32Packet()
////	if err != nil {
////		t.Error(err)
////	}
////
////	var raw1, raw2 []byte
////	if raw, err := p.ToByte(); err != nil {
////		t.Error("ToByte() Err:", err)
////	} else {
////		raw1 = raw[:len(raw)-3]
////		raw2 = raw[len(raw)-3:]
////	}
////
////	srv.sendtestpacket(raw1)
////	<-time.After(time.Millisecond)
////	srv.sendtestpacket(raw2)
////
////	select {
////	case r := <-cli.RecvChan():
////		if r == nil {
////			t.Error("Recieved nil")
////			return
////		}
////		err = packet.Validatepackets(p, r)
////	case <-time.After(100 * time.Millisecond):
////		t.Error("Cli Recieve timeout")
////	}
////
////	srv.Cancel()
////
////	select {
////	case <-cliCloseCh:
////	case <-time.After(100 * time.Millisecond):
////		t.Error("Cli Close Timeout")
////	}
////
////	select {
////	case <-srvCloseCh:
////	case <-time.After(100 * time.Millisecond):
////		t.Error("Srv Close Timeout")
////	}
////}
////
//-----------------------------------------------------------------------------

// -
//
// -
func connection(cx *ctx.Ctx) (srvconn *TcpStruct, cliconn *TcpStruct, err error) {

	srv, cli, err := createPairedConnections()
	if err != nil {
		return nil, nil, err
	}

	srvconn = new(cx, srv)
	cliconn = new(cx, cli)

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
