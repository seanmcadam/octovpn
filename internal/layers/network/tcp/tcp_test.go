package tcp

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"testing"

	"github.com/seanmcadam/octovpn/internal/msgbus"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
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

	parent := msgbus.MsgTarget("Parent")

	mb := msgbus.New()

	var ts *TcpStruct
	ts.goSend()
	ts.sendclose()
	ts.sendpacket(nil)
	ts.sendtestpacket(nil)
	ts.Cancel()
	_ = ts.doneChan()
	ts.goRecv()

	srv, cli, err := connection(cx, mb, parent)
	if err != nil {
		t.Error(err)
	}

	cli.doneChan()
	cli.Cancel()
	srv.doneChan()
	srv.Cancel()
}

// -
// Send a packet, and recieve it
// -
func TestNewTCP_send_and_compare_size(t *testing.T) {
	var wg sync.WaitGroup
	cx := ctx.NewContext()
	defer cx.Cancel()

	mb := msgbus.New()
	parent := msgbus.MsgTarget("Parent")

	srv, cli, err := connection(cx, mb, parent)
	if err != nil {
		t.Error(err)
	}

	p, err := packet.TestConn32Packet()
	if err != nil {
		t.Error(err)
	}

	err = mb.ReceiveHandler(parent, func(data ...interface{}) {
		if len(data) == 0 {
			t.Fatal("Retrned Data is zero length")
		}
		defer wg.Done()

		switch tp := data[0].(type) {
		case *packet.PacketStruct:
			log.Infof("Received *Packet:%s", tp.Sig().String())
		case *msgbus.StateStruct:
			log.Infof("Received *State:%s", tp.State)
		case *msgbus.NoticeStruct:
			log.Infof("Received *Notice:%s", tp.Notice)
		default:
			t.Fatalf("Handler Default reached: %T", data[0])
		}
	})
	if err != nil {
		t.Fatalf("ReceiveHandle Err:%s", err)
	}

	wg.Add(1)
	mb.Send(cli.InstanceName(), p)
	wg.Add(1)
	mb.Send(srv.InstanceName(), p)
	wg.Add(2)
	srv.Cancel()

	wg.Wait()

}

//// -
////
//// -
//func TestNewTCP_cli_send_bad_sig(t *testing.T) {
//	cx := ctx.NewContext()
//	defer cx.Cancel()
//
//	srv, cli, err := connection(cx)
//	if err != nil {
//		t.Error(err)
//	}
//
//	srvCloseCh := srv.link.LinkCloseCh()
//	cliCloseCh := cli.link.LinkCloseCh()
//
//	raw := []byte{0, 0, 0, 0}
//	cli.sendtestpacket(raw)
//
//	select {
//	case <-srvCloseCh:
//	case <-time.After(100 * time.Millisecond):
//		t.Error("Srv Close Timeout")
//	}
//
//	select {
//	case <-cliCloseCh:
//	case <-time.After(100 * time.Millisecond):
//		t.Error("Cli Close Timeout")
//	}
//}
//
//// -
////
//// -
//func TestNewTCP_srv_send_bad_sig(t *testing.T) {
//	cx := ctx.NewContext()
//	defer cx.Cancel()
//
//	srv, cli, err := connection(cx)
//	if err != nil {
//		t.Error(err)
//	}
//
//	srvCloseCh := srv.link.LinkCloseCh()
//	cliCloseCh := cli.link.LinkCloseCh()
//
//	raw := []byte{0, 0, 0, 0}
//	srv.sendtestpacket(raw)
//
//	select {
//	case <-cliCloseCh:
//	case <-time.After(100 * time.Millisecond):
//		t.Error("Cli Close Timeout")
//	}
//
//	select {
//	case <-srvCloseCh:
//	case <-time.After(100 * time.Millisecond):
//		t.Error("Srv Close Timeout")
//	}
//}
//
//// -
////
//// -
//func TestNewTCP_srv_send_short_packet(t *testing.T) {
//	cx := ctx.NewContext()
//	defer cx.Cancel()
//
//	srv, cli, err := connection(cx)
//	if err != nil {
//		t.Error(err)
//	}
//
//	srvCloseCh := srv.link.LinkCloseCh()
//	cliCloseCh := cli.link.LinkCloseCh()
//
//	p, err := packet.TestConn32Packet()
//	if err != nil {
//		t.Error(err)
//	}
//
//	var raw1, raw2 []byte
//	if raw, err := p.ToByte(); err != nil {
//		t.Error("ToByte() Err:", err)
//	} else {
//		raw1 = raw[:len(raw)-3]
//		raw2 = raw[len(raw)-3:]
//	}
//
//	srv.sendtestpacket(raw1)
//	<-time.After(time.Millisecond)
//	srv.sendtestpacket(raw2)
//
//	select {
//	case r := <-cli.RecvChan():
//		if r == nil {
//			t.Error("Recieved nil")
//			return
//		}
//		err = packet.Validatepackets(p, r)
//	case <-time.After(100 * time.Millisecond):
//		t.Error("Cli Recieve timeout")
//	}
//
//	srv.Cancel()
//
//	select {
//	case <-cliCloseCh:
//	case <-time.After(100 * time.Millisecond):
//		t.Error("Cli Close Timeout")
//	}
//
//	select {
//	case <-srvCloseCh:
//	case <-time.After(100 * time.Millisecond):
//		t.Error("Srv Close Timeout")
//	}
//}
//
//// -
//// Server Send a bad Sig packet
//// -
//func TestNewTCP_cli_recv_bad_sig(t *testing.T) {
//	cx := ctx.NewContext()
//	defer cx.Cancel()
//
//	srv, cli, err := connection(cx)
//	if err != nil {
//		t.Error(err)
//	}
//
//	srvCloseCh := srv.link.LinkCloseCh()
//	cliCloseCh := cli.link.LinkCloseCh()
//
//	p, err := packet.TestChan32Packet()
//	if err != nil {
//		t.Error(err)
//	}
//
//	srv.Send(p)
//
//	select {
//	case <-srvCloseCh:
//	case <-time.After(100 * time.Millisecond):
//		t.Error("Srv Close Timeout")
//	}
//
//	select {
//	case <-cliCloseCh:
//	case <-time.After(100 * time.Millisecond):
//		t.Error("Cli Close Timeout")
//	}
//
//}
//
//// -
////
//// -
//func TestNewTCP_cli_recv_short_packet(t *testing.T) {
//	cx := ctx.NewContext()
//	defer cx.Cancel()
//
//	srv, cli, err := connection(cx)
//	if err != nil {
//		t.Error(err)
//	}
//
//	srvCloseCh := srv.link.LinkCloseCh()
//	cliCloseCh := cli.link.LinkCloseCh()
//
//	p, err := packet.TestConn32Packet()
//	if err != nil {
//		t.Error(err)
//	}
//
//	var raw1, raw2 []byte
//	if raw, err := p.ToByte(); err != nil {
//		t.Error("ToByte() Err:", err)
//	} else {
//		raw1 = raw[:len(raw)-3]
//		raw2 = raw[len(raw)-3:]
//	}
//
//	srv.sendtestpacket(raw1)
//	<-time.After(time.Millisecond)
//	srv.sendtestpacket(raw2)
//
//	select {
//	case r := <-cli.RecvChan():
//		if r == nil {
//			t.Error("Recieved nil")
//			return
//		}
//		err = packet.Validatepackets(p, r)
//	case <-time.After(100 * time.Millisecond):
//		t.Error("Cli Recieve timeout")
//	}
//
//	srv.Cancel()
//
//	select {
//	case <-cliCloseCh:
//	case <-time.After(100 * time.Millisecond):
//		t.Error("Cli Close Timeout")
//	}
//
//	select {
//	case <-srvCloseCh:
//	case <-time.After(100 * time.Millisecond):
//		t.Error("Srv Close Timeout")
//	}
//}
//
//-----------------------------------------------------------------------------

// -
//
// -
func connection(cx *ctx.Ctx, mb *msgbus.MsgBus, p msgbus.MsgTarget) (srvconn *TcpStruct, cliconn *TcpStruct, err error) {

	srv, cli, err := createPairedConnections()
	if err != nil {
		return nil, nil, err
	}

	srvconn = NewTCP(cx, mb, p, srv)
	cliconn = NewTCP(cx, mb, p, cli)

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
