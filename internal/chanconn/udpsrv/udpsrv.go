package udpsrv

import (
	"fmt"
	"net"
	"time"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/chanconn/udp"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type UdpServerStruct struct {
	config  *settings.NetworkStruct
	address string
	udpaddr *net.UDPAddr
	udpconn *udp.UdpStruct
	auth    bool
	closech chan interface{}
	resetch chan interface{}
}

func New(config *settings.NetworkStruct) (udpserver interfaces.ChannelInterface, err error) {

	u := &UdpServerStruct{
		config:  config,
		address: fmt.Sprintf("%s:%d", config.GetHost(), config.GetPort()),
		udpaddr: nil,
		udpconn: nil,
		auth:    false,
		closech: make(chan interface{}),
		resetch: make(chan interface{}),
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
			u.udpconn.Close()
			u.udpconn = nil
		}
		close(u.closech)
		close(u.resetch)

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
			case <-u.closech:
				log.Debug("udpcli goRun() closed")
				return
			default:
			}

			continue
		}

		log.Info("New UDP Connection")

		u.udpconn = udp.NewUDPSrv(conn)
		if u.udpconn == nil {
			log.Error("udpconn == nil")
			continue
		}

		closech := u.udpconn.Closech

		for {
			select {
			case <-u.closech:
				log.Debug("UDPSrv Closing Down")
				return
			case <-closech:
				log.Debug("UDPSrv Channel Closed")
				u.udpconn = nil
				break
			}
		}
	}
}
