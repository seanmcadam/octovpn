package udpsrv

import (
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/internal/settings"
)

func TestNewUdpServer_new(t *testing.T) {

	config := &settings.NetworkStruct{
		Name:  "testing",
		Proto: "udp",
		Host:  "127.0.0.1",
		Port:  "50000",
		Auth:  "",
	}

	udpclient, err := New(config)
	if err != nil {
		t.Fatalf("New Error:%s", err)
	}

	time.Sleep(time.Second)

	udpclient.Close()

}
