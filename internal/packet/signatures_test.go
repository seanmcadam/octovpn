package packet

import "testing"

func TestAll(t *testing.T) {

	tests := []struct {
		name        string
		data        PacketSigType
		V1          bool
		RouterLayer bool
		SiteLayer   bool
		ChanLayer   bool
		ConnLayer   bool
		Parent      bool
		Raw         bool
		Auth        bool
		Ack         bool
		Nak         bool
		Ping32      bool
		Pong32      bool
		Ping64      bool
		Pong64      bool
		Error       bool
		ID          bool
		IPV4        bool
		IPV6        bool
		Eth         bool
		Router      bool
	}{
		{"All Zeros", PacketSigType(0x0000), false, false, false, false, false, true, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
	}

	for _, x := range tests {
		if x.data.V1() != x.V1 {
			t.Fatal()
		}
		if x.data.RouterLayer() != x.RouterLayer {
			t.Fatal()
		}
		if x.data.SiteLayer() != x.SiteLayer {
			t.Fatal()
		}
		if x.data.ChanLayer() != x.ChanLayer {
			t.Fatal()
		}
		if x.data.ConnLayer() != x.ConnLayer {
			t.Fatal()
		}
		if x.data.Parent() != x.Parent {
			t.Fatal()
		}
		if x.data.Raw() != x.Raw {
			t.Fatal()
		}
		if x.data.Auth() != x.Auth {
			t.Fatal()
		}
		if x.data.Ack() != x.Ack {
			t.Fatal()
		}
		if x.data.Nak() != x.Nak {
			t.Fatal()
		}
		if x.data.Ping32() != x.Ping32 {
			t.Fatal()
		}
		if x.data.Pong32() != x.Pong32 {
			t.Fatal()
		}
		if x.data.Ping64() != x.Ping64 {
			t.Fatal()
		}
		if x.data.Pong64() != x.Pong64 {
			t.Fatal()
		}
		if x.data.Error() != x.Error {
			t.Fatal()
		}
		if x.data.ID() != x.ID {
			t.Fatal()
		}
		if x.data.IPV4() != x.IPV4 {
			t.Fatal()
		}
		if x.data.IPV6() != x.IPV6 {
			t.Fatal()
		}
		if x.data.Eth() != x.Eth {
			t.Fatal()
		}
		if x.data.Router() != x.Router {
			t.Fatal()
		}
	}
}
func TestV1(t *testing.T) {
	var p1 PacketSigType = 0x1000
	var p2 PacketSigType = 0x2000
	if !p1.V1() {
		t.Fatal("V1")
	}
	if p2.V1() {
		t.Fatal("V2")
	}
}

func TestLayer(t *testing.T) {
	var p1 PacketSigType = Packet_ROUT
	var p2 PacketSigType = Packet_SITE
	if !p1.V1() {
		t.Fatal("V1")
	}
	if p2.V1() {
		t.Fatal("V2")
	}
}

//func (p PacketSigType) RouterLayer() bool {
//	return (p&Packet_LAYER)^Packet_ROUT == 0
//}
//func (p PacketSigType) SiteLayer() bool {
//	return (p&Packet_LAYER)^Packet_SITE == 0
//}
//func (p PacketSigType) ChanLayer() bool {
//	return (p&Packet_LAYER)^Packet_CHAN == 0
//}
//func (p PacketSigType) ConnLayer() bool {
//	return (p&Packet_LAYER)^Packet_CONN == 0
//}
//
//func (p PacketSigType) Parent() bool {
//	return p&Packet_TYPE == 0
//}
//func (p PacketSigType) Raw() bool {
//	return (p&Packet_TYPE)^Packet_RAW == 0
//}
//func (p PacketSigType) Auth() bool {
//	return (p&Packet_TYPE)^Packet_AUTH == 0
//}
//func (p PacketSigType) Ack() bool {
//	return (p&Packet_TYPE)^Packet_ACK == 0
//}
//func (p PacketSigType) Nak() bool {
//	return (p&Packet_TYPE)^Packet_NAK == 0
//}
//func (p PacketSigType) Ping32() bool {
//	return (p&Packet_TYPE)^Packet_PING32 == 0
//}
//func (p PacketSigType) Pong32() bool {
//	return (p&Packet_TYPE)^Packet_PONG32 == 0
//}
//func (p PacketSigType) Ping64() bool {
//	return (p&Packet_TYPE)^Packet_PING64 == 0
//}
//func (p PacketSigType) Pong64() bool {
//	return (p&Packet_TYPE)^Packet_PONG64 == 0
//}
//func (p PacketSigType) Error() bool {
//	return (p&Packet_TYPE)^Packet_ERROR == 0
//}
//func (p PacketSigType) ID() bool {
//	return (p&Packet_TYPE)^Packet_ID == 0
//}
//func (p PacketSigType) IPV4() bool {
//	return (p&Packet_TYPE)^Packet_IPV4 == 0
//}
//func (p PacketSigType) IPV6() bool {
//	return (p&Packet_TYPE)^Packet_IPV6 == 0
//}
//func (p PacketSigType) Eth() bool {
//	return (p&Packet_TYPE)^Packet_ETH == 0
//}
//func (p PacketSigType) Router() bool {
//	return (p&Packet_TYPE)^Packet_ROUTER == 0
//}
//
