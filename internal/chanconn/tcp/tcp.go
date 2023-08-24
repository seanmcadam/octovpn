package tcp

import (
	"net"
)

type TcpStruct struct {
	conn    *net.TCPConn
	sendch  chan []byte
	recvch  chan []byte
	Closech chan interface{}
}

func NewTCP(conn *net.TCPConn) (tcp *TcpStruct){

	tcp = &TcpStruct{
		conn:    conn,
		sendch:  make(chan []byte),
		recvch:  make(chan []byte),
		Closech: make(chan interface{}),
	}

	go tcp.goRecv()
	go tcp.goSend()

	return tcp
}
