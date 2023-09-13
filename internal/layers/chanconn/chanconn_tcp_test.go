package chanconn

import (
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/internal/layers/chanconn/loopconn"
	"github.com/seanmcadam/octovpn/internal/layers/conn/tcpcli"
	"github.com/seanmcadam/octovpn/internal/layers/conn/tcpsrv"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/netlib"
)

func TestNewTCP_SetupSrv(t *testing.T) {

	cx := ctx.NewContext()
	defer cx.Cancel()

	config := &settings.ConnectionStruct{
		Name:  "testing",
		Proto: "tcp",
		Host:  "127.0.0.1",
		Port:  settings.ConfigPortType(uint16(netlib.GetRandomNetworkPort())),
		Auth:  "",
	}

	_, _ = NewConn32(cx, config, tcpsrv.New)

	cx.Cancel()

}

func TestNewTCP_SetupCli(t *testing.T) {

	cx := ctx.NewContext()
	defer cx.Cancel()

	config := &settings.ConnectionStruct{
		Name:  "testing",
		Proto: "tcp",
		Host:  "127.0.0.1",
		Port:  settings.ConfigPortType(uint16(netlib.GetRandomNetworkPort())),
		Auth:  "",
	}

	NewConn32(cx, config, tcpcli.New)

}

func TestNewTCP_CliSrv(t *testing.T) {

	cx := ctx.NewContext()
	defer cx.Cancel()

	srv, cli, err := loopconn.NewTcpConnLoop(cx)

	srvUpCh := srv.Link().LinkUpCh()
	cliUpCh := cli.Link().LinkUpCh()
	srvLinkCh := srv.Link().LinkStateCh()
	cliLinkCh := cli.Link().LinkStateCh()

	<-srvLinkCh
	<-cliLinkCh
	<-srvUpCh
	<-cliUpCh

	p, err := packet.TestConn32Packet()
	if err != nil {
		t.Errorf("Testpacket() Err:%s", err)
	}

	cli.Send(p)

	select {
	case r := <-srv.RecvChan():
		if r == nil {
			t.Error("srv.RecvChan() nil")
			return
		}
		err = packet.Validatepackets(p, r)
		if r == nil {
			t.Error("Validatepacket() Recieved nil")
			return
		}

	case <-time.After(4 * time.Second):
		t.Error("Srv Recieve timeout")
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
	case <-time.After(4 * time.Second):
		t.Error("Cli Recieve timeout")
	}
}
