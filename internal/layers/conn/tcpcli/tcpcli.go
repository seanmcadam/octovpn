package tcpcli

import (
	"fmt"
	"net"
	"time"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/layers/network/tcp"
	"github.com/seanmcadam/octovpn/internal/link"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type TcpClientStruct struct {
	cx      *ctx.Ctx
	link    *link.LinkStateStruct
	config  *settings.ConnectionStruct
	address string
	tcpaddr *net.TCPAddr
	tcpconn *tcp.TcpStruct
	recvch  chan *packet.PacketStruct
}

func New(ctx *ctx.Ctx, config *settings.ConnectionStruct) (tcpclient interfaces.ConnInterface, err error) {
	return new(ctx, config)
}

func new(ctx *ctx.Ctx, config *settings.ConnectionStruct) (tcpclient *TcpClientStruct, err error) {

	t := &TcpClientStruct{
		cx:      ctx,
		link:    link.NewLinkState(ctx),
		config:  config,
		address: fmt.Sprintf("%s:%d", config.Host, config.Port),
		tcpaddr: nil,
		tcpconn: nil,
		recvch:  make(chan *packet.PacketStruct),
	}

	// Do an initial check and fail if it fails
	// Recheck this each time, the IP could change or rotate
	t.tcpaddr, err = net.ResolveTCPAddr(string(t.config.Proto), t.address)
	if err != nil {
		return nil, err
	}

	t.link.Down()

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
	if t == nil {
		log.ErrorStack("Nil Method Pointer")
		return
	}

	defer t.Cancel()

TCPFOR:
	for {
		var err error
		var conn *net.TCPConn

		if t.tcpconn != nil {
			log.FatalStack("should be nil")
		}

		// Dial it and keep trying forever
		conn, err = net.DialTCP(string(t.config.Proto), nil, t.tcpaddr)

		if err != nil || conn == nil{
			log.Warnf("connection failed %s: %s, wait", t.address, err)
			t.tcpconn = nil
			t.link.Down()
			time.Sleep(1 * time.Second)
			continue TCPFOR
		}

		log.Info("New TCP Connection")

		t.tcpconn = tcp.NewTCP(t.cx.NewWithCancel(), conn)
		if t.tcpconn == nil {
			log.Fatal("tcpconn == nil")
		}

		t.link.Link()
		t.link.Up()
		t.link.AddLink(t.tcpconn.Link().LinkStateCh)

	TCPCLOSE:
		for {
			select {
			case <-t.tcpconn.Link().LinkCloseCh():
				log.Debug("TCPCli TCP Closed, restart")
				t.tcpconn = nil
				break TCPCLOSE

			case <-t.doneChan():
				log.Debug("TCPCli Closing Down")
				return
			}
		}
	}
}

func (t *TcpClientStruct) Link() *link.LinkStateStruct {
	return t.link
}

