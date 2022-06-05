package tcpserver

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"net"

	"github.com/seanmcadam/octovpn/connection"
	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/octoconfig"
	"github.com/seanmcadam/octovpn/octolib"
)

type listenTCPStruct struct {
	ctx           *ctx.Ctx
	protocol      octoconfig.ConnectionProtocol
	ip            string
	port          uint16
	mtu           uint16
	listener      net.Listener
	addConnection chan interface{}
}

type TCPStruct struct {
	ctx               *ctx.Ctx
	ID                uint64
	protocol          octoconfig.ConnectionProtocol
	ip                string
	port              uint16
	mtu               uint16
	tcpConn           net.Conn
	readConnFrameChan chan *connection.ConnFrame
	writeBuf          io.Writer
	decoder           *gob.Decoder
	encoder           *gob.Encoder
	pinger            *Pinger
}

func newListenTCP(cx *ctx.Ctx, l *octoconfig.ConfigListen, addConnChan chan interface{}) (listen *listenTCPStruct, e error) {

	cx = cx.NewWithCancel()

	listen = &listenTCPStruct{
		ctx:           cx,
		protocol:      l.Protocol,
		ip:            l.IP,
		port:          l.Port,
		mtu:           l.MTU,
		addConnection: addConnChan,
	}

	address := fmt.Sprintf("%s:%d", l.IP, l.Port)

	switch l.Protocol {
	case "tcp":
		fallthrough
	case "tcp4":
		fallthrough
	case "tcp6":

		listen.listener, e = net.Listen(string(l.Protocol), address)

	default:
		cx.Logf(ctx.LogLevelPanic, "Unhandled Protocol: %s", l.Protocol)
	}

	if e != nil {
		return nil, e
	}

	go listen.goRunListener()

	return listen, e
}

//
// goRunListener()
// Returns a TCPStruct to the connection channel
//
func (t *listenTCPStruct) goRunListener() {

	for {
		netConn, e := t.listener.Accept()
		if e != nil {
			t.ctx.Logf(ctx.LogLevelPanic, " error %s", e)
		}

		addr := netConn.RemoteAddr()
		ipport := addr.String()

		ip, port, e := octolib.SplitAddr(ipport)
		if e != nil {
			t.ctx.Logf(ctx.LogLevelPanic, " SplitAddr error %s", e)
		}

		ctx := t.ctx.NewWithCancel()
		tcpConn := &TCPStruct{
			ctx:      ctx,
			protocol: t.protocol,
			ip:       ip,
			port:     port,
			mtu:      t.mtu,
			tcpConn:  netConn,
		}

		pinger := NewPinger(ctx, 200)
		tcpConn.pinger = pinger

		tcpConn.writeBuf = bufio.NewWriter(tcpConn.tcpConn)
		tcpConn.encoder = gob.NewEncoder(tcpConn.writeBuf)
		tcpConn.decoder = gob.NewDecoder(tcpConn.tcpConn)

		go tcpConn.goRunReadFrom()

		// Fire up monitor go routine
		// Fire up reader go routine

		select {
		case <-t.ctx.DoneChan():
			return
		default:
			t.addConnection <- tcpConn
		}

		go tcpConn.goSendPinger()
		tcpConn.pinger.Run()

	}
}

//
//
//
func (t *TCPStruct) goRunReadFrom() {
	for {
		select {
		case <-t.ctx.DoneChan():
			return
		default:
			var header ProtoHeader
			e := t.decoder.Decode(&header)
			if e != nil {
				t.ctx.Logf(ctx.LogLevelPanic, "Decode() returned error %s", e)
			}

			switch header.Payload.(type) {
			case connection.ConnFrame:
				frame := header.Payload.(connection.ConnFrame)
				t.readConnFrameChan <- &frame

			case Ping:
				ping := header.Payload.(Ping)
				t.pinger.GotPing(&ping)

			case Pong:
				pong := header.Payload.(Pong)
				t.pinger.GotPong(&pong)

			default:
				t.ctx.Logf(ctx.LogLevelPanic, "Unhandled type Recieved:%v", header.Payload)
			}
		}
	}

}

//
//
//
func (t *TCPStruct) Recv() (connFrame *connection.ConnFrame, e error) {

	connFrame = <-t.readConnFrameChan
	return connFrame, e
}

//
// SendConnFrame()
//
func (t *TCPStruct) SendConnFrame(connFrame *connection.ConnFrame) (e error) {
	header, e := NewProtoHeader(tcpHeaderSignature, connFrame)
	if e != nil {
		t.ctx.Logf(ctx.LogLevelPanic, " newHeaderFrame error:%s", e)
	}

	e = t.send(header)

	return e
}

//
// goSendPinger()
// Read pinger send channel, wrap and send the ping and pong structs
//
func (t *TCPStruct) goSendPinger() {

	var pinger interface{}
	for {
		select {
		case <-t.ctx.DoneChan():
			return
		case pinger = <-t.pinger.SendChan():
			switch pinger.(type) {
			case *Ping:
			case *Pong:
			default:
				t.ctx.Logf(ctx.LogLevelPanic, "Bad type:%t", pinger)
			}

			header, e := NewProtoHeader(tcpHeaderSignature, pinger)
			if e != nil {
				t.ctx.Logf(ctx.LogLevelPanic, " newHeaderFrame() error:%s", e)
			}
			e = t.send(header)
			if e != nil {
				t.ctx.Logf(ctx.LogLevelPanic, " send() error:%s", e)
			}
		}
	}
}

//
//
//
func (t *TCPStruct) send(header *ProtoHeader) (e error) {
	e = t.encoder.Encode(header)
	return e
}

//
//
//
func (t *TCPStruct) State() (c ConnState) {
	return c
}

//
//
//
func (t *TCPStruct) Status() (c ConnStatus) {
	return c
}

//
//
//
func (t *TCPStruct) Protocol() (p octoconfig.ConnectionProtocol) {
	return t.protocol
}

//
//
//
func (t *TCPStruct) HostIP() (h string) {
	return t.ip
}

//
//
//
func (t *TCPStruct) Port() (h uint16) {
	return t.port
}

//
//
//
func (t *TCPStruct) MTU() (m uint16) {
	return t.mtu
}
