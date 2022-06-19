package connmgmt

import (
	"net"
	"time"

	"github.com/seanmcadam/octovpn/connection"
	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/octoconfig"
	"github.com/seanmcadam/octovpn/packet"
	"github.com/seanmcadam/octovpn/pinger"
)

type connectionInterface interface {
	//
	// Start the interface goroutines
	Start()
	//
	// Stop the server
	Stop()
	//
	// Send a predefined GOB struct vi this interface
	Write(*packet.ProtoHeader) error
	//
	//
	Protocol() octoconfig.ConnectionProtocol
	//
	// Get the local address of the interface
	LocalAddr() net.Addr
	//
	// Get the remote address of the interface
	RemoteAddr() net.Addr
	//
	// Get the channel that delivers the predefined GOB structures send over the interface connection
	ReadChan() chan *packet.ProtoHeader
}

type ConnMgmtStruct struct {
	ctx        *ctx.Ctx
	connection connectionInterface
	startTime  time.Time
	pinger     *pinger.PingerStruct
	readChan   chan *packet.EthFrame
}

//
// New() connection.
//
func New(cx *ctx.Ctx, t *octoconfig.ConfigTarget) (cm *ConnMgmtStruct, e error) {

	cx = cx.NewWithCancel()

	var c connectionInterface

	c, e = connection.NewConn(cx, t)

	cm = &ConnMgmtStruct{
		ctx:        cx,
		connection: c,
		startTime:  time.Now(),
		readChan:   make(chan *packet.EthFrame),
	}
	cm.pinger = pinger.NewPinger(cx, 5)

	return cm, e
}

//
//
//
func (c *ConnMgmtStruct) Start() {
	c.connection.Start()
	c.pinger.Start()
}

//
//
//
func (c *ConnMgmtStruct) Stop() {
	c.pinger.Stop()
	c.connection.Stop()
}

//
//
//

func (c *ConnMgmtStruct) ReadChan() chan *packet.EthFrame {
	return c.readChan
}

//
//
//
func (c *ConnMgmtStruct) Online() bool {
	return false
}

//
//
//
func (s *ConnMgmtStruct) Latency() pinger.Latency {
	return 0
}

//
//
//
func (s *ConnMgmtStruct) Loss() pinger.Loss {
	return 0
}

//
//
//
func (s *ConnMgmtStruct) Deviation() pinger.Deviation {
	return 0
}

//
//
//
func (c *ConnMgmtStruct) SendProto(data *packet.ProtoHeader) error {
	return c.connection.Write(data)
}

//
//
//
func (c *ConnMgmtStruct) LocalAddr() net.Addr {
	return c.connection.LocalAddr()
}

//
//
//
func (c *ConnMgmtStruct) RemoteAddr() net.Addr {
	return c.connection.RemoteAddr()
}
