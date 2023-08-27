package chanconn

import (
	"testing"

	"github.com/seanmcadam/octovpn/internal/chanconn/tcpcli"
	"github.com/seanmcadam/octovpn/internal/chanconn/tcpsrv"
	"github.com/seanmcadam/octovpn/internal/chanconn/udpcli"
	"github.com/seanmcadam/octovpn/internal/chanconn/udpsrv"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestNewTCPSrv(t *testing.T) {

	cx := ctx.NewContext()

	config := &settings.NetworkStruct{
		Name:  "testing",
		Proto: "tcp",
		Host:  "127.0.0.1",
		Port:  "50000",
		Auth:  "",
	}

	New(cx, config, tcpsrv.New)

	cx.Cancel()
}

func TestNewTCPCli(t *testing.T) {

	cx := ctx.NewContext()

	config := &settings.NetworkStruct{
		Name:  "testing",
		Proto: "tcp",
		Host:  "127.0.0.1",
		Port:  "50000",
		Auth:  "",
	}

	New(cx, config, tcpcli.New)

	cx.Cancel()
}

func TestNewUDPSrv(t *testing.T) {

	cx := ctx.NewContext()

	config := &settings.NetworkStruct{
		Name:  "testing",
		Proto: "udp",
		Host:  "127.0.0.1",
		Port:  "50000",
		Auth:  "",
	}

	New(cx, config, udpsrv.New)

	cx.Cancel()

}

func TestNewUDPCli(t *testing.T) {

	cx := ctx.NewContext()

	config := &settings.NetworkStruct{
		Name:  "testing",
		Proto: "udp",
		Host:  "127.0.0.1",
		Port:  "50000",
		Auth:  "",
	}

	New(cx, config, udpcli.New)

	cx.Cancel()
}
