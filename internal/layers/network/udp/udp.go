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
	link   *link.LinkStateStruct
	srv    bool
	conn   *net.UDPConn
	addr   *net.UDPAddr
	sendch chan *packet.PacketStruct
	recvch chan *packet.PacketStruct
}

func NewUDPSrv(ctx *ctx.Ctx, conn *net.UDPConn) (udp *UdpStruct) {
	if ctx == nil || conn == nil{
		log.ErrorStack("NewUDPSrv()")
		return
	}

	log.Debug("Local Addr %s", conn.LocalAddr())

	udp = &UdpStruct{
		cx:     ctx,
		link:   link.NewLinkState(ctx),
		srv:    true,
		conn:   conn,
		addr:   nil,
		sendch: make(chan *packet.PacketStruct),
		recvch: make(chan *packet.PacketStruct),
	}

	udp.link.Down()
	udp.run()
	return udp
}

func NewUDPCli(ctx *ctx.Ctx, conn *net.UDPConn) (udp *UdpStruct) {
	if ctx == nil || conn == nil{
		log.ErrorStack("NewUDPCli()")
		return
	}

	log.Debug("Local Addr %s", conn.LocalAddr())

	udp = &UdpStruct{
		cx:     ctx,
		link:   link.NewLinkState(ctx),
		srv:    false,
		conn:   conn,
		addr:   nil,
		sendch: make(chan *packet.PacketStruct),
		recvch: make(chan *packet.PacketStruct),
	}

	udp.link.Down()
	udp.run()
	return udp
}

func (u *UdpStruct) endpoint() (v string) {
	if u == nil {
		return
	}

	if u.srv {
		v = "SRV"
	} else {
		v = "CLI"
	}
	return v
}

func (u *UdpStruct) Link() *link.LinkStateStruct {
	if u == nil {
		return nil
	}

	return u.link
}

func (u *UdpStruct) run() {
	if u == nil {
		return
	}

	go u.goRecv()
	go u.goSend()
	u.link.Link()
}
