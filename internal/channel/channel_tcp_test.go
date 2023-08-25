package channel

import (
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/internal/chanconn/tcpcli"
	"github.com/seanmcadam/octovpn/internal/chanconn/tcpsrv"
	"github.com/seanmcadam/octovpn/internal/settings"
)

func TestNewChannel_Tcp(t *testing.T) {
	var testval = []byte("test")

	config := &settings.NetworkStruct{
		Name:  "testing",
		Proto: "tcp",
		Host:  "127.0.0.1",
		Port:  "50000",
		Auth:  "",
	}

	// Get Client and Server

	serv, err := tcpsrv.New(config)
	if err != nil {
		t.Fatalf("tcpsrv New err:%s", err)
	}
	if serv == nil {
		t.Fatal("serv == nil")
	}

	client, err := tcpcli.New(config)
	if err != nil {
		t.Fatalf("tcpcli New err:%s", err)
	}
	if client == nil {
		t.Fatal("client == nil")
	}

	// Create Channels

	chanServ, err := NewChannel(serv)
	if err != nil {
		t.Fatalf("NewChannel Server error:%s", err)
	}

	chanClient, err := NewChannel(client)
	if err != nil {
		t.Fatalf("NewChannel Client error:%s", err)
	}

	time.Sleep(2 * time.Second)

	err = chanClient.Send(testval)
	if err != nil {
		t.Fatalf("Channel Send() error:%s", err)
	}

	b, err := chanServ.Recv()
	if err != nil {
		t.Fatalf("NewChannel Client error:%s", err)
	}

	if string(b) != string(testval) {
		t.Fatalf("Send/Recv %s != %s", string(b), string(testval))

	}

	chanServ.Close()
	chanClient.Close()

}
