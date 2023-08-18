package tcp

import (
	"testing"

	"github.com/seanmcadam/octovpn/octotypes"
)

func TestNewTcpClient_host(t *testing.T) {

	randomport := octotypes.GetRandomNetworkPort()
	var err error

	_, err = NewTcpClient("127.0.0.1", randomport)
	if err != nil {
		t.Fatalf("Get 127.0.0.1 failed")
	}

	_, err = NewTcpClient("localhost", randomport)
	if err != nil {
		t.Fatalf("Get localhost failed")
	}

	_, err = NewTcpClient("google.com", randomport)
	if err != nil {
		t.Fatalf("Get google.com failed")
	}

}
