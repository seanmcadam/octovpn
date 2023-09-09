package channel

import (
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/internal/layer/chanconn"
	"github.com/seanmcadam/octovpn/internal/layer/conn/udpcli"
	"github.com/seanmcadam/octovpn/internal/layer/conn/udpsrv"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestNewChannel_Udp(t *testing.T) {

	cx := ctx.NewContext()

	config := &settings.NetworkStruct{
		Name:  "testing",
		Proto: "udp",
		Host:  "127.0.0.1",
		Port:  "50000",
		Auth:  "",
	}

	// Get Client and Server

	serv, err := chanconn.NewConn32(cx, config, udpsrv.New)
	if err != nil {
		t.Fatalf("udpsrv New err:%s", err)
	}
	if serv == nil {
		t.Fatal("serv == nil")
	}

	client, err := chanconn.NewConn32(cx, config, udpcli.New)
	if err != nil {
		t.Fatalf("udpcli New err:%s", err)
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

	time.Sleep(2 * time.Second)

	//err = chanClient.Send(testval)
	//if err != nil {
	//	t.Fatalf("Channel Send() error:%s", err)
	//}

	//_ = chanServ.RecvChan()

	//if string(b) != string(testval) {
	//	t.Fatalf("Send/Recv %s != %s", string(b), string(testval))
	//}

	//chanServ.Close()
	//chanClient.Close()

	cx.Cancel()

}
