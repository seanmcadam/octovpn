package packetconn

import (
	"testing"
)

func TestNewPacket(t *testing.T) {

	cp, err := NewPacket(PACKET_TYPE_ERROR, []byte(""))
	if err != ErrChanConnPayloadLength {
		t.Error("Zero Payload did not return error")
	}

	data := []byte("data")

	cp, err = NewPacket(PACKET_TYPE_ERROR, data)
	if err != nil {
		t.Fatalf("Err:%s", err)
	}

	if cp.GetType() != PACKET_TYPE_ERROR {
		t.Error("Packet Type is not correct")
	}

	if cp.GetLength() != PacketLength(len(data)) {
		t.Error("Payload Length is not correct")
	}


}
