package conn

import (
	"testing"
	"time"

	"github.com/seanmcadam/ctx"
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/internal/layers/conn/tcpcli"
	"github.com/seanmcadam/octovpn/internal/layers/conn/tcpsrv"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/netlib"
)

func Test_TCPLoop(t *testing.T) {
	cx := ctx.New()
	config := &settings.ConnectionStruct{
		Name:  "Server Loopback TCP",
		Proto: "tcp",
		Host:  "127.0.0.1",
		Port:  settings.ConfigPortType(uint16(netlib.GetRandomNetworkPort())),
		Auth:  "",
	}

	srv, err := tcpsrv.New(cx, config)
	if err != nil {
		t.Fatalf("tcpsrv New() Err:%s", err)
	}

	<-srv.Link().LinkListenCh()
	srvupch := srv.Link().LinkStateCh()

	cli, err := tcpcli.New(cx, config)
	if err != nil {
		t.Fatalf("tcpsrv New() Err:%s", err)
	}

	cliupch := cli.Link().LinkStateCh()

	<-time.After(time.Second)

	state := <-srvupch
	log.Debugf("Srv state:%s", state)
	state = <-cliupch
	log.Debugf("Cli state:%s", state)

	cx.Cancel()
}
