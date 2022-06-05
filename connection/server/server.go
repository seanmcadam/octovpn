package server

import (
	"errors"
	"net"
	"time"

	"github.com/seanmcadam/octovpn/connection/client/tcp"
	"github.com/seanmcadam/octovpn/connection/client/udp"
	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/octoconfig"
)

var ErrBadProtocol = errors.New("bad Server Protocol")

type serverInterface interface {
	Start()
	Send(interface{}) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
}

type serverStruct struct {
	ctx       *ctx.Ctx
	ID        uint64
	protocol  octoconfig.ConnectionProtocol
	hostname  string
	port      uint16
	mtu       uint16
	server    serverInterface
	startTime time.Time
	//	state             connection.ConnState
	//	status            connection.ConnStatus
	//	readConnFrameChan chan *connection.ConnFrame
	//	writeBuf          *bufio.Writer
	//	decoder           *gob.Decoder
	//	encoder           *gob.Encoder
}

//
// New() server connection.
//
func New(cx *ctx.Ctx, t *octoconfig.ConfigTarget, readChan <-chan interface{}) (server *serverStruct, e error) {

	cx = cx.NewWithCancel()

	server = &serverStruct{
		ctx:       cx,
		protocol:  t.Protocol,
		hostname:  t.Hostname,
		port:      t.Port,
		mtu:       t.MTU,
		startTime: time.Now(),
	}

	switch t.Protocol {
	case octoconfig.TCP:
		fallthrough
	case octoconfig.TCP4:
		fallthrough
	case octoconfig.TCP6:
		server.server = tcp.New(server.protocol, server.hostname, server.port, readChan)
	case octoconfig.UDP:
		fallthrough
	case octoconfig.UDP4:
		fallthrough
	case octoconfig.UDP6:
		server.server = udp.New(server.protocol, server.hostname, server.port, readChan)
	default:
		e = ErrBadProtocol
		return nil, e
	}

	//
	// other initilization stuff
	//
	return server, e
}

//
//
//
func (c *serverStruct) Start() {

}

//
//
//
func (c *serverStruct) Send(data interface{}) error {
	return c.server.Send(data)
}
