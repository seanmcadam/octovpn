package tcpcli

import (
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/internal/settings"
)

func TestNewTcpClient_host(t *testing.T) {

config := &settings.NetworkStruct{
	Name:  "testing",
	Proto: "tcp",
	Host:  "127.0.0.1",
	Port:  "50000",
	Auth:  "",
}

	tcpclient, err := New(config)
	if err != nil {
		t.Fatalf("New Error:%s", err)
	}

	time.Sleep(time.Second)

	tcpclient.Close()


}
