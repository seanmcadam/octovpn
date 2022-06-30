package routes

import (
	"net"

	"github.com/seanmcadam/octovpn/connmgmt"
	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/octoconfig"
	"github.com/seanmcadam/octovpn/octolib"
	"github.com/seanmcadam/octovpn/packet"
)

// Package to manage all routes to the target vpn system
// New()(c *ConnectionStruct)
// (c *ConnectionStruct) Start()
// (c *ConnectionStruct) Stop()
//
// Send filtered Eth packets from Transit module
// (c *ConnectionStruct) Write (eth *packet.EthFrame)
//
// Read filtered Eth Packets in Transit module
// (c *ConnectionStruct) ReadEthChan () <-chan *packet.EthFrame
//
// Handler for new incoming connections from listeners
// (c *ConnectionStruct) goAddConnection ()
//
// Collect and handle incoming packets from the connection layer
// (c *ConnectionStruct) goReadConnChan ()
//
// Collect and handle incoming packets from the connection layer
// (c *ConnectionStruct) goWriteConn ()
//
// Create wrapper for Eth Frame
// (c *ConnectionStruct) newConnFrame (frame *packet.EthFrame) (cf *packet.ConnFrame)
//

const ethChanDepth = 2

//
//  Used to manage both Client and Server
//
type RouteInterface interface {
	Start()
	Stop()
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	ReadChan() chan *packet.ConnFrame
	Write(*packet.ConnFrame) error
	Online() bool
	Loss() packet.Loss           // Loss calculation 0-1000 - 0 = best
	Latency() packet.Latency     // Latency calulation 0-1000 - 0 = best
	Deviation() packet.Deviation // Deviation calulation 0-1000 - 0 = best
}

//
// Used for both originating and recieved sockets
//
type RouteStruct struct {
	//
	// Context structure
	ctx *ctx.Ctx
	//
	// Configuration
	conf octoconfig.ConfigV1
	//
	// All connections from both client and server
	connections map[uint64]RouteInterface
	//
	// Listener Sockets on this device
	listeners []*connmgmt.ConnMgmtListen
	//
	// Tracks ConnFrame packets to make sure that only the first frame recieved passes
	frameTracker packet.ConnFrameTrackerStruct
	//
	// The channel passes new accepts from listener connections
	addConnChan chan *connmgmt.ConnMgmtStruct
	//
	// Read channel for Transit layer, only one unique EthFrame ID is passed up
	readEthChan chan *packet.EthFrame
	//
	// Read channel for Connection layer, all ConnFrames pass back this way to be filtered for uniqueness
	readConnChan chan *packet.ConnFrame
	//
	// Write Channel to connections, pushed into the processing section of sending over which Conn
	// Feed routng go routine
	writeConnChan chan *packet.ConnFrame
}

var counterConnectionChan chan uint64
var counterConnChanIDChan chan uint64

//
//
//
func init() {
	counterConnectionChan = octolib.RunGoCounter64()
	counterConnChanIDChan = octolib.RunGoCounter64()
}

//
// New() Connecton
// Takes a config struct, and returns a ConnectionStruct
// This will create listsning sockets, and start trying to connect to the target server(s)
//
func New(cx *ctx.Ctx, conf octoconfig.ConfigV1) (r *RouteStruct, e error) {

	cx = cx.NewWithCancel()

	if len(conf.Targ) == 0 && len(conf.List) == 0 {
		cx.Logf(ctx.LogLevelPanic, "")
	}

	r = &RouteStruct{
		ctx:           cx,
		conf:          conf,
		connections:   make(map[uint64]RouteInterface),
		listeners:     make([]*connmgmt.ConnMgmtListen, 0),
		frameTracker:  *packet.NewFrameTracker(cx),
		addConnChan:   make(chan *connmgmt.ConnMgmtStruct),
		writeConnChan: make(chan *packet.ConnFrame, ethChanDepth),
		readConnChan:  make(chan *packet.ConnFrame, ethChanDepth),
	}

	return r, e
}

//
// Start() connection
// This will launch the connection process
//
func (r *RouteStruct) Start() {

	go r.goReadConnChan()
	go r.goWriteConn()
	go r.goAddConnection()

	r.startListeners()
	r.startClientConnections()

}

//
// Stop() connection
// Orderly shutdown of all connections
//
func (r *RouteStruct) Stop() {

	r.ctx.Cancel()

	for _, j := range r.listeners {
		j.Stop()
	}

	for _, j := range r.connections {
		j.Stop()
	}
}

//
//
//
func (r *RouteStruct) startClientConnections() {

	for _, j := range r.conf.Targ {
		if j.Active {
			connection, e := connmgmt.NewClient(r.ctx, j)
			if e != nil {
				r.ctx.Logf(ctx.LogLevelPanic, "error:%s", e)
			}
			r.addConnChan <- connection
		}
	}
}

//
//
//
func (r *RouteStruct) startListeners() {

	for _, j := range r.conf.List {
		if j.Active {
			listen, e := connmgmt.NewListen(r.ctx, j, r.addConnChan)
			if e != nil {
				r.ctx.Logf(ctx.LogLevelPanic, "server.New() error:%s", e)
			}
			r.listeners = append(r.listeners, listen)
		}
	}

	for _, j := range r.listeners {
		j.Start()
	}

}

//
// goReadConnChan()
// Read packets from the lower stack
// Determine if the packet has passed already
// Pass it to the transit layer if not
//
func (r *RouteStruct) goReadConnChan() {

	for {
		select {
		case <-r.ctx.DoneChan():
			return
		case conn := <-r.readConnChan:
			if r.frameTracker.PassFrame(conn.ID) {
				r.readEthChan <- conn.Frame
			} else {
				// Stale Frame
				_ = conn
			}
		}
	}
}

//
// ReadEthChan()
// Let the transit function read the eth packets coming in over a connection
//
func (r *RouteStruct) ReadEthChan() <-chan *packet.EthFrame {
	return r.readEthChan
}

//
// Write()
// Get EthFrame packet to write
// Wrap in ConnFrame for tracking
// Wrap in standard ProtoHeader (what is actually sent over the Conn)
// Push into writeEthChan
// Send an ether packet to one or more connections
//
func (r *RouteStruct) Write(eth *packet.EthFrame) {
	conn := packet.NewConnFrame(eth)
	// comm := packet.NewHeaderV1Payload(conn)
	r.writeConnChan <- conn

}

//
// goWriteConn()
// Figure out which connection(s) to send packs over
//
func (r *RouteStruct) goWriteConn() {

	for {
		select {
		case <-r.ctx.DoneChan():
			return
		case proto := <-r.writeConnChan:
			//
			// Lots of calculations here... which Conn circuits do we send the packets on?
			//
			// For now, write to all online Conn channels

			for _, j := range r.connections {
				if j.Online() {

					//ch := packet.NewHeaderV1Payload(proto)
					//j.Write(ch)
					j.Write(proto)
				}
			}
		}
	}
}

//
//
//
func (r *RouteStruct) goAddConnection() {

	for {
		var conn interface{}

		select {
		case <-r.ctx.DoneChan():
			return
		case conn = <-r.addConnChan:

			connID := <-counterConnectionChan
			var ri RouteInterface

			switch conn := conn.(type) {
			case *connmgmt.ConnMgmtStruct:
				ri = conn
			default:
				r.ctx.Logf(ctx.LogLevelPanic, "Type not supported:%t", conn)
			}

			r.connections[connID] = ri
			ri.Start()
		}
	}
}

//
//
//
func (r *RouteStruct) newConnFrame(frame *packet.EthFrame) (cf *packet.ConnFrame) {
	r.ctx.LogLocation()
	cf = &packet.ConnFrame{
		ID:    packet.ConnFrameID(<-counterConnChanIDChan),
		Frame: frame,
	}
	return cf
}
