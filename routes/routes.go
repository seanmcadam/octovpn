package routes

import (
	"net"

	"github.com/seanmcadam/octovpn/connections"
	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/octoconfig"
	"github.com/seanmcadam/octovpn/octolib"
	"github.com/seanmcadam/octovpn/packet"
	"github.com/seanmcadam/octovpn/pinger"
	"github.com/seanmcadam/octovpn/server"
)

// Package to manage all routes to the target vpn system
// New()(c *ConnectionStruct)
// (c *ConnectionStruct) Start()
// (c *ConnectionStruct) Stop()
//
// Send filtered Eth packets from Transit module
// (c *ConnectionStruct) Write(eth *packet.EthFrame)
//
// Read filtered Eth Packets in Transit module
// (c *ConnectionStruct) ReadEthChan() <-chan *packet.EthFrame
//
// Handler for new incoming connections from listeners
// (c *ConnectionStruct) goAddConnection()
//
// Collect and handle incoming packets from the connection layer
// (c *ConnectionStruct) goReadVPNChan()
//
// Collect and handle incoming packets from the connection layer
// (c *ConnectionStruct) goWriteVPN()
//
// Create wrapper for Eth Frame
// (c *ConnectionStruct) newConnFrame(frame *packet.EthFrame) (cf *packet.ConnFrame)
//

//type ConnState string
//
//const ConnStateNew = "new"
//const ConnStateRunning = "running"
//const ConnStateError = "error"
//const ConnStateClosed = "closed"

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
	Loss() pinger.Loss           // Loss calculation 0-1000 - 0 = best
	Latency() pinger.Latency     // Latency calulation 0-1000 - 0 = best
	Deviation() pinger.Deviation // Deviation calulation 0-1000 - 0 = best
}

//
// Used for both originating and recieved sockets
//
type RouteStruct struct {
	//
	// Context structure
	//
	ctx *ctx.Ctx
	//
	// Configuration
	//
	conf octoconfig.ConfigV1
	//
	// All connections from both client and server
	//
	connections map[uint64]RouteInterface
	//
	// Listener Sockets on this device
	//
	listeners []*server.ListenerStruct
	//
	// Tracks ConnFrame packets to make sure that only the first frame recieved passes
	//
	frameTracker packet.ConnFrameTrackerStruct
	//
	// The channel passes new accepts from listener connections
	//
	addVPNChan chan interface{}
	//
	// Read channel for Transit layer, only one unique EthFrame ID is passed up
	//
	readEthChan chan *packet.EthFrame
	//
	// Read channel for VPN layer
	//
	readVPNChan chan *packet.ConnFrame
	//
	// Write Channel to connections, pushed into the processing section of sending over which VPNs
	//
	writeVPNChan chan *packet.ConnFrame
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
		ctx:          cx,
		conf:         conf,
		connections:  make(map[uint64]RouteInterface),
		listeners:    make([]*server.ListenerStruct, 0),
		frameTracker: *packet.NewFrameTracker(cx),
		addVPNChan:   make(chan interface{}),
		writeVPNChan: make(chan *packet.ConnFrame, ethChanDepth),
		readVPNChan:  make(chan *packet.ConnFrame, ethChanDepth),
	}

	return r, e
}

//
// Start() connection
// This will launch the connection process
//
func (r *RouteStruct) Start() {

	go r.goReadVPNChan()
	go r.goWriteVPN()
	go r.goAddConnection()

	r.startClientConnections()
	r.startListeners()

}

//
//
//
func (r *RouteStruct) startClientConnections() {

	for _, j := range r.conf.Targ {
		if j.Active {
			connection, e := connections.New(r.ctx, j)
			if e != nil {
				r.ctx.Logf(ctx.LogLevelPanic, "error:%s", e)
			}
			r.addVPNChan <- connection
		}
	}
}

//
//
//
func (r *RouteStruct) startListeners() {

	for _, j := range r.conf.List {
		if j.Active {
			listen, e := server.NewListener(r.ctx, j, r.addVPNChan, r.readVPNChan)
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
// Stop() connection
// Orderly shutdown of all connections
//
func (r *RouteStruct) Stop() {

	r.ctx.Cancel()

	for _, j := range r.listeners {
		j.Stop()
	}

	for _, j := range r.vpn {
		j.Stop()
	}
}

//
// goReadVPNChan()
// Read packets from the lower stack
// Determine if the packet has passed already
// Pass it to the transit layer if not
//
func (r *RouteStruct) goReadVPNChan() {

	for {
		select {
		case <-r.ctx.DoneChan():
			return
		case conn := <-r.readVPNChan:
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
// Wrap in standard ProtoHeader (what is actually sent over the VPN)
// Push into writeEthChan
// Send an ether packet to one or more connections
//
func (r *RouteStruct) Write(eth *packet.EthFrame) {
	conn := packet.NewConnFrame(eth)
	r.writeVPNChan <- conn

}

//
//
//
func (r *RouteStruct) goWriteVPN() {

	select {
	case <-r.ctx.DoneChan():
		return
	case proto := <-r.writeVPNChan:
		//
		// Lots of calculations here... which VPN circuits do we send the packets on?
		//
		// For now, write to all online VPN channels

		for _, j := range r.vpn {
			if j.Online() {
				j.Write(proto)
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
		case conn = <-r.addVPNChan:
		}

		connID := <-counterConnectionChan
		var ri RouteInterface

		switch conn.(type) {
		case *connections.ConnectionsStruct:
			ri = conn.(*connections.ConnectionsStruct)
		default:
			r.ctx.Logf(ctx.LogLevelPanic, "Type not supported:%t", conn)
		}

		r.connections[connID] = ri
		ri.Start()
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
