package chanconn

import (
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/internal/chanconn/loopconn"
	"github.com/seanmcadam/octovpn/internal/chanconn/udpcli"
	"github.com/seanmcadam/octovpn/internal/chanconn/udpsrv"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestNewUDP_SetupSrv(t *testing.T) {

	cx := ctx.NewContext()
	defer cx.Cancel()

	config := &settings.NetworkStruct{
		Name:  "testing",
		Proto: "udp",
		Host:  "127.0.0.1",
		Port:  "50000",
		Auth:  "",
	}

	NewConn(cx, config, udpsrv.New)

}

func TestNewUDP_SetupCli(t *testing.T) {

	cx := ctx.NewContext()
	defer cx.Cancel()

	config := &settings.NetworkStruct{
		Name:  "testing",
		Proto: "udp",
		Host:  "127.0.0.1",
		Port:  "50000",
		Auth:  "",
	}

	NewConn(cx, config, udpcli.New)

}

func TestNewUdp_CliSrv(t *testing.T) {

	cx := ctx.NewContext()
	defer cx.Cancel()

	srv, cli, err := loopconn.NewUdpConnLoop(cx)

	p, err := packet.Testpacket()
	if err != nil {
		t.Error(err)
	}

	cli.Send(p)
	select {
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

}
