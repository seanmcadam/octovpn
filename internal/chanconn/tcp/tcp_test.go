package tcp

import (
	"fmt"
	"math/rand"
	"net"
	"testing"

	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestNewTCP(t *testing.T) {

	cx := ctx.NewContext()

	conn1, conn2, err := createPairedConnections()
	if err != nil {
		t.Errorf("Error:%s", err)
		return
	}

	defer conn1.Close()
	defer conn2.Close()

	go handleConnection(conn1, "Connection 1")
	go handleConnection(conn2, "Connection 2")

	tcp1 := NewTCP(cx, conn1)
	tcp2 := NewTCP(cx, conn2)

	_ = tcp1
	_ = tcp2

	cx.Cancel()

}

func createPairedConnections() (*net.TCPConn, *net.TCPConn, error) {
	port := getRandomPort()

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: port})
	if err != nil {
		return nil, nil, err
	}

	conn1, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: port})
	if err != nil {
		listener.Close()
		return nil, nil, err
	}

	conn2, err := listener.AcceptTCP()
	if err != nil {
		listener.Close()
		conn1.Close()
		return nil, nil, err
	}

	fmt.Printf("Connections established on port %d\n", port)
	return conn1, conn2, nil
}

func handleConnection(conn *net.TCPConn, name string) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("[%s] Error reading: %v\n", name, err)
			return
		}

		receivedData := buffer[:n]
		fmt.Printf("[%s] Received: %s\n", name, receivedData)
	}
}

func getRandomPort() int {
	return rand.Intn(60000) + 1025
}