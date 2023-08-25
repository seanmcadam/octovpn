package packetconn

import (
	"testing"

	"github.com/seanmcadam/octovpn/octolib/errors"
)

func TestCompilerCheck(t *testing.T) {
	_, _ = NewPacket(PACKET_TYPE_ERROR, nil)
}

func TestNewPacket(t *testing.T) {

	nodata := []byte("")
	cp, err := NewPacket(PACKET_TYPE_ERROR, nodata)
	if err != errors.ErrChanConnPayloadLength {
		t.Error("Zero Payload did not return error")
	}

	data := []byte("data")

	cp, err = NewPacket(PACKET_TYPE_LOOP, data)
	if err != nil {
		t.Fatalf("Err:%s", err)
	}

	if cp.GetType() != PACKET_TYPE_LOOP {
		t.Error("Packet Type is not correct")
	}

	if cp.GetLength() != PacketLength(len(data)) {
		t.Error("Payload Length is not correct")
	}

	retdata := cp.ToByte()
	if len(retdata) != PacketOverhead+len(data) {
		t.Error("Length is not correct")
	}

}
