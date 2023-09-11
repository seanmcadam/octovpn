package udpsrv

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

type UdpServerStruct struct {
	cx      *ctx.Ctx
	link    *link.LinkStateStruct
	config  *settings.ConnectionStruct
	address string
	udpaddr *net.UDPAddr
	udpconn *udp.UdpStruct
	auth    bool
	recvch  chan *packet.PacketStruct
}

func New(ctx *ctx.Ctx, config *settings.ConnectionStruct) (udpserver interfaces.ConnInterface, err error) {
	return new(ctx, config)
}

func new(ctx *ctx.Ctx, config *settings.ConnectionStruct) (udpserver *UdpServerStruct, err error) {

	u := &UdpServerStruct{
		cx:      ctx,
		link:    link.NewLinkState(ctx),
		config:  config,
		address: fmt.Sprintf("%s:%d", config.Host, config.Port),
		udpaddr: nil,
		udpconn: nil,
		auth:    false,
		recvch:  make(chan *packet.PacketStruct, 16),
	}

	// Do an initial check and fail if it fails
	// Recheck this each time, the IP could change or rotate
	u.udpaddr, err = net.ResolveUDPAddr(string(u.config.Proto), u.address)
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
func (u *UdpServerStruct) goRun() {
	if u == nil {
		return
	}

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
		log.Debug("ListenUDP()")
		conn, err = net.ListenUDP(string(u.config.Proto), u.udpaddr)

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

		log.Debug("UDP Srv Conn UP")
		u.link.AddLink(u.udpconn.Link().LinkStateCh)
		u.link.Up()

		for {
			select {
			case <-u.link.LinkUpCh():
			case <-u.link.LinkCloseCh():
				log.Debug("UDPSrv Link Down")
				u.udpconn = nil
				break
			case <-u.link.LinkDownCh():
				log.Debug("UDPSrv Link Down")
				u.udpconn = nil
				break
			case <-u.udpconn.DoneChan():
				log.Debug("UDPSrv Channel Closed")
				u.udpconn = nil
				break
			case <-u.cx.DoneChan():
				log.Debug("UDPSrv Closing Down")
				return
			}
		}
	}
}

func (t *UdpServerStruct) Link() *link.LinkStateStruct {
	return t.link
}

func (u *UdpServerStruct) GetLinkNoticeStateCh() link.LinkNoticeStateCh {
	if u == nil {
		return nil
	}

	return u.link.LinkNoticeStateCh()
}

func (u *UdpServerStruct) GetLinkStateCh() link.LinkNoticeStateCh {
	if u == nil {
		return nil
	}

	return u.link.LinkStateCh()
}
func (u *UdpServerStruct) GetUpCh() link.LinkNoticeStateCh {
	if u == nil {
		return nil
	}

	return u.link.LinkUpCh()
}

func (u *UdpServerStruct) GetLinkCh() link.LinkNoticeStateCh {
	if u == nil {
		return nil
	}

	return u.link.LinkLinkCh()
}

func (u *UdpServerStruct) GetDownCh() link.LinkNoticeStateCh {
	if u == nil {
		return nil
	}

	return u.link.LinkDownCh()
}

func (u *UdpServerStruct) GetState() link.LinkStateType {
	if u == nil {
		return link.LinkStateERROR
	}

	return u.link.GetState()
}
