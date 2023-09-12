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
	config      *settings.ConnectionStruct
	address     string
	tcplistener *net.TCPListener
	tcpaddr     *net.TCPAddr
	tcpconn     *tcp.TcpStruct
	tcpconnch   chan *tcp.TcpStruct
	recvch      chan *packet.PacketStruct
}

func New(ctx *ctx.Ctx, config *settings.ConnectionStruct) (tcpserver interfaces.ConnInterface, err error) {
	return new(ctx, config)
}

func new(ctx *ctx.Ctx, config *settings.ConnectionStruct) (tcpserver *TcpServerStruct, err error) {

	t := &TcpServerStruct{
		cx:          ctx,
		link:        link.NewLinkState(ctx, link.LinkModeUpAND), // If more then 1 is connected, they all have to be up
		config:      config,
		address:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		tcplistener: nil,
		tcpaddr:     nil,
		tcpconn:     nil,
		tcpconnch:   make(chan *tcp.TcpStruct),
		recvch:      make(chan *packet.PacketStruct, 16),
	}

	// Recheck this each time, the IP could change or rotate
	t.tcpaddr, err = net.ResolveTCPAddr(string(t.config.Proto), t.address)
	if err != nil {
		return nil, fmt.Errorf("ResolveTCPAddr Failed:%s", err)
	}

	t.tcplistener, err = net.ListenTCP(string(t.config.Proto), t.tcpaddr)
	if err != nil {
		return nil, fmt.Errorf("ListenTCP Failed:%s", err)
	}

	// This is the server, so the connection is down to start with.
	t.link.NoLink()

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

	if t == nil {
		log.ErrorStack("Nil Method Pointer")
		return
	}
	
	defer t.Cancel()

	for {
		var tcpconnclosech chan interface{}

		select {
		case conn := <-t.tcpconnch:
			log.Debugf("New incoming TCP Connection")

			// Terminate last connection
			log.Debug("Need to do more here with closing old connections")

			if t.tcpconn != nil {
				log.Debug("Shutdown Previous connection")
				t.tcpconn.Cancel()
				t.tcpconn = nil
			}

			t.tcpconn = conn
			t.link.Link()
			t.link.Connected()
			t.link.AddLinkStateCh(t.tcpconn.Link())

			log.Debugf("TCP Srv state:%s", conn.Link().GetState())

		case <-tcpconnclosech:
			continue

		case <-t.cx.DoneChan():
			return

		}
	}
}

func (t *TcpServerStruct) emptyconn() {
	if t == nil {
		return
	}

	for {
		select {
		case <-t.tcpconnch:
		default:
			return
		}
	}
}

func (t *TcpServerStruct) Link() *link.LinkStateStruct {
	return t.link
}