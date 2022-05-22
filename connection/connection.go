package connection

import (
	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/iface"
	"github.com/seanmcadam/octovpn/octoconfig"
	"github.com/seanmcadam/octovpn/octolib"
)

//  Package to manage the conntions to the target vpn system
//
//

type listenStruct struct {
	protocol octoconfig.ConnectionProtocol
	ip       string
	port     uint16
	mtu      uint16
}

type targetStruct struct {
	protocol octoconfig.ConnectionProtocol
	hostname string
	port     uint16
	mtu      uint16
}

type Connection struct {
	ctx          *ctx.Ctx
	iface        *iface.IFace
	inFrame      chan *ConnFrame
	frameTracker *ConnFrameTrackerStruct
	target       map[string]*targetStruct
	listen       map[string]*listenStruct
}

var counterConnIDChan chan uint64

//
//
//
func init() {
	counterConnIDChan = octolib.RunGoCounter64()
}

//
//
//
func New(cx *ctx.Ctx, conf octoconfig.ConfigV1, iface *iface.IFace) (c *Connection) {

	cx = cx.NewWithCancel()

	if len(conf.Targ) == 0 && len(conf.List) == 0 {
		cx.Logf(ctx.LogLevelPanic, "")
	}

	c = &Connection{
		ctx:          cx,
		iface:        iface,
		inFrame:      make(chan *ConnFrame),
		frameTracker: newFrameTracker(),
		target:       make(map[string]*targetStruct),
		listen:       make(map[string]*listenStruct),
	}

	for i, j := range conf.Targ {
		c.target[i] = newTarget(j)
	}

	for i, j := range conf.List {
		c.listen[i] = newListen(j)
	}

	return c
}

//
//
//
func newConnFrame(frame *iface.Frame) (cf *ConnFrame) {
	cf = &ConnFrame{
		id:    connFrameID(<-counterConnIDChan),
		frame: frame,
	}
	return cf
}

//
// Ctx()
//
func (c *Connection) Ctx() *ctx.Ctx {
	return c.ctx
}

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
		case <-c.ctx.Done():
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
		case <-c.ctx.Done():
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
	c.iface.Write(conn.frame)
}

//
//
//
func (c *Connection) handleIfFrame(frame *iface.Frame) {
	conn := newConnFrame(frame)
	c.ctx.Logf(ctx.LogLevelTrace, " got Frame ID:%d:%s", conn.id, conn)

	// Need to handle multiple paths here...
	// c..Write(conn.frame)
}

func newTarget(t *octoconfig.ConfigTarget) (target *targetStruct) {
	target = &targetStruct{
		protocol: t.Protocol,
		hostname: t.Hostname,
		port:     t.Port,
		mtu:      t.MTU,
	}

	return target
}
func newListen(l *octoconfig.ConfigListen) (listen *listenStruct) {
	listen = &listenStruct{
		protocol: l.Protocol,
		ip:       l.IP,
		port:     l.Port,
		mtu:      l.MTU,
	}
	return listen
}
