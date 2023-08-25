package udp

import (
	"net"
	"time"

	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
	"github.com/seanmcadam/octovpn/octolib/pinger"
)

type UdpStruct struct {
	cx     *ctx.Ctx
	srv    bool
	conn   *net.UDPConn
	addr   *net.UDPAddr
	pinger *pinger.Pinger64Struct
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
		pinger: pinger.NewPinger64(ctx, time.Second, 5*time.Second),
		sendch: make(chan *packetconn.ConnPacket),
		recvch: make(chan *packetconn.ConnPacket),
	}

	udp.run()
	return udp
}

func NewUDPCli(ctx *ctx.Ctx, conn *net.UDPConn) (udp *UdpStruct) {

	log.Debug("Local Addr %s", conn.LocalAddr())

	udp = &UdpStruct{
		cx: ctx,
		srv:     false,
		conn:    conn,
		addr:    nil,
		pinger:  pinger.NewPinger64(ctx, time.Second, 5*time.Second),
		sendch:  make(chan *packetconn.ConnPacket),
		recvch:  make(chan *packetconn.ConnPacket),
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
