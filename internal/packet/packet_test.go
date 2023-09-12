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

func TestNewPacket_nil_methods(t *testing.T) {

	var p *PacketStruct

	NewPacket(0x0000)
	ReadPacketBuffer([]byte{})
	MakePacket([]byte{})

	p.ToByte()
	p.Sig()
	p.Size()
	p.Width()
	p.Counter()
	p.Ping()
	p.Pong()
	p.Router()
	p.IPv4()
	p.IPv6()
	p.Eth()
	p.Auth()
	p.ID()
	p.Packet()
	p.Raw()
	p.DebugPacket("")

}

func TestNewPacket_close_packets(t *testing.T) {
	var p *PacketStruct
	var buf []byte
	var err error

	if p, err = NewPacket(SIG_CONN_CLOSE); err != nil {
		t.Errorf("NewPacket Err:%s", err)
	}

	if buf = p.ToByte(); len(buf) != 4 {
		t.Errorf("NewPacket ToByte return %d", len(buf))
	}

	if _, err := MakePacket(buf); err != nil {
		t.Errorf("MakePacket Err:%s", err)
	}

}


func TestNewPacket_short_packets(t *testing.T) {
	var p *PacketStruct
	var buf []byte
	var err error

	if p, err = NewPacket(SIG_CONN_64_RAW,counter.MakeCounter32(0), []byte("data")); err != nil {
		t.Fatalf("NewPacket Err:%s", err)
	}

	if buf = p.ToByte(); len(buf) < 10 {
		t.Errorf("NewPacket ToByte return %d", len(buf))
	}

	if _, err := MakePacket(buf[:10]); err == nil {
		t.Errorf("MakePacket did not Err")
	}

}
