package loopconn

import (
	"reflect"
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestNewUdpLoop_OpenClose(t *testing.T) {

	cx := ctx.NewContext()
	l1, l2, err := NewUdpConnLoop(cx)

	if err != nil {
		t.Fatalf("UDP Error:%s", err)
	}

	time.Sleep(10 * time.Second)

	if !l1.Active() {
		t.Fatal("UDP L1 Active failed")
	}

	if !l2.Active() {
		t.Fatal("UDP L2 Active failed")
	}
}

func TestNewUdpLoop_SendRecv(t *testing.T) {

	cx := ctx.NewContext()
	data := []byte("data")
	cp, err := packet.NewPacket(packet.SIG_CONN_32_RAW, data)

	srv, cli, err := NewUdpConnLoop(cx)

	if err != nil {
		t.Fatalf("UDP Error:%s", err)
	}

	time.Sleep(2 * time.Second)
	if !srv.Active() {
		t.Fatal("UDP L1 Active failed")
	}

	if !cli.Active() {
		t.Fatal("UDP L2 Active failed")
	}

	err = cli.Send(cp)
	if err != nil {
		t.Fatalf("UDP Send Error:%s", err)
	}

	recv := <-srv.RecvChan()
	if err != nil {
		t.Fatalf("UDP Recv Error:%s", err)
	}

	if recv == nil {
		t.Fatalf("UDP Recv Returned nil")
	}

	if !reflect.DeepEqual(recv.ToByte(), cp.ToByte()) {
		t.Fatalf("UDP Recv Returned bad Data: '%v', '%v'", recv, cp)
	}
}
