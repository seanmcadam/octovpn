package packetchan

import (
	"testing"
)

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

	if cp.GetType() != CHAN_TYPE_DATA {
		t.Error("Packet Type is not correct")
	}

	if cp.GetLength() != ChanLength(len(data)) {
		t.Error("Payload Length is not correct")
	}

	_, err = MakePacket(cp.ToByte())
	if err != nil {
		t.Errorf("MakePacket DATA Err:%s", err)
	}

	cp.ConvertDataToAck()
	_, err = MakePacket(cp.ToByte())
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
