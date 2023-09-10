package udpcli

import (
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestNewUdpClient_new(t *testing.T) {

	cx := ctx.NewContext()
	config := &settings.NetworkStruct{
		Name:  "testing",
		Proto: "udp",
		Host:  "127.0.0.1",
		Port:  "50006",
		Auth:  "",
	}

	_, err := New(cx, config)
	if err != nil {
		t.Fatalf("New Error:%s", err)
	}

	time.Sleep(time.Second)

	cx.Cancel()

}

func TestNewUdpClient_test_nil_returns(t *testing.T) {

	cx := ctx.NewContext()
	defer cx.Cancel()

	config := &settings.NetworkStruct{
		Name:  "testing",
		Proto: "udp",
		Host:  "127.0.0.1",
		Port:  "50008",
		Auth:  "",
	}

	bigdata := make([]byte, 2048)
	for i := range bigdata {
		bigdata[i] = byte(i % 256)
	}

	bigpacket, err := packet.NewPacket(packet.SIG_CONN_32_RAW, counter.MakeCounter32(0), bigdata)
	packet, err := packet.NewPacket(packet.SIG_CONN_32_RAW, counter.MakeCounter32(0), []byte{})

	var u *UdpClientStruct

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
