package loopchan

import (
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func TestNewTcpLoop_compile(t *testing.T) {
	cx := ctx.NewContext()

	_, _, err := NewTcpChanLoop(cx)
	if err != nil {
		t.Fatalf("NewTcpChanLoop Err:%s", err)
	}

	cx.Cancel()
}

func TestNewTcpLoop_OpenClose(t *testing.T) {

	cx := ctx.NewContext()

	_, _, err := NewTcpChanLoop(cx)

	if err != nil {
		t.Fatalf("TCP Error:%s", err)
	}

	time.Sleep(5 * time.Second)

	cx.Cancel()
}

func TestNewTcpLoop_SendRecv(t *testing.T) {

	cx := ctx.NewContext()

	data := []byte("data")
	cp, err := packet.NewPacket(packet.SIG_CONN_32_RAW, data, counter.MakeCounter32(33))
	if err != nil {
		t.Fatalf("NewPacket Err:%s", err)
	}

	l1, l2, err := NewTcpChanLoop(cx)

	if err != nil {
		t.Fatalf("TCP Error:%s", err)
	}

	select {
	case <-l1.Link().LinkUpCh():
		log.Debug("L1 Up")
	case <-l2.Link().LinkUpCh():
		log.Debug("L2 Up")
	case <-time.After(2 * time.Second):
		t.Error("Timeout waiting for channels to come up")
	}

	if l1.Link().IsDown()  {
		t.Fatalf("TCP Send Channel 1 down")
	}

	if l2.Link().IsDown()  {
		t.Fatalf("TCP Send Channel 2 down")
	}

	err = l1.Send(cp)
	if err != nil {
		t.Fatalf("TCP Send Error:%s", err)
	}

	recv := <-l2.RecvChan()

	if recv == nil {
		t.Fatalf("TCP Recv Returned nil")
	}

	recvbyte, err := recv.ToByte()
	cpbyte, err := cp.ToByte()

	if string(recvbyte) != string(cpbyte) {
		t.Fatalf("TCP Recv Returned bad Data: '%v', '%v'", recv, cp)
	}

}

func TestNewTcpLoop_SendRecvReset(t *testing.T) {

	cx := ctx.NewContext()

	data := []byte("data")
	cp, err := packet.NewPacket(packet.SIG_CONN_32_RAW, data, counter.MakeCounter32(25))
	if err != nil {
		t.Fatalf("NewPacket Err:%s", err)
	}

	l1, l2, err := NewTcpChanLoop(cx)

	if err != nil {
		t.Fatalf("TCP Error:%s", err)
	}

	time.Sleep(2 * time.Second)

	err = l1.Send(cp)
	if err != nil {
		t.Fatalf("TCP Send Error:%s", err)
	}

	recv := <-l2.RecvChan()

	if recv == nil {
		t.Fatalf("TCP Recv Returned nil")
	}

	//l2.Reset()

	err = l2.Send(cp)
	if err != nil {
		t.Fatalf("TCP Send Error:%s", err)
	}

	recv = <-l1.RecvChan()

	if recv == nil {
		t.Fatalf("TCP Recv Returned nil")
	}

}
