package channel

import (
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/internal/chanconn/udpcli"
	"github.com/seanmcadam/octovpn/internal/chanconn/udpsrv"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestNewChannel_Udp(t *testing.T) {
	//var testval = []byte("test")

	cx := ctx.NewContext()

	config := &settings.NetworkStruct{
		Name:  "testing",
		Proto: "udp",
		Host:  "127.0.0.1",
		Port:  "50000",
		Auth:  "",
	}

	// Get Client and Server

	serv, err := udpsrv.New(cx, config)
	if err != nil {
		t.Fatalf("udpsrv New err:%s", err)
	}
	if serv == nil {
		t.Fatal("serv == nil")
	}

	client, err := udpcli.New(cx, config)
	if err != nil {
		t.Fatalf("udpcli New err:%s", err)
	}
	if client == nil {
		t.Fatal("client == nil")
	}

	// Create Channels

	chanServ, err := NewChannel(cx, serv)
	if err != nil {
		t.Fatalf("NewChannel Server error:%s", err)
	}
	_ = chanServ

	chanClient, err := NewChannel(cx, client)
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
