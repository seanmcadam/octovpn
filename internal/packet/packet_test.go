package packet

import (
	"fmt"
	"testing"

	"github.com/seanmcadam/octovpn/internal/counter"
)

func TestNewPacket_compile(t *testing.T) {
	var p *PacketStruct
	var err error
	var c32 counter.Counter32 = 1
	var c64 counter.Counter64 = 1



	p, err = NewPacket(SIG_CONN_32_RAW, []byte(""), counter.Counter(&c32))
	if err != nil {
		t.Fatal(fmt.Sprintf("NewPacket Err:%s", err))
	}
	if p == nil {
		t.Fatal("NewPacket failed")
	}
	p, err = NewPacket(SIG_CONN_64_RAW, []byte(""), counter.Counter(&c64))
	if err != nil {
		t.Fatal(fmt.Sprintf("NewPacket Err:%s", err))
	}
	if p == nil {
		t.Fatal("NewPacket failed")
	}
}

//func TestNewPacket_layer(t *testing.T) {
//	var p *PacketStruct
//	var err error
//
//	p, err = NewPacket(Packet_ROUTE_RAW, []byte(""))
//	if p == nil || err != nil {
//		t.Fatal("NewPacket Route failed")
//	}
//	p, err = NewPacket(Packet_SITE_RAW, []byte(""))
//	if p == nil || err != nil {
//		t.Fatal("NewPacket SITE failed")
//	}
//	p, err = NewPacket(Packet_CHAN_RAW, []byte(""))
//	if p == nil || err != nil {
//		t.Fatal("NewPacket CHAN failed")
//	}
//	p, err = NewPacket(Packet_CONN_RAW, []byte(""))
//	if p == nil || err != nil {
//		t.Fatal("NewPacket CONN failed")
//	}
//}
//