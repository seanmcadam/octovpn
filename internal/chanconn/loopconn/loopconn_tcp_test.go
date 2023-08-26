package loopconn

import (
	"reflect"
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestNewTcpLoop_OpenClose(t *testing.T) {

	cx := ctx.NewContext()

	l1, l2, err := NewTcpLoop(cx)

	if err != nil {
		t.Fatalf("TCP Error:%s", err)
	}

	time.Sleep(5 * time.Second)

	if !l1.Active() {
		t.Fatal("TCP L1 Active failed")
	}

	if !l2.Active() {
		t.Fatal("TCP L2 Active failed")
	}
}

func TestNewTcpLoop_SendRecv(t *testing.T) {

	cx := ctx.NewContext()

	data := []byte("data")
	l1, l2, err := NewTcpLoop(cx)

	if err != nil {
		t.Fatalf("TCP Error:%s", err)
	}

	time.Sleep(2 * time.Second)

	if !l1.Active() {
		t.Fatal("TCP L1 Active failed")
	}

	if !l2.Active() {
		t.Fatal("TCP L2 Active failed")
	}

	err = l1.Send(data)
	if err != nil {
		t.Fatalf("TCP Send Error:%s", err)
	}

	recv := <-l2.RecvChan()

	if recv == nil {
		t.Fatalf("TCP Recv Returned nil")
	}

	if !reflect.DeepEqual(recv.ToByte(), data) {
		t.Fatalf("TCP Recv Returned bad Data: '%s', '%s'", string(recv.GetPayload()), string(data))
	}

}

func TestNewTcpLoop_SendRecvReset(t *testing.T) {

	cx := ctx.NewContext()

	data := []byte("data")
	l1, l2, err := NewTcpLoop(cx)

	if err != nil {
		t.Fatalf("TCP Error:%s", err)
	}

	time.Sleep(2 * time.Second)

	if !l1.Active() {
		t.Fatal("TCP L1 Active failed")
	}

	if !l2.Active() {
		t.Fatal("TCP L2 Active failed")
	}

	err = l1.Send(data)
	if err != nil {
		t.Fatalf("TCP Send Error:%s", err)
	}

	recv := <-l2.RecvChan()

	if recv == nil {
		t.Fatalf("TCP Recv Returned nil")
	}

	//l2.Reset()

	err = l2.Send(data)
	if err != nil {
		t.Fatalf("TCP Send Error:%s", err)
	}

	recv = <-l1.RecvChan()

	if recv == nil {
		t.Fatalf("TCP Recv Returned nil")
	}

}
