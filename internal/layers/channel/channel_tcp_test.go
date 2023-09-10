package channel

import (
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/internal/layers/chanconn"
	"github.com/seanmcadam/octovpn/internal/layers/conn/tcpcli"
	"github.com/seanmcadam/octovpn/internal/layers/conn/tcpsrv"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestNewChannel_Tcp(t *testing.T) {
	cx := ctx.NewContext()

	config := &settings.NetworkStruct{
		Name:  "testing",
		Proto: "tcp",
		Host:  "127.0.0.1",
		Port:  "50000",
		Auth:  "",
	}

	// Get Client and Server

	serv, err := chanconn.NewConn32(cx, config, tcpsrv.New)
	if err != nil {
		t.Fatalf("tcpsrv New err:%s", err)
	}
	if serv == nil {
		t.Fatal("serv == nil")
	}

	client, err := chanconn.NewConn32(cx, config, tcpcli.New)
	if err != nil {
		t.Fatalf("tcpcli New err:%s", err)
	}
	if client == nil {
		t.Fatal("client == nil")
	}

	// Create Channels

	chanServ, err := NewChannel32(cx, serv)
	if err != nil {
		t.Fatalf("NewChannel Server error:%s", err)
	}
	_ = chanServ

	chanClient, err := NewChannel32(cx, client)
	if err != nil {
		t.Fatalf("NewChannel Client error:%s", err)
	}
	_ = chanClient

	time.Sleep(10 * time.Second)

	//err = chanClient.Send(testval)
	//if err != nil {
	//	t.Fatalf("Channel Send() error:%s", err)
	//}

	//_ = <-chanServ.RecvChan()

	//if string(b) != string(testval) {
	//	t.Fatalf("Send/Recv %s != %s", string(b), string(testval))
	//}

	//chanServ.Close()
	//chanClient.Close()

	cx.Cancel()

}
