package netconn

import (
	"testing"

	"github.com/seanmcadam/octovpn/octolib/octolibtest"
)


func TestNewConn_tcp(t *testing.T){

	a,b, err := octolibtest.NewNetworkTCPConnPair()
	if err != nil{
		t.Fatalf("NewNetworkTCPConnPair() error:%s", err)
	}

	ncA := NewNetConn(a)
	ncB := NewNetConn(b)

	ncA.Run()
	ncB.Run()

	testb := []byte("testing")
	l, err := ncA.Write(testb)

	if err != nil{
		t.Fatalf("Write A error:%s", err)
	}
	if l != len(testb) {
		t.Fatalf("Write A length mismatch:%d %d", l, len(testb))
	}

	readb, err := ncB.Read()
	if len(testb) != len(readb){
		t.Fatalf("Read A length mismatch:%d %d", len(testb),len(readb))
	}

	ncA.Close()
	ncB.Close()

}

func TestNewConn_udp(t *testing.T){

	a,b, err := octolibtest.NewNetworkUDPConnPair()
	if err != nil{
		t.Fatalf("NewNetworkUDPConnPair() error:%s", err)
	}

	ncA := NewNetConn(a)
	ncB := NewNetConn(b)

	ncA.Run()
	ncB.Run()

	ncA.Close()
	ncB.Close()

}



