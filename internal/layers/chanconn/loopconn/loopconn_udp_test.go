package loopconn

import (
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestNewUdpLoop_OpenClose(t *testing.T) {

	cx := ctx.NewContext()
	_, _, err := NewUdpConnLoop(cx)

	if err != nil {
		t.Fatalf("UDP Error:%s", err)
	}

	time.Sleep(10 * time.Second)

}

func TestNewUdpLoop_SendRecv(t *testing.T) {

	cx := ctx.NewContext()
	data := []byte("data")
	cp, err := packet.NewPacket(packet.SIG_CONN_32_RAW, data, counter.MakeCounter32(1))

	srv, cli, err := NewUdpConnLoop(cx)

	if err != nil {
		t.Fatalf("UDP Error:%s", err)
	}

	time.Sleep(2 * time.Second)

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

	recvbyte, err := recv.ToByte()
	cpbyte, err := cp.ToByte()

	if string(recvbyte) != string(cpbyte) {
		t.Fatalf("UDP Recv Returned bad Data: '%v', '%v'", recv, cp)
	}
}
