package packetchan

import (
	"testing"

	"github.com/seanmcadam/octovpn/octolib/log"
)

func TestCompilerCheck(t *testing.T) {
	_, _ = NewPacket(CHAN_TYPE_ERROR, nil)
}

func TestNewPacket(t *testing.T) {

	cp, err := NewPacket(CHAN_TYPE_ERROR, []byte(""))
	if err != nil {
		t.Errorf("Zero Payload Err:%s", err)
	}

	data := []byte("data")

	cp, err = NewPacket(CHAN_TYPE_DATA, data)
	if err != nil {
		t.Fatalf("Err:%s", err)
	}

	log.Infof("Size ConnPacket:%d", cp.GetSize())

	if cp.GetType() != CHAN_TYPE_DATA {
		t.Error("Packet Type is not correct")
	}

	if cp.GetPayloadLength() != ChanLength(len(data)) {
		t.Error("Payload Length is not correct")
	}

	_, err = MakePacket(cp.ToByte())
	if err != nil {
		t.Errorf("MakePacket DATA Err:%s", err)
	}

	ackcp := cp.CopyDataToAck()
	_, err = MakePacket(ackcp.ToByte())
	if err != nil {
		t.Errorf("MakePacket ACK Err:%s", err)
	}

}

//func TestNewPacketTcp(t *testing.T) {
//
//	srv, cli, err := loopconn.NewTcpLoop()
//	if err != nil {
//		t.Fatalf("Error:%s", err)
//	}
//
//	cssrv, err := NewChannel(srv)
//	cscli, err := NewChannel(cli)
//
//	time.Sleep(3 * time.Second)
//
//	cssrv.Close()
//	cscli.Close()
//
//}
//
