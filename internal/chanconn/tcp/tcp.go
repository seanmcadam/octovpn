package tcp

import (
	"net"
	"time"

	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
	"github.com/seanmcadam/octovpn/octolib/pinger"
)

const TcpPingFreq = 1 * time.Second
const TcpPingTimeout = 2 * time.Second

type TcpStruct struct {
	conn    *net.TCPConn
	pinger  *pinger.Pinger64Struct
	sendch  chan *packetconn.ConnPacket
	recvch  chan *packetconn.ConnPacket
	Closech chan interface{}
}

func NewTCP(conn *net.TCPConn) (tcp *TcpStruct) {

	closech := make(chan interface{})
	tcp = &TcpStruct{
		conn:    conn,
		pinger:  pinger.NewPinger64(TcpPingFreq, TcpPingTimeout, closech),
		sendch:  make(chan *packetconn.ConnPacket),
		recvch:  make(chan *packetconn.ConnPacket),
		Closech: closech,
	}

	tcp.run()
	tcp.pinger.TurnOn() // Tcp Ping is always on, if the the connection drops the object is terminated

	return tcp
}

func (t *TcpStruct) run() {
	go t.goRecv()
	go t.goSend()
}
