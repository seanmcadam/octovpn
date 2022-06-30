package connmgmt

import (
	"net"
	"time"

	"github.com/seanmcadam/octovpn/connection"
	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/octoconfig"
	"github.com/seanmcadam/octovpn/packet"
)

type connectionInterface interface {
	//
	// Start the interface goroutines
	Start()
	//
	// Stop the server
	Stop()
	//
	//
	Write(*packet.CommonHeader) (int, error)
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
	//
	ReadChan() chan *connection.ConnReadStruct
}

type ConnMgmtStruct struct {
	ctx        *ctx.Ctx
	connection connectionInterface
	startTime  time.Time
	pinger     *packet.PingerStruct
	readChan   chan *packet.ConnFrame
	writeChan  chan *packet.CommonHeader
}

//	Start()
//	Stop()
//	LocalAddr() net.Addr
//	RemoteAddr() net.Addr
//	ReadChan() chan *packet.ConnFrame
//	Write(*packet.ConnFrame) error
//	Online() bool
//	Loss() pinger.Loss           // Loss calculation 0-1000 - 0 = best
//	Latency() pinger.Latency     // Latency calulation 0-1000 - 0 = best
//	Deviation() pinger.Deviation // Deviation calulation 0-1000 - 0 = best

//
// NewStruct() connection.
// MgmtStruct generated from Accepted connections (server)
//
func NewMgmtStruct(cx *ctx.Ctx, conn *connection.ConnectionStruct) (cm *ConnMgmtStruct, e error) {

	cx = cx.NewWithCancel()
	var c connectionInterface
	c = conn

	cm = &ConnMgmtStruct{
		ctx:        cx,
		connection: c,
		startTime:  time.Now(),
		readChan:   make(chan *packet.ConnFrame),
		writeChan:  make(chan *packet.CommonHeader),
	}
	cm.pinger = packet.NewPinger(cx, 5)

	return cm, e
}

//
// NewConfigStruct() connection
// MgmtStruct generated from config files (client)
//
func NewClient(cx *ctx.Ctx, t *octoconfig.ConfigTarget) (cm *ConnMgmtStruct, e error) {

	cx.LogLocation()

	cx = cx.NewWithCancel()
	var c connectionInterface
	c, e = connection.NewConn(cx, t)

	if e == nil {
		cm = &ConnMgmtStruct{
			ctx:        cx,
			connection: c,
			startTime:  time.Now(),
			readChan:   make(chan *packet.ConnFrame),
			writeChan:  make(chan *packet.CommonHeader),
		}
		cm.pinger = packet.NewPinger(cx, 5)
	}

	return cm, e
}

//
//
//
func (c *ConnMgmtStruct) Start() {
	c.connection.Start()
	c.pinger.Start()
	go c.goRunReader()
	go c.goRunWriter()
	go c.goSendPingPong()
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
func (c *ConnMgmtStruct) ReadChan() chan *packet.ConnFrame {
	return c.readChan
}

//
//
//
func (c *ConnMgmtStruct) Write(data *packet.ConnFrame) error {
	proto := packet.NewHeaderV1Payload(data)
	c.writeChan <- proto

	return nil
}

//
//
//
func (c *ConnMgmtStruct) goSendPingPong() {
	pingerchan := c.pinger.SendChan()
	for {
		select {
		case <-c.ctx.DoneChan():
			return
		case pingpong := <-pingerchan:
			proto := packet.NewHeaderV1Payload(pingpong)
			c.writeChan <- proto

			select {
			case c.writeChan <- proto:
			default:
				c.ctx.Logf(ctx.LogLevelError, "pingpong write failed")
			}
		}
	}
}

//
// goRunWriter()
// Serializes writes to the connection
//
func (c *ConnMgmtStruct) goRunWriter() {
	c.ctx.LogLocation()
	for {
		select {
		case <-c.ctx.DoneChan():
			return
		case proto := <-c.writeChan:
			count, e := c.connection.Write(proto)
			_ = count
			if e != nil {
				c.ctx.Logf(ctx.LogLevelPanic, " connection.Write() error:%s", e)
			}
		}
	}
}

//
//
//
func (c *ConnMgmtStruct) goRunReader() {
	c.ctx.LogLocation()
	readConnChan := c.connection.ReadChan()
	for {
		select {
		case <-c.ctx.DoneChan():
			return
		case readStruct := <-readConnChan:
			err := readStruct.Err()
			if err != nil {
				panic("")
			}

			packettype := readStruct.Header().GetType()
			payload := readStruct.Header().GetPayload()

			switch packettype {
			case packet.EthPacket:
				panic("")
			case packet.CommPacket:
				switch payload := payload.(type) {
				case *packet.ConnFrame:
					connFrame := payload
					c.readChan <- connFrame
				default:
					panic("")
				}
			case packet.PingPacket:
				c.pinger.GotPing(payload.(*packet.Ping))
			case packet.PongPacket:
				c.pinger.GotPong(payload.(*packet.Pong))

			default:
				c.ctx.Logf(ctx.LogLevelPanic, "Default reached for type:%t", payload)
			}
		}
	}
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
func (s *ConnMgmtStruct) Latency() packet.Latency {
	return 0
}

//
//
//
func (s *ConnMgmtStruct) Loss() packet.Loss {
	return 0
}

//
//
//
func (s *ConnMgmtStruct) Deviation() packet.Deviation {
	return 0
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
