package tcp

import (
	"net"

	"github.com/seanmcadam/octovpn/internal/chanconn"
)

type TcpStruct struct {
	conn    *net.TCPConn
	sendch  chan *chanconn.ConnPacket
	recvch  chan *chanconn.ConnPacket
	Closech chan interface{}
}

func NewTCP(conn *net.TCPConn) (tcp *TcpStruct) {

	tcp = &TcpStruct{
		conn:    conn,
		sendch:  make(chan *chanconn.ConnPacket),
		recvch:  make(chan *chanconn.ConnPacket),
		Closech: make(chan interface{}),
	}

	go tcp.goRecv()
	go tcp.goSend()

	return tcp
}
