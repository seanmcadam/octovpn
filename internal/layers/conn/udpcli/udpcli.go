package udpcli

import (
	"fmt"
	"net"
	"time"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/layers/network/udp"
	"github.com/seanmcadam/octovpn/internal/link"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type UdpClientStruct struct {
	cx      *ctx.Ctx
	link    *link.LinkStateStruct
	config  *settings.ConnectionStruct
	address string
	udpaddr *net.UDPAddr
	udpconn *udp.UdpStruct
	recvch  chan *packet.PacketStruct
}

func New(ctx *ctx.Ctx, config *settings.ConnectionStruct) (udpclient interfaces.ConnInterface, err error) {
	return new(ctx, config)
}

func new(ctx *ctx.Ctx, config *settings.ConnectionStruct) (udpclient *UdpClientStruct, err error) {

	u := &UdpClientStruct{
		cx:      ctx,
		link:    link.NewLinkState(ctx),
		config:  config,
		address: fmt.Sprintf("%s:%d", config.Host, config.Port),
		udpaddr: nil,
		udpconn: nil,
		recvch:  make(chan *packet.PacketStruct, 16),
	}

	// Do an initial check and fail if it fails
	// Recheck this each time, the IP could change or rotate
	u.udpaddr, err = net.ResolveUDPAddr(string(u.config.Proto), u.address)
	if err != nil {
		return nil, err
	}
	u.link.NoLink()

	go u.goRun()
	return u, err
}

// goRun()
// Loop on
//
//	Establish Connection
//	Start Send and Recv Goroutines
//	Monitor reset request
func (u *UdpClientStruct) goRun() {
	if u == nil {
		return
	}

	defer u.Cancel()

	for {
		var err error
		var conn *net.UDPConn

		// Dial it and keep trying forever
		conn, err = net.DialUDP(string(u.config.Proto), nil, u.udpaddr)

		if err != nil {
			log.Warnf("connection failed %s: %s, wait", u.address, err)
			u.udpconn = nil
			time.Sleep(1 * time.Second)

			if u.closed(){
				return
			}

			continue
		}

		log.Info("New UDP Connection")

		u.udpconn = udp.NewUDPCli(u.cx.NewWithCancel(), conn)
		if u.udpconn == nil {
			log.Error("udpconn == nil")
			continue
		}

		log.Debug("UDP Cli Conn UP")
		u.link.AddLinkStateCh(u.udpconn.Link())
		u.link.Connected()
		for {
			select {
			case <-u.udpconn.Link().LinkCloseCh():
				log.Debug("UDPLink CLosed")
				u.udpconn = nil
				break
			case <-u.doneChan():
				log.Debug("UDPCli Closing Down")
				return
			}
		}
	}
}

func (t *UdpClientStruct) Link() *link.LinkStateStruct {
	return t.link
}
