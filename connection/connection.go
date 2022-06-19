package connection

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"

	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/octoconfig"
	"github.com/seanmcadam/octovpn/packet"
)

type ConnectionStruct struct {
	ctx          *ctx.Ctx
	conn         net.Conn
	local        net.Addr
	remote       net.Addr
	protocol     octoconfig.ConnectionProtocol
	writeBuf     *bufio.Writer
	decoder      *gob.Decoder
	encoder      *gob.Encoder
	readChan     chan *packet.ProtoHeader
	clientConfig *octoconfig.ConfigTarget
}

// func init() {
//
// }

//
//
//
func NewConn(cx *ctx.Ctx, config *octoconfig.ConfigTarget) (conn *ConnectionStruct, e error) {

	protocol := config.Protocol
	port := config.Port
	host := config.Hostname
	//mtu := config.MTU

	var netconn net.Conn

	switch protocol {
	case octoconfig.TCP:
		fallthrough
	case octoconfig.TCP4:
		fallthrough
	case octoconfig.TCP6:
		netconn, e = dialTCP(protocol, host, port)
		conn = NewStruct(cx, netconn, config)

	case octoconfig.UDP:
		fallthrough
	case octoconfig.UDP4:
		fallthrough
	case octoconfig.UDP6:
		netconn, e = dialUDP(protocol, host, port)
		conn = NewStruct(cx, netconn, config)
	default:
		e = octoconfig.ErrBadProtocol
		return nil, e
	}

	return conn, e
}

//
//
//
func NewStruct(cx *ctx.Ctx, net net.Conn, config *octoconfig.ConfigTarget) (conn *ConnectionStruct) {
	readchan := make(chan *packet.ProtoHeader)
	conn = &ConnectionStruct{
		ctx:          cx,
		conn:         net,
		readChan:     readchan,
		local:        net.LocalAddr(),
		remote:       net.RemoteAddr(),
		writeBuf:     bufio.NewWriter(conn.conn),
		encoder:      gob.NewEncoder(conn.writeBuf),
		decoder:      gob.NewDecoder(conn.conn),
		clientConfig: config,
	}
	return conn
}

func (conn *ConnectionStruct) Protocol() octoconfig.ConnectionProtocol { return conn.protocol }
func (conn *ConnectionStruct) LocalAddr() net.Addr                     { return conn.local }
func (conn *ConnectionStruct) RemoteAddr() net.Addr                    { return conn.remote }
func (conn *ConnectionStruct) ReadChan() chan *packet.ProtoHeader      { return conn.readChan }

//
//
//
func (conn *ConnectionStruct) Start() {
	go conn.goRead()
}

//
//
//
func (conn *ConnectionStruct) Stop() {
	conn.ctx.Cancel()
}

//
// Write()
// Encode ProtoHeaders and send them over the network
//
func (conn *ConnectionStruct) Write(p *packet.ProtoHeader) (e error) {
	e = conn.encoder.Encode(p)
	if e == nil {
		conn.writeBuf.Flush()
	}
	return e
}

//
// goRead()
// Read and decode incoming packets from the network connection
//
func (conn *ConnectionStruct) goRead() {
	for {
		select {
		case <-conn.ctx.DoneChan():
			return
		default:
			var data interface{}
			var proto *packet.ProtoHeader
			e := conn.decoder.Decode(&data)
			if e != nil {
				conn.ctx.Logf(ctx.LogLevelPanic, "Decode() returned error %s", e)
			}

			switch data := data.(type) {
			case packet.ProtoHeader:
				proto = &data
			case *packet.ProtoHeader:
				proto = data
			default:
				conn.ctx.Logf(ctx.LogLevelPanic, "Decode() returned bad type %t", data)
			}
			conn.readChan <- proto
		}
	}
}

//
//
//
func dialTCP(protocol octoconfig.ConnectionProtocol, host string, port uint16) (conn *net.TCPConn, e error) {

	var remoteaddr *net.TCPAddr = nil
	var localaddr *net.TCPAddr = nil

	remoteaddr, e = net.ResolveTCPAddr(string(protocol), fmt.Sprintf("%s:%d", host, port))
	if e != nil {
		return nil, e
	}

	conn, e = net.DialTCP(string(protocol), localaddr, remoteaddr)

	return conn, e
}

//
//
//
func dialUDP(protocol octoconfig.ConnectionProtocol, host string, port uint16) (conn *net.UDPConn, e error) {

	var remoteaddr *net.UDPAddr = nil
	var localaddr *net.UDPAddr = nil

	remoteaddr, e = net.ResolveUDPAddr(string(protocol), fmt.Sprintf("%s:%d", host, port))
	if e != nil {
		return nil, e
	}

	conn, e = net.DialUDP(string(protocol), localaddr, remoteaddr)

	return conn, e
}
