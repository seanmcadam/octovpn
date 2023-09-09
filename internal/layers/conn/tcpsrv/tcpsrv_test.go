package tcpsrv

import (
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestNewTcpServer_host(t *testing.T) {

	cx := ctx.NewContext()

	config := &settings.NetworkStruct{
		Name:  "testing",
		Proto: "tcp",
		Host:  "127.0.0.1",
		Port:  "50000",
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
