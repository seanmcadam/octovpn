package packetchan

import (
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/internal/chanconn/loopconn"
)

func TestNewPacket(t *testing.T) {

	cp, err := NewPacket(CHAN_TYPE_ERROR, []byte(""))
	if err != ErrChanPayloadLength {
		t.Error("Zero Payload did not return error")
	}

	data := []byte("data")

	cp, err = NewPacket(CHAN_TYPE_ERROR, data)
	if err != nil {
		t.Fatalf("Err:%s", err)
	}

	if cp.GetType() != CHAN_TYPE_ERROR {
		t.Error("Packet Type is not correct")
	}

	if cp.GetLength() != ChanLength(len(data)) {
		t.Error("Payload Length is not correct")
	}

}

func TestNewPacketTcp(t *testing.T) {

	srv, cli, err := loopconn.NewTcpLoop()
	if err != nil {
		t.Fatalf("Error:%s", err)
	}

	cssrv, err := NewChannel(srv)
	cscli, err := NewChannel(cli)

	time.Sleep(3 * time.Second)

	cssrv.Close()
	cscli.Close()

}
