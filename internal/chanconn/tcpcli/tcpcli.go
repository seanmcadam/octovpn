package tcpcli

import (
	"fmt"
	"net"
	"time"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/chanconn/tcp"
	"github.com/seanmcadam/octovpn/internal/link"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type TcpClientStruct struct {
	cx      *ctx.Ctx
	link    link.LinkStateStruct
	config  *settings.NetworkStruct
	address string
	tcpaddr *net.TCPAddr
	tcpconn *tcp.TcpStruct
	recvch  chan *packet.PacketStruct
}

func New(ctx *ctx.Ctx, config *settings.NetworkStruct) (tcpclient interfaces.ConnInterface, err error) {

	t := &TcpClientStruct{
		cx:      ctx,
		link:    *link.NewLinkState(ctx),
		config:  config,
		address: fmt.Sprintf("%s:%d", config.GetHost(), config.GetPort()),
		tcpaddr: nil,
		tcpconn: nil,
		recvch:  make(chan *packet.PacketStruct),
	}

	// Do an initial check and fail if it fails
	// Recheck this each time, the IP could change or rotate
	t.tcpaddr, err = net.ResolveTCPAddr(t.config.Proto, t.address)
	if err != nil {
		return nil, err
	}

	go t.goRun()
	return t, err
}

// goRun()
// Loop on
//
//	Establish Connection
//	Start Send and Recv Goroutines
//	Monitor reset request
func (t *TcpClientStruct) goRun() {

	defer func(t *TcpClientStruct) {
		t.cx.Cancel()
		if t.tcpconn != nil {
			t.tcpconn.Cancel()
			t.tcpconn = nil
		}

	}(t)

TCPFOR:
	for {
		var err error
		var conn *net.TCPConn

		// Dial it and keep trying forever
		conn, err = net.DialTCP(t.config.Proto, nil, t.tcpaddr)

		if err != nil {
			log.Warnf("connection failed %s: %s, wait", t.address, err)
			t.tcpconn = nil
			time.Sleep(1 * time.Second)

			if conn != nil {
				select {
				case <-t.cx.DoneChan():
					log.Debug("tcpcli goRun() closed")
					return
				default:
				}
			}
			continue TCPFOR
		}

		if conn == nil {
			log.Fatal("conn == nil")
		}

		log.Info("New TCP Connection")

		t.tcpconn = tcp.NewTCP(t.cx.NewWithCancel(), conn)
		if t.tcpconn == nil {
			log.Fatal("tcpconn == nil")
		}
		t.link.ToggleState(link.LinkStateUp)

		//closech := t.tcpconn.Closech

	TCPCLOSE:
		for {
			tcplink := t.tcpconn.LinkToggleCh()
			select {
			case state := <-tcplink:
				log.Debug("TCPCli Link Toggled Down")
				if state == link.LinkStateDown {
					t.link.ToggleState(link.LinkStateDown)
					break TCPCLOSE
				}
			case <-t.cx.DoneChan():
				log.Debug("TCPCli Closing Down")
				return

			case <-t.tcpconn.DoneChan():
				log.Debug("TCPCli Channel Closed")
				t.tcpconn = nil
				break TCPCLOSE
			}
		}
	}
}



func (t *TcpClientStruct) StateToggleCh() <- chan link.LinkStateType {
	return t.link.StateToggleCh()
}
func (t *TcpClientStruct) GetState() link.LinkStateType {
	return t.link.GetState()
}

