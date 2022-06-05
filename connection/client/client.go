package client

import (
	"errors"
	"net"
	"time"

	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/octoconfig"
	"github.com/seanmcadam/octovpn/tcpclient"
	"github.com/seanmcadam/octovpn/udpclient"
)

var ErrBadProtocol = errors.New("bad Client Protocol")

type clientInterface interface {
	Start()
	Send(interface{}) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
}

type clientStruct struct {
	ctx       *ctx.Ctx
	ID        uint64
	protocol  octoconfig.ConnectionProtocol
	hostname  string
	port      uint16
	mtu       uint16
	client    clientInterface
	startTime time.Time
	//	state             connection.ConnState
	//	status            connection.ConnStatus
	//	writeBuf          *bufio.Writer
	//	decoder           *gob.Decoder
	//	encoder           *gob.Encoder
}

//
// New() client connection.
//
func New(cx *ctx.Ctx, t *octoconfig.ConfigTarget, readChan <-chan interface{}) (client *clientStruct, e error) {

	cx = cx.NewWithCancel()

	client = &clientStruct{
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
		client.client = tcpclient.New(client.protocol, client.hostname, client.port, readChan)
	case octoconfig.UDP:
		fallthrough
	case octoconfig.UDP4:
		fallthrough
	case octoconfig.UDP6:
		client.client = udpclient.New(client.protocol, client.hostname, client.port, readChan)
	default:
		e = ErrBadProtocol
		return nil, e
	}

	//
	// other initilization stuff
	//
	return client, e
}

//
//
//
func (c *clientStruct) Start() {

}

//
//
//
func (c *clientStruct) Send(data interface{}) error {
	return c.client.Send(data)
}
