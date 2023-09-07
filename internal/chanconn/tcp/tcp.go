package tcp

import (
	"net"

	"github.com/seanmcadam/octovpn/internal/link"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/ctx"
)

type TcpStruct struct {
	cx     *ctx.Ctx
	link   link.LinkStateStruct
	conn   *net.TCPConn
	sendch chan *packet.PacketStruct
	recvch chan *packet.PacketStruct
}

func NewTCP(ctx *ctx.Ctx, conn *net.TCPConn) (tcp *TcpStruct) {

	tcp = &TcpStruct{
		cx:     ctx,
		link:   *link.NewLinkState(ctx),
		conn:   conn,
		sendch: make(chan *packet.PacketStruct),
		recvch: make(chan *packet.PacketStruct),
	}

	tcp.run()

	return tcp
}

func (t *TcpStruct)LinkToggleCh() <-chan link.LinkStateType {
	return t.link.StateToggleCh()
}

func (t *TcpStruct) run() {
	t.link.ToggleState(link.LinkStateUp)
	go t.goRecv()
	go t.goSend()
}
