package udp

import (
	"fmt"
	"math/rand"
	"net"
	"testing"

	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestNewUdp(t *testing.T) {

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

	udp1 := NewUDPCli(cx, conn1)
	udp2 := NewUDPSrv(cx, conn2)

	_ = udp1
	_ = udp2

	cx.Cancel()

}

func createPairedConnections() (*net.UDPConn, *net.UDPConn, error) {
	port := getRandomPort()

	addr := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: port}

	conn1, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, nil, err
	}

	conn2, err := net.ListenUDP("udp", addr)
	if err != nil {
		conn1.Close()
		return nil, nil, err
	}

	fmt.Printf("Connections established on port %d\n", port)
	return conn1, conn2, nil
}

func handleConnection(conn *net.UDPConn, name string) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("[%s] Error reading: %v\n", name, err)
			return
		}

		receivedData := buffer[:n]
		fmt.Printf("[%s] Received from %s: %s\n", name, addr.String(), receivedData)
	}
}

func getRandomPort() int {
	return rand.Intn(60000) + 1025
}
