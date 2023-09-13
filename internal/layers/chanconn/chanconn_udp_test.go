package chanconn

import (
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/internal/layers/chanconn/loopconn"
	"github.com/seanmcadam/octovpn/internal/layers/conn/udpcli"
	"github.com/seanmcadam/octovpn/internal/layers/conn/udpsrv"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/netlib"
)

func TestNewUDP_SetupSrv(t *testing.T) {

	cx := ctx.NewContext()
	defer cx.Cancel()

	config := &settings.ConnectionStruct{
		Name:  "testing",
		Proto: "udp",
		Host:  "127.0.0.1",
		Port:  settings.ConfigPortType(uint16(netlib.GetRandomNetworkPort())),
		Auth:  "",
	}

	NewConn32(cx, config, udpsrv.New)

}

func TestNewUDP_SetupCli(t *testing.T) {

	cx := ctx.NewContext()
	defer cx.Cancel()

	config := &settings.ConnectionStruct{
		Name:  "testing",
		Proto: "udp",
		Host:  "127.0.0.1",
		Port:  settings.ConfigPortType(uint16(netlib.GetRandomNetworkPort())),
		Auth:  "",
	}

	NewConn32(cx, config, udpcli.New)

}

func TestNewUdp_CliSrv(t *testing.T) {

	cx := ctx.NewContext()
	defer cx.Cancel()

	srv, cli, err := loopconn.NewUdpConnLoop(cx)

	srvUpCh := srv.Link().LinkUpCh()
	cliUpCh := cli.Link().LinkUpCh()

	p, err := packet.TestConn32Packet()
	if err != nil {
		t.Error(err)
	}

	cli.Send(p)
	select {
	case <-srv.Link().LinkUpCh():
	case r := <-srv.RecvChan():
		if r == nil {
			t.Error("Recieved nil")
			return
		}
		err = packet.Validatepackets(p, r)
		if r == nil {
			t.Error("Recieved nil")
			return
		}
	case <-time.After(2 * time.Second):
		t.Error("Recieve timeout")
		return
	}

	srv.Send(p)
	select {
	case r := <-cli.RecvChan():
		if r == nil {
			t.Error("Recieved nil")
			return
		}
		err = packet.Validatepackets(p, r)
	case <-time.After(2 * time.Second):
		t.Error("Recieve timeout")
	}

	<-srvUpCh
	<-cliUpCh
}
