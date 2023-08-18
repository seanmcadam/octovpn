package channel

import (
	"testing"

	"github.com/seanmcadam/octovpn/internal/channel/tcp"
)

func TestNewChannel_Tcp(t *testing.T) {
	var testval = []byte("test")
	serv := tcp.NewTcpServer()
	client := tcp.NewTcpClient()

	chanServ, err := NewChannel(serv)
	if err != nil {
		t.Fatalf("NewChannel Server error:%s", err)
	}

	chanClient, err := NewChannel(client)
	if err != nil {
		t.Fatalf("NewChannel Client error:%s", err)
	}

	err = chanServ.Send(testval)
	if err != nil {
		t.Fatalf("Channel Send() error:%s", err)
	}

	b, err := chanClient.Recv()
	if err != nil {
		t.Fatalf("NewChannel Client error:%s", err)
	}

	if string(b) != string(testval) {
		t.Fatalf("Send/Recv %s != %s", string(b), string(testval))

	}

}
