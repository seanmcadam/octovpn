package packetconn

import (
	"testing"

	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/packet"
)

func TestCompilerCheck(t *testing.T) {
	_, _ = NewPacket(packet.CONN_TYPE_RAW, nil)
}

func TestNewPacket(t *testing.T) {

	nodata := []byte("")
	cp, err := NewPacket(packet.CONN_TYPE_RAW, nodata)
	if err != errors.ErrConnPayloadLength {
		t.Error("Zero Payload did not return error")
	}

	data := []byte("data")

	cp, err = NewPacket(packet.CONN_TYPE_RAW, data)
	if err != nil {
		t.Fatalf("Err:%s", err)
	}

	log.Infof("Size ConnPacket:%d", cp.Size())

	if cp.Type() != packet.CONN_TYPE_RAW {
		t.Error("Packet Type is not correct")
	}

	if cp.PayloadSize() != packet.PacketPayloadSize(len(data)) {
		t.Error("Payload Length is not correct")
	}

	retdata := cp.ToByte()
	if len(retdata) != Overhead+len(data) {
		t.Error("Length is not correct")
	}

}
