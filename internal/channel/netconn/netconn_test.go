package netconn

import (
	"crypto/rand"
	"fmt"
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

	err = push_data_through_test(t,ncA,ncB,1)
	if err != nil{
		t.Fatalf("TCP push_data_through_test err: %s", err)
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

	err = push_data_through_test(t,ncA,ncB,1)
	if err != nil{
		t.Fatalf("UDP push_data_through_test err: %s", err)
	}

	ncA.Close()
	ncB.Close()

}



func push_data_through_test( t *testing.T, a, b *NetConnStruct, size int)(err error){

	testb, err := generateRandomBytes(size)
	if err != nil{
		return err
	}

	l, err := a.Write(testb)
	if err != nil{
		return err	
	}
	if l != len(testb) {
		return fmt.Errorf("Write() data lengths do not match %d:%d", l, len(testb))
	}

	readb, err := b.Read()
	if len(testb) != len(readb){
		return fmt.Errorf("Length of data does not match %d:%d", len(testb), len(readb))
	}

	return nil
}


func generateRandomBytes(length int) ([]byte, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	return randomBytes, nil
}
