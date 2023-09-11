package udpsrv

import (
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/netlib"
)

func TestNewUdpServer_new(t *testing.T) {

	cx := ctx.NewContext()

	config := &settings.ConnectionStruct{
		Name:  "testing",
		Proto: "udp",
		Host:  "127.0.0.1",
		Port:  settings.ConfigPortType(uint16(netlib.GetRandomNetworkPort())),
		Auth:  "",
	}

	_, err := New(cx, config)
	if err != nil {
		t.Fatalf("New Error:%s", err)
	}

	time.Sleep(time.Second)

	cx.Cancel()

}

func TestNewUdpServer_test_nil_returns(t *testing.T) {

	cx := ctx.NewContext()
	defer cx.Cancel()

	config := &settings.ConnectionStruct{
		Name:  "testing",
		Proto: "udp",
		Host:  "127.0.0.1",
		Port:  settings.ConfigPortType(uint16(netlib.GetRandomNetworkPort())),
		Auth:  "",
	}

	bigdata := make([]byte, 2048)
	for i := range bigdata {
		bigdata[i] = byte(i % 256)
	}

	bigpacket, err := packet.NewPacket(packet.SIG_CONN_32_RAW, counter.MakeCounter32(0), bigdata)
	packet, err := packet.NewPacket(packet.SIG_CONN_32_RAW, counter.MakeCounter32(0), []byte{})

	var u *UdpServerStruct

	u.goRun()
	u.GetLinkNoticeStateCh()
	u.GetLinkStateCh()
	u.GetUpCh()
	u.GetLinkCh()
	u.GetDownCh()
	u.GetState()
	u.Send(nil)
	u.Reset()
	u.RecvChan()

	u, err = new(cx, config)
	if err != nil {
		t.Fatalf("New Error:%s", err)
	}

	u.GetLinkNoticeStateCh()
	u.GetLinkStateCh()
	u.GetUpCh()
	u.GetLinkCh()
	u.GetDownCh()
	u.GetState()
	u.Send(nil)
	u.Send(packet)
	u.Send(bigpacket)
	u.Reset()
	u.RecvChan()

}
