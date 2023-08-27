package udpcli

import (
	"fmt"
	"net"
	"time"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/chanconn/udp"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

type UdpClientStruct struct {
	cx      *ctx.Ctx
	config  *settings.NetworkStruct
	address string
	udpaddr *net.UDPAddr
	udpconn *udp.UdpStruct
	recvch  chan *packetconn.ConnPacket
}

func New(ctx *ctx.Ctx, config *settings.NetworkStruct) (udpclient interfaces.ConnInterface, err error) {

	u := &UdpClientStruct{
		cx:      ctx,
		config:  config,
		address: fmt.Sprintf("%s:%d", config.GetHost(), config.GetPort()),
		udpaddr: nil,
		udpconn: nil,
		recvch:  make(chan *packetconn.ConnPacket, 16),
	}

	// Do an initial check and fail if it fails
	// Recheck this each time, the IP could change or rotate
	u.udpaddr, err = net.ResolveUDPAddr(u.config.Proto, u.address)
	if err != nil {
		return nil, err
	}

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

	defer func(u *UdpClientStruct) {
		if u.udpconn != nil {
			u.udpconn.Cancel()
			u.udpconn = nil
		}
	}(u)

	for {
		var err error
		var conn *net.UDPConn

		// Dial it and keep trying forever
		conn, err = net.DialUDP(u.config.Proto, nil, u.udpaddr)

		if err != nil {
			log.Warnf("connection failed %s: %s, wait", u.address, err)
			u.udpconn = nil
			time.Sleep(1 * time.Second)

			select {
			case <-u.cx.DoneChan():
				log.Debug("udpcli goRun() closed")
				return
			default:
			}

			continue
		}

		log.Info("New UDP Connection")

		u.udpconn = udp.NewUDPCli(u.cx.NewWithCancel(), conn)
		if u.udpconn == nil {
			log.Error("udpconn == nil")
			continue
		}

		for {
			select {
			case <-u.cx.DoneChan():
				log.Debug("UDPCli Closing Down")
				return
			case <-u.udpconn.DoneChan():
				log.Debug("UDPCli Channel Closed")
				u.udpconn = nil
				break
			}
		}
	}
}
