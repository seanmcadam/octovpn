package loopconn

import (
	"testing"
	"time"
)

func TestNewUdpLoop_OpenClose(t *testing.T) {

	l1, l2, err := NewUdpLoop()

	if err != nil {
		t.Fatalf("UDP Error:%s", err)
	}

	time.Sleep(2 * time.Second)
	if !l1.Active() {
		t.Fatal("UDP L1 Active failed")
	}

	if !l2.Active() {
		t.Fatal("UDP L2 Active failed")
	}
}

func TestNewUdpLoop_SendRecv(t *testing.T) {

	data := []byte("data")
	srv, cli, err := NewUdpLoop()

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

	err = cli.Send(data)
	if err != nil {
		t.Fatalf("UDP Send Error:%s", err)
	}

	recv, err := srv.Recv()
	if err != nil {
		t.Fatalf("UDP Recv Error:%s", err)
	}

	if recv == nil {
		t.Fatalf("UDP Recv Returned nil")
	}

}

func TestNewTcpLoop_OpenClose(t *testing.T) {

	l1, l2, err := NewTcpLoop()

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
}

func TestNewTcpLoop_SendRecv(t *testing.T) {

	data := []byte("data")
	l1, l2, err := NewTcpLoop()

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

	recv, err := l2.Recv()
	if err != nil {
		t.Fatalf("TCP Recv Error:%s", err)
	}

	if recv == nil {
		t.Fatalf("TCP Recv Returned nil")
	}

}

func TestNewTcpLoop_SendRecvReset(t *testing.T) {

	data := []byte("data")
	l1, l2, err := NewTcpLoop()

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

	recv, err := l2.Recv()
	if err != nil {
		t.Fatalf("TCP Recv Error:%s", err)
	}

	if recv == nil {
		t.Fatalf("TCP Recv Returned nil")
	}

	//l2.Reset()

	err = l2.Send(data)
	if err != nil {
		t.Fatalf("TCP Send Error:%s", err)
	}

	recv, err = l1.Recv()
	if err != nil {
		t.Fatalf("TCP Recv Error:%s", err)
	}

	if recv == nil {
		t.Fatalf("TCP Recv Returned nil")
	}

}
