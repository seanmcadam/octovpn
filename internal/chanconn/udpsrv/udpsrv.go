package udpsrv

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

type UdpServerStruct struct {
	cx      *ctx.Ctx
	config  *settings.NetworkStruct
	address string
	udpaddr *net.UDPAddr
	udpconn *udp.UdpStruct
	auth    bool
	recvch  chan *packetconn.ConnPacket
}

func New(ctx *ctx.Ctx, config *settings.NetworkStruct) (udpserver interfaces.ConnInterface, err error) {

	u := &UdpServerStruct{
		cx:      ctx,
		config:  config,
		address: fmt.Sprintf("%s:%d", config.GetHost(), config.GetPort()),
		udpaddr: nil,
		udpconn: nil,
		auth:    false,
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
func (u *UdpServerStruct) goRun() {

	defer func(u *UdpServerStruct) {
		if u.udpconn != nil {
			u.udpconn.Cancel()
			u.udpconn = nil
		}
	}(u)

	for {
		var err error
		var conn *net.UDPConn

		//
		// conn, err = net.DialUDP(u.config.Proto, nil, u.udpaddr)
		conn, err = net.ListenUDP(u.config.Proto, u.udpaddr)

		if err != nil {
			log.Warnf("UDP Listener failed %s: %s, wait", u.address, err)
			u.udpconn = nil
			time.Sleep(1 * time.Second)

			select {
			case <-u.cx.DoneChan():
				log.Debug("udpsrv goRun() closed")
				return
			default:
			}

			continue
		}

		log.Info("New UDP Connection")

		u.udpconn = udp.NewUDPSrv(u.cx.NewWithCancel(), conn)
		if u.udpconn == nil {
			log.Error("udpconn == nil")
			continue
		}

		for {
			select {
			case <-u.cx.DoneChan():
				log.Debug("UDPSrv Closing Down")
				return
			case <-u.udpconn.DoneChan():
				log.Debug("UDPSrv Channel Closed")
				u.udpconn = nil
				break
			}
		}
	}
}
