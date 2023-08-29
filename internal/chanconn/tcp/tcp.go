package tcp

import (
	"net"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/octolib/ctx"
)

type TcpStruct struct {
	cx     *ctx.Ctx
	conn   *net.TCPConn
	sendch chan interfaces.PacketInterface
	recvch chan interfaces.PacketInterface
}

func NewTCP(ctx *ctx.Ctx, conn *net.TCPConn) (tcp *TcpStruct) {

	tcp = &TcpStruct{
		cx:     ctx,
		conn:   conn,
		sendch: make(chan interfaces.PacketInterface),
		recvch: make(chan interfaces.PacketInterface),
	}

	tcp.run()

	return tcp
}

func (t *TcpStruct) run() {
	go t.goRecv()
	go t.goSend()
}
