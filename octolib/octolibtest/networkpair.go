package octolibtest

import (
	"fmt"
	"net"

	"github.com/seanmcadam/octovpn/octolib/netlib"
)

func NewNetworkTCPConnPair() (a, b net.Conn, err error) {
	port := netlib.GetRandomNetworkPort()
	address := fmt.Sprintf("127.0.0.1:%d", port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, nil, err
	}
	defer listener.Close()

	// Get the listener's address
	listenAddr := listener.Addr().String()

	b, err = net.Dial("tcp", listenAddr)
	if err != nil {
		return nil, nil, err
	}

	a, err = listener.Accept()
	if err != nil {
		return nil, nil, err
	}

	return a, b, nil
}

func NewNetworkUDPConnPair() (a, b net.Conn, err error) {
	port := netlib.GetRandomNetworkPort()
	address := fmt.Sprintf("127.0.0.1:%d", port)

	localAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, nil, err
	}

	conn1, err := net.DialUDP("udp", nil, localAddr)
	if err != nil {
		return nil, nil, err
	}

	conn2, err := net.DialUDP("udp", nil, localAddr)
	if err != nil {
		conn1.Close()
		return nil, nil, err
	}

	return conn1, conn2, nil
}
