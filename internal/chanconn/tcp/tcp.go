package tcp

import (
	"net"

	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

type TcpStruct struct {
	cx     *ctx.Ctx
	conn   *net.TCPConn
	sendch chan *packetconn.ConnPacket
	recvch chan *packetconn.ConnPacket
}

func NewTCP(ctx *ctx.Ctx, conn *net.TCPConn) (tcp *TcpStruct) {

	tcp = &TcpStruct{
		cx:     ctx,
		conn:   conn,
		sendch: make(chan *packetconn.ConnPacket),
		recvch: make(chan *packetconn.ConnPacket),
	}

	tcp.run()

	return tcp
}

func (t *TcpStruct) run() {
	go t.goRecv()
	go t.goSend()
}
