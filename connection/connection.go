package connection

import (
	"errors"
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
	clientConfig *octoconfig.ConfigTarget
	readChan     chan *ConnReadStruct
}

type ConnReadStruct struct {
	e      error
	count  int
	header *packet.CommonHeader
}

func (c *ConnReadStruct) Err() error {
	return c.e
}

func (c *ConnReadStruct) Count() int {
	return c.count
}

func (c *ConnReadStruct) Header() *packet.CommonHeader {
	return c.header
}

//
// NewConn()
// Create a new ConnectionStruct based on a target config
//
func NewConn(cx *ctx.Ctx, config *octoconfig.ConfigTarget) (conn *ConnectionStruct, e error) {

	protocol := config.Protocol
	port := config.Port
	host := config.Hostname
	//mtu := config.MTU

	conn = nil

	var netconn net.Conn

	switch protocol {
	case octoconfig.TCP:
		fallthrough
	case octoconfig.TCP4:
		fallthrough
	case octoconfig.TCP6:
		netconn, e = dialTCP(protocol, host, port)
	case octoconfig.UDP:
		fallthrough
	case octoconfig.UDP4:
		fallthrough
	case octoconfig.UDP6:
		netconn, e = dialUDP(protocol, host, port)
	default:
		e = octoconfig.ErrBadProtocol
	}

	if e == nil {
		conn = newStruct(cx, netconn, config)
	}

	return conn, e
}

//
// NewStruct()
// Create a new ConnectionStruct based on net.Conn
//
func newStruct(cx *ctx.Ctx, net net.Conn, config *octoconfig.ConfigTarget) (conn *ConnectionStruct) {

	cx.LogLocation()

	readchan := make(chan *ConnReadStruct)
	conn = &ConnectionStruct{
		ctx:          cx,
		conn:         net,
		readChan:     readchan,
		clientConfig: config,
	}
	conn.local = net.LocalAddr()
	conn.remote = net.RemoteAddr()

	return conn
}

func (conn *ConnectionStruct) Protocol() octoconfig.ConnectionProtocol { return conn.protocol }
func (conn *ConnectionStruct) LocalAddr() net.Addr                     { return conn.local }
func (conn *ConnectionStruct) RemoteAddr() net.Addr                    { return conn.remote }
func (conn *ConnectionStruct) ReadChan() chan *ConnReadStruct          { return conn.readChan }

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
	e := conn.conn.Close()
	if e != nil {
		conn.ctx.Logf(ctx.LogLevelError, "Close error:%s", e)
	}
	conn.ctx.Cancel()
}

//
// Write()
// Encode ProtoHeaders and send them over the network
//
func (conn *ConnectionStruct) Write(ch *packet.CommonHeader) (count int, e error) {
	return conn.conn.Write(ch.ToByte())
}

//
// goRead()
// Read and decode incoming packets from the network connection
// read the first packets to determine the packet type and lenght
//
func (conn *ConnectionStruct) goRead() {
	for {
		select {
		case <-conn.ctx.DoneChan():
			return
		default:
			packetheader := make([]byte, packet.PacketHeaderSize)
			count, e := conn.conn.Read(packetheader)
			if e != nil {
				conn.ctx.Logf(ctx.LogLevelPanic, "Read header error:%e", e)
			}
			if count != packet.PacketHeaderSize {
				conn.ctx.Logf(ctx.LogLevelPanic, "Bad packetheadersize count:%d", count)
			}
			if packet.ValidHeaderV1 != packet.PacketBlockBuf(packetheader) {
				conn.ctx.Logf(ctx.LogLevelPanic, "Bad PacketHeader")
			}
			if packet.VersionV1 != packet.PacketVersionBuf(packetheader) {
				conn.ctx.Logf(ctx.LogLevelPanic, "Bad PacketHeader")
			}

			ch, size := packet.NewHeaderRead(packetheader)

			payload := make([]byte, size)
			count, e = conn.conn.Read(payload)
			if e != nil {
				conn.ctx.Logf(ctx.LogLevelPanic, "Read payload error:%e", e)
			}
			if count != packet.PacketHeaderSize {
				str := fmt.Sprintf("Bad payload size count:%d", count)
				conn.ctx.Logf(ctx.LogLevelPanic, str)
				e = errors.New(str)
			}

			ch.AddPayload(payload)

			rs := &ConnReadStruct{
				e:      e,
				count:  int(size) + packet.PacketHeaderSize,
				header: ch,
			}

			conn.readChan <- rs
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
