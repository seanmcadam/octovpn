package connmgmt

import (
	"net"
	"time"

	"github.com/seanmcadam/octovpn/connection"
	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/octoconfig"
)

// func NewListen(cx *ctx.Ctx, t *octoconfig.ConfigListen, addConnChan chan *ConnMgmtStruct) (cm *ConnMgmtListen, e error)
// func (c *ConnMgmtListen) Start()
// func (c *ConnMgmtListen) Stop()
// func (c *ConnMgmtListen) goAccept()

type listenerInterface interface {
	//
	// Start the listener gooutines
	Start()
	//
	// Stop the listener
	Stop()
	//
	// Get Local Listen Address
	LocalAddr() net.Addr
	//
	// Get Local Listen Address
	// AcceptChan() chan *ConnMgmtStruct
	AcceptChan() chan interface{}
}

type ConnMgmtListen struct {
	ctx         *ctx.Ctx
	startTime   time.Time
	listen      listenerInterface
	addConnChan chan *ConnMgmtStruct
	acceptChan  chan interface{}
}

//
// NewListen() connection.
// Create Listener Struct
// Track addConnChan
//
func NewListen(cx *ctx.Ctx, t *octoconfig.ConfigListen, addConnChan chan *ConnMgmtStruct) (cm *ConnMgmtListen, e error) {
	cx = cx.NewWithCancel()

	cm = &ConnMgmtListen{
		ctx:         cx,
		startTime:   time.Now(),
		addConnChan: addConnChan,
		acceptChan:  make(chan interface{}),
	}

	switch t.Protocol {
	case octoconfig.TCP:
		fallthrough
	case octoconfig.TCP4:
		fallthrough
	case octoconfig.TCP6:
		cm.listen, e = connection.NewTCPListen(cx, t)
	case octoconfig.UDP:
		fallthrough
	case octoconfig.UDP4:
		fallthrough
	case octoconfig.UDP6:
		cm.listen, e = connection.NewUDPListen(cx, t)
		//fallthrough
	default:
		cx.Logf(ctx.LogLevelPanic, "default reached %s", t.Protocol)

	}

	return cm, e
}

func (c *ConnMgmtListen) Start() {
	go c.goAccept()
}

func (c *ConnMgmtListen) Stop() {
	c.ctx.Cancel()
}

//
//
//
func (c *ConnMgmtListen) goAccept() {
	for {
		select {
		case <-c.ctx.DoneChan():
			return

		case accept := <-c.acceptChan:
			switch accept := accept.(type) {
			case *connection.ConnectionStruct:
				ctx := c.ctx.NewWithCancel()
				conn, e := NewMgmtStruct(ctx, accept)
				if e != nil {
					c.ctx.LogPanicf(" NewAccept() error:%s", e)
				}
				c.addConnChan <- conn
			default:
				c.ctx.Logf(ctx.LogLevelPanic, "default %t", accept)
			}

		}
	}
}

//
// NewAccept() connection.
//
//func newAccept(cx *ctx.Ctx, net net.Conn) (cm *ConnMgmtStruct, e error) {
//
//	cx = cx.NewWithCancel()
//	var c connectionInterface = connection.NewStruct(cx, net, nil)
//
//	cm = &ConnMgmtStruct{
//		ctx:        cx,
//		connection: c,
//		startTime:  time.Now(),
//		readChan:   make(chan *packet.ConnFrame),
//	}
//	cm.pinger = pinger.NewPinger(cx, 5)
//
//	return cm, e
//}
//
