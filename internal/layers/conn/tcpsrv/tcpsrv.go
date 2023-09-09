package tcpsrv

import (
	"fmt"
	"net"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/layers/network/tcp"
	"github.com/seanmcadam/octovpn/internal/link"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type TcpServerStruct struct {
	cx          *ctx.Ctx
	link        *link.LinkStateStruct
	config      *settings.NetworkStruct
	address     string
	tcplistener *net.TCPListener
	tcpaddr     *net.TCPAddr
	tcpconn     *tcp.TcpStruct
	tcpconnch   chan *tcp.TcpStruct
	recvch      chan *packet.PacketStruct
}

// NewTcpServer()
// Returns a TcpServerStruct and error value
func New(ctx *ctx.Ctx, config *settings.NetworkStruct) (tcpserver interfaces.ConnInterface, err error) {

	t := &TcpServerStruct{
		cx:          ctx,
		link:        link.NewLinkState(ctx, link.LinkModeUpAND), // If more then 1 is connected, they all have to be up
		config:      config,
		address:     fmt.Sprintf("%s:%d", config.GetHost(), config.GetPort()),
		tcplistener: nil,
		tcpaddr:     nil,
		tcpconn:     nil,
		tcpconnch:   make(chan *tcp.TcpStruct),
		recvch:      make(chan *packet.PacketStruct, 16),
	}

	// Recheck this each time, the IP could change or rotate
	t.tcpaddr, err = net.ResolveTCPAddr(t.config.Proto, t.address)
	if err != nil {
		return nil, fmt.Errorf("ResolveTCPAddr Failed:%s", err)
	}

	t.tcplistener, err = net.ListenTCP(t.config.Proto, t.tcpaddr)
	if err != nil {
		return nil, fmt.Errorf("ListenTCP Failed:%s", err)
	}

	// This is the server, so the connection is down to start with.
	t.link.Down()

	go t.goListen()
	go t.goRun()

	return t, err

}

// goRun()
// Loop on
// 	Establish Connection
// 	Start Send and Recv Goroutines
// 	Monitor reset request
//

func (t *TcpServerStruct) goRun() {

	defer func(t *TcpServerStruct) {
		if t.tcplistener != nil {
			t.tcplistener.Close()
			t.tcplistener = nil
		}
		if t.tcpconn != nil {
			t.tcpconn.Cancel()
			t.emptyconn()
			t.tcpconn = nil
		}
		close(t.tcpconnch)
	}(t)

	for {
		var tcpconnclosech chan interface{}

		select {
		case conn := <-t.tcpconnch:
			log.Debugf("New incoming TCP Connection")

			// Terminate last connection

			if t.tcpconn != nil {
				log.Debug("Shutdown Previous connection")
				t.tcpconn.Cancel()
				t.tcpconn = nil
			}

			t.tcpconn = conn
			t.link.Link()
			t.link.Up()
			t.link.AddLink(t.tcpconn.Link().LinkNoticeStateCh)

			log.Debugf("TCP Srv state:%s", conn.Link().GetState())

		case <-tcpconnclosech:
			continue

		case <-t.cx.DoneChan():
			return

		}
	}
}

func (t *TcpServerStruct) emptyconn() {
	for {
		select {
		case <-t.tcpconnch:
		default:
			return
		}
	}
}

func (t *TcpServerStruct) GetLinkNoticeStateCh() link.LinkNoticeStateCh {
	return t.link.LinkNoticeStateCh()
}

func (t *TcpServerStruct) GetLinkStateCh() link.LinkNoticeStateCh {
	return t.link.LinkStateCh()
}

func (t *TcpServerStruct) GetUpCh() link.LinkNoticeStateCh {
	return t.link.LinkUpCh()
}

func (t *TcpServerStruct) GetDownCh() link.LinkNoticeStateCh {
	return t.link.LinkDownCh()
}

func (t *TcpServerStruct) GetLinkCh() link.LinkNoticeStateCh {
	return t.link.LinkLinkCh()
}

func (t *TcpServerStruct) GetState() link.LinkStateType {
	return t.link.GetState()
}
