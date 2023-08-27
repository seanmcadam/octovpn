package udp

import (
	"net"

	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

type UdpStruct struct {
	cx     *ctx.Ctx
	srv    bool
	conn   *net.UDPConn
	addr   *net.UDPAddr
	sendch chan *packetconn.ConnPacket
	recvch chan *packetconn.ConnPacket
}

func NewUDPSrv(ctx *ctx.Ctx, conn *net.UDPConn) (udp *UdpStruct) {

	log.Debug("Local Addr %s", conn.LocalAddr())

	udp = &UdpStruct{
		cx:     ctx,
		srv:    true,
		conn:   conn,
		addr:   nil,
		sendch: make(chan *packetconn.ConnPacket),
		recvch: make(chan *packetconn.ConnPacket),
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
		sendch: make(chan *packetconn.ConnPacket),
		recvch: make(chan *packetconn.ConnPacket),
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
