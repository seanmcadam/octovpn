package udp

import (
	"net"

	"github.com/seanmcadam/octovpn/internal/link"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type UdpStruct struct {
	cx     *ctx.Ctx
	link   link.LinkStateStruct
	srv    bool
	conn   *net.UDPConn
	addr   *net.UDPAddr
	sendch chan *packet.PacketStruct
	recvch chan *packet.PacketStruct
}

func NewUDPSrv(ctx *ctx.Ctx, conn *net.UDPConn) (udp *UdpStruct) {

	log.Debug("Local Addr %s", conn.LocalAddr())

	udp = &UdpStruct{
		cx:     ctx,
		link:   *link.NewLinkState(ctx),
		srv:    true,
		conn:   conn,
		addr:   nil,
		sendch: make(chan *packet.PacketStruct),
		recvch: make(chan *packet.PacketStruct),
	}

	udp.run()
	return udp
}

func NewUDPCli(ctx *ctx.Ctx, conn *net.UDPConn) (udp *UdpStruct) {

	log.Debug("Local Addr %s", conn.LocalAddr())

	udp = &UdpStruct{
		cx:     ctx,
		srv:    false,
		conn:   conn,
		addr:   nil,
		sendch: make(chan *packet.PacketStruct),
		recvch: make(chan *packet.PacketStruct),
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

func (u *UdpStruct) LinkToggleCh() <-chan link.LinkStateType {
	if u == nil{
		return nil
	}
	return u.link.StateToggleCh()
}

func (u *UdpStruct) run() {
	u.link.ToggleState(link.LinkStateUp)
	go u.goRecv()
	go u.goSend()
}
