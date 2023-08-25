package channel

import (
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/internal/chanconn/udpcli"
	"github.com/seanmcadam/octovpn/internal/chanconn/udpsrv"
	"github.com/seanmcadam/octovpn/internal/settings"
)

func TestNewChannel_Udp(t *testing.T) {
	var testval = []byte("test")

	config := &settings.NetworkStruct{
		Name:  "testing",
		Proto: "udp",
		Host:  "127.0.0.1",
		Port:  "50000",
		Auth:  "",
	}

	// Get Client and Server

	serv, err := udpsrv.New(config)
	if err != nil {
		t.Fatalf("udpsrv New err:%s", err)
	}
	if serv == nil {
		t.Fatal("serv == nil")
	}

	client, err := udpcli.New(config)
	if err != nil {
		t.Fatalf("udpcli New err:%s", err)
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
