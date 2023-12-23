package tcpsrv

import (
	"testing"
	"time"

	"github.com/seanmcadam/ctx"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/netlib"
)

func TestNewTcpServer_host(t *testing.T) {

	cx := ctx.NewContext()

	config := &settings.ConnectionStruct{
		Name:  "testing",
		Proto: "tcp",
		Host:  "127.0.0.1",
		Port:  settings.ConfigPortType(uint16(netlib.GetRandomNetworkPort())),
		Auth:  "",
	}

	tcpserver, err := New(cx, config)
	_ = tcpserver
	if err != nil {
		t.Fatalf("New Error:%s", err)
	}

	time.Sleep(time.Second)

	cx.Cancel()

}

func TestNewTcpClient_test_nil_returns(t *testing.T) {

	cx := ctx.NewContext()
	defer cx.Cancel()

	config := &settings.ConnectionStruct{
		Name:  "testing",
		Proto: "tcp",
		Host:  "127.0.0.1",
		Port:  settings.ConfigPortType(uint16(netlib.GetRandomNetworkPort())),
		Auth:  "",
	}

	bigdata := make([]byte, 2048)
	for i := range bigdata {
		bigdata[i] = byte(i % 256)
	}

	//bigpacket, err := packet.NewPacket(packet.SIG_CONN_32_RAW, counter.MakeCounter32(0), bigdata)
	//packet, err := packet.NewPacket(packet.SIG_CONN_32_RAW, counter.MakeCounter32(0), []byte{})

	var tcp *TcpServerStruct

	tcp, err := new(cx, config)
	if err != nil {
		t.Fatalf("New Error:%s", err)
	}

	tcp.Cancel()

}
