package loopconn

import (
	"reflect"
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestNewTcpLoop_OpenClose(t *testing.T) {

	cx := ctx.NewContext()

	_, _, err := NewTcpConnLoop(cx)

	if err != nil {
		t.Fatalf("TCP Error:%s", err)
	}

	time.Sleep(5 * time.Second)

}

func TestNewTcpLoop_SendRecv(t *testing.T) {

	cx := ctx.NewContext()

	var data []byte
	data = append(data, []byte("data")...)

	cp, err := packet.NewPacket(packet.SIG_CONN_32_RAW, data, counter.MakeCounter32(1))
	if err != nil {
		t.Fatalf("NewPacket Err:%s", err)
	}

	l1, l2, err := NewTcpConnLoop(cx)

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

	if !reflect.DeepEqual(recv.ToByte(), cp.ToByte()) {
		t.Fatalf("TCP Recv Returned bad Data: '%v', '%v'", recv, cp)
	}

}

func TestNewTcpLoop_SendRecvReset(t *testing.T) {

	cx := ctx.NewContext()

	data := []byte("data")
	cp, err := packet.NewPacket(packet.SIG_CONN_32_RAW, data, counter.MakeCounter32(1))
	if err != nil {
		t.Fatalf("NewPacket Err:%s", err)
	}

	l1, l2, err := NewTcpConnLoop(cx)

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
