package connection

import (
	"errors"
	"net"
	"time"

	"github.com/seanmcadam/octovpn/connection/client"
	"github.com/seanmcadam/octovpn/connection/server"
	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/iface"
	"github.com/seanmcadam/octovpn/octoconfig"
	"github.com/seanmcadam/octovpn/octolib"
	"github.com/seanmcadam/octovpn/packet"
)

//  Package to manage the conntions to the target vpn system
//
//

type ConnState string

const ConnStateNew = "new"
const ConnStateRunning = "running"
const ConnStateError = "error"
const ConnStateClosed = "closed"

var ErrConnectionClosed = errors.New("connection closed")

type ConnectionInterface interface {
}

type Connection struct {
	ctx          *ctx.Ctx
	iface        *iface.IFace
	inFrame      chan *ConnFrame
	frameTracker *ConnFrameTrackerStruct
	connections  map[uint64]interface{}
	addConnChan  chan interface{}
}

type ConnStatus struct {
	FramesIn      uint64
	FramesOut     uint64
	packetLossOut uint64
	packetLossIn  uint64
	// RT Latency
	// Bandwidth in
	// Bandwidth out
}

type ConnStruct struct {
	ctx       *ctx.Ctx
	conn      net.Conn
	startTime time.Time
	state     ConnState
	status    ConnStatus
}

//type connInterface interface {
//	Recv() (*ConnFrame, error)
//	Send(*ConnFrame) error
//	Protocol() octoconfig.ConnectionProtocol
//	HostIP() string
//	Port() uint16
//	MTU() uint16
//	State() ConnState
//	Status() ConnStatus
//}

var counterConnectionChan chan uint64
var counterConnFrameIDChan chan uint64
var counterConnChanIDChan chan uint64

//
//
//
func init() {
	counterConnectionChan = octolib.RunGoCounter64()
	counterConnFrameIDChan = octolib.RunGoCounter64()
	counterConnChanIDChan = octolib.RunGoCounter64()
}

//
//
//
func New(cx *ctx.Ctx, conf octoconfig.ConfigV1, iface *iface.IFace) (c *Connection, e error) {

	cx = cx.NewWithCancel()

	if len(conf.Targ) == 0 && len(conf.List) == 0 {
		cx.Logf(ctx.LogLevelPanic, "")
	}

	c = &Connection{
		ctx:          cx,
		iface:        iface,
		inFrame:      make(chan *ConnFrame),
		frameTracker: newFrameTracker(cx),
		connections:  make(map[uint64]interface{}),
		addConnChan:  make(chan interface{}),
	}
	go c.goAddConnection()

	for _, j := range conf.Targ {
		if j.Active {
			client, e := client.New(cx, j, c.inFrame)
			if e != nil {
				cx.Logf(ctx.LogLevelPanic, "error:%s", e)
			}
			c.addConnChan <- client
		}
	}

	var listen interface{}
	for _, j := range conf.List {
		if j.Active {
			switch j.Protocol {
			case "udp":
				fallthrough
			case "udp4":
				fallthrough
			case "udp6":
				listen, e = server.New(cx, j, c.inFrame)

			case "tcp":
				fallthrough
			case "tcp4":
				fallthrough
			case "tcp6":
				listen, e = server.New(cx, j, c.addConnChan)

			default:
				cx.Logf(ctx.LogLevelPanic, "protocol nt supported:%s", j.Protocol)
			}

			if e != nil {
				cx.Logf(ctx.LogLevelPanic, "error:%s", e)
			}
			c.addConnChan <- listen
		}
	}

	return c, e
}

//
//
//
func (c *Connection) goAddConnection() {

	for {
		var conn interface{}

		connID := <-counterConnectionChan
		select {
		case conn = <-c.addConnChan:
		case <-c.ctx.DoneChan():
			return
		}

		switch conn.(type) {
		case *listenTCPStruct:
			// conn.(*listenTCPStruct).ID = connID
		case *TCPStruct:
			conn.(*TCPStruct).ID = connID
		case *listenUDPStruct:
			conn.(*listenUDPStruct).ID = connID
		case *targetStruct:
			conn.(*targetStruct).ID = connID
		default:
			c.ctx.Logf(ctx.LogLevelPanic, "Type not supported:%t", conn)
		}

		c.connections[connID] = conn
	}
}

//
//
//
func (c *Connection) newConnFrame(frame *packet.EthFrame) (cf *ConnFrame) {
	c.ctx.LogLocation()
	cf = &ConnFrame{
		ID:    connFrameID(<-counterConnChanIDChan),
		Frame: frame,
	}
	return cf
}

//
// Ctx()
//
//func (c *Connection) Ctxx() *ctx.Ctx {
//	return c.ctx
//}

//
//
//
func (c *Connection) Start() {
	go c.goRunIFRead()
	go c.goRunConnRead()
}

//
//
//
func (c *Connection) Stop() {
	c.ctx.Cancel()
}

//
// goRunIFRead()
//
func (c *Connection) goRunIFRead() {
	for {
		select {
		case i := <-c.iface.ReadChan():
			c.handleIfFrame(i)
		case <-c.ctx.DoneChan():
			return
		}
	}
}

//
//
//
func (c *Connection) goRunConnRead() {
	for {
		select {
		case frame := <-c.inFrame:
			c.handleConnFrame(frame)
		case <-c.ctx.DoneChan():
			return
		}
	}
}

//
//
//
func (c *Connection) handleConnFrame(conn *ConnFrame) {
	//
	// Need to track incoming IDs to see if they have passed through already, and drop if they have
	//

	if c.frameTracker.PassFrame(conn.ID) {
		c.iface.Write(conn.Frame)
	}
}

//
// hadleIfFrame()
// Take inbound frames from the interface an process them
//
func (c *Connection) handleIfFrame(frame *packet.EthFrame) {

	conn := c.newConnFrame(frame)
	c.ctx.Logf(ctx.LogLevelTrace, " got Frame ID:%d:%s", conn.ID, conn)

	if 0 == len(c.connections) {
		c.ctx.Logf(ctx.LogLevelError, " No Connections")
	}
	//
	// Send on every connection for now
	//
	for i, j := range c.connections {
		c.ctx.Logf(ctx.LogLevelTrace, " Send() on connection# %d", i)

		var e error
		switch j.(type) {
		case *targetStruct:
			e = j.(*targetStruct).Send(conn)
		case *listenUDPStruct:
			e = j.(*listenUDPStruct).Send(conn)
		case *TCPStruct:
			e = j.(*TCPStruct).SendConnFrame(conn)
		case nil:
			c.ctx.Logf(ctx.LogLevelError, " j == nil")
		default:
			c.ctx.Logf(ctx.LogLevelPanic, " default reached type:%t", j)
		}
		if e != nil {
			c.ctx.Logf(ctx.LogLevelPanic, " error:%s", e)
		}
	}
}
