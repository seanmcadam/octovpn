package udp

import (
	"net"
	"time"

	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
	"github.com/seanmcadam/octovpn/octolib/pinger"
)

type UdpStruct struct {
	srv     bool
	conn    *net.UDPConn
	addr    *net.UDPAddr
	pinger  *pinger.Pinger64Struct
	sendch  chan *packetconn.ConnPacket
	recvch  chan *packetconn.ConnPacket
	Closech chan interface{}
}

func NewUDPSrv(conn *net.UDPConn) (udp *UdpStruct) {

	log.Debug("Local Addr %s", conn.LocalAddr())

	closech := make(chan interface{})
	udp = &UdpStruct{
		srv:     true,
		conn:    conn,
		addr:    nil,
		pinger:  pinger.NewPinger64(time.Second, 5*time.Second, closech),
		sendch:  make(chan *packetconn.ConnPacket),
		recvch:  make(chan *packetconn.ConnPacket),
		Closech: closech,
	}

	udp.run()
	return udp
}

func NewUDPCli(conn *net.UDPConn) (udp *UdpStruct) {

	log.Debug("Local Addr %s", conn.LocalAddr())

	closech := make(chan interface{})
	udp = &UdpStruct{
		srv:     false,
		conn:    conn,
		addr:    nil,
		pinger:  pinger.NewPinger64(time.Second, 5*time.Second, closech),
		sendch:  make(chan *packetconn.ConnPacket),
		recvch:  make(chan *packetconn.ConnPacket),
		Closech: closech,
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
