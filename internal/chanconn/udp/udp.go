package udp

import (
	"net"

	"github.com/seanmcadam/octovpn/internal/chanconn"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type UdpStruct struct {
	srv     bool
	conn    *net.UDPConn
	addr    *net.UDPAddr
	sendch  chan *chanconn.ConnPacket
	recvch  chan *chanconn.ConnPacket
	Closech chan interface{}
}

func NewUDPSrv(conn *net.UDPConn) (udp *UdpStruct) {

	log.Debug("Local Addr %s", conn.LocalAddr())
	udp = &UdpStruct{
		srv:     true,
		conn:    conn,
		addr:    nil,
		sendch:  make(chan *chanconn.ConnPacket),
		recvch:  make(chan *chanconn.ConnPacket),
		Closech: make(chan interface{}),
	}

	udp.run()
	return udp
}

func NewUDPCli(conn *net.UDPConn) (udp *UdpStruct) {

	log.Debug("Local Addr %s", conn.LocalAddr())
	udp = &UdpStruct{
		srv:     false,
		conn:    conn,
		addr:    nil,
		sendch:  make(chan *chanconn.ConnPacket),
		recvch:  make(chan *chanconn.ConnPacket),
		Closech: make(chan interface{}),
	}

	udp.run()
	return udp
}

func (u *UdpStruct) endpoint() (v string) {
	if u.srv {
		v = "SRV"
	} else {
		v = "CLI"
	}
	return v
}

func (u *UdpStruct) run() {
	go u.goRecv()
	go u.goSend()
}
