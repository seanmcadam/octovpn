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
	config  *settings.NetworkStruct
	address string
	udpaddr *net.UDPAddr
	udpconn *udp.UdpStruct
	recvch  chan *packet.PacketStruct
}

func New(ctx *ctx.Ctx, config *settings.NetworkStruct) (udpclient interfaces.ConnInterface, err error) {
	return new(ctx, config)
}

func new(ctx *ctx.Ctx, config *settings.NetworkStruct) (udpclient *UdpClientStruct, err error) {

	u := &UdpClientStruct{
		cx:      ctx,
		link:    link.NewLinkState(ctx),
		config:  config,
		address: fmt.Sprintf("%s:%d", config.GetHost(), config.GetPort()),
		udpaddr: nil,
		udpconn: nil,
		recvch:  make(chan *packet.PacketStruct, 16),
	}

	// Do an initial check and fail if it fails
	// Recheck this each time, the IP could change or rotate
	u.udpaddr, err = net.ResolveUDPAddr(u.config.Proto, u.address)
	if err != nil {
		return nil, err
	}
	u.link.Down()

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

		log.Debug("UDP Cli Conn UP")
		u.link.AddLink(u.udpconn.Link().LinkStateCh)
		u.link.Up()
		for {
			select {
			case <-u.udpconn.Link().LinkUpCh():
			case <-u.udpconn.Link().LinkDownCh():
				log.Debug("UDPLink Down")
				u.udpconn = nil
				break
			case <-u.udpconn.Link().LinkCloseCh():
				log.Debug("UDPLink CLosed")
				u.udpconn = nil
				break
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

func (u *UdpClientStruct) GetLinkNoticeStateCh() link.LinkNoticeStateCh {
	if u == nil {
		return nil
	}

	return u.link.LinkNoticeStateCh()
}

func (u *UdpClientStruct) GetLinkStateCh() link.LinkNoticeStateCh {
	if u == nil {
		return nil
	}

	return u.link.LinkStateCh()
}

func (u *UdpClientStruct) GetUpCh() link.LinkNoticeStateCh {
	if u == nil {
		return nil
	}

	return u.link.LinkUpCh()
}

func (u *UdpClientStruct) GetLinkCh() link.LinkNoticeStateCh {
	if u == nil {
		return nil
	}

	return u.link.LinkLinkCh()
}

func (u *UdpClientStruct) GetDownCh() link.LinkNoticeStateCh {
	if u == nil {
		return nil
	}

	return u.link.LinkDownCh()
}

func (u *UdpClientStruct) GetState() link.LinkStateType {
	if u == nil {
		return link.LinkStateERROR
	}

	return u.link.GetState()
}
