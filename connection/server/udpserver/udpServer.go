package udpserver

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"net"

	"github.com/seanmcadam/octovpn/connection"
	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/octoconfig"
)

type listenUDPStruct struct {
	ctx               *ctx.Ctx
	ID                uint64
	protocol          octoconfig.ConnectionProtocol
	ip                string
	port              uint16
	mtu               uint16
	udpConn           *net.UDPConn
	readConnFrameChan chan *connection.ConnFrame
	writeBuf          io.Writer
	decoder           *gob.Decoder
	encoder           *gob.Encoder
}

func newListenUDP(cx *ctx.Ctx, l *octoconfig.ConfigListen, readConnFrame chan *connection.ConnFrame) (listen *listenUDPStruct, e error) {

	cx = cx.NewWithCancel()
	var protocol string

	listen = &listenUDPStruct{
		ctx:               cx,
		protocol:          l.Protocol,
		ip:                l.IP,
		port:              l.Port,
		mtu:               l.MTU,
		readConnFrameChan: readConnFrame,
	}

	switch l.Protocol {
	case "udp":
		protocol = "udp"
	case "udp4":
		protocol = "udp4"
	case "udp6":
		protocol = "udp6"
	default:
		cx.Logf(ctx.LogLevelPanic, "Unhandled Protocol: %s", l.Protocol)
	}

	address, e := net.ResolveUDPAddr(protocol, fmt.Sprintf("%s:%d", l.IP, l.Port))
	if e != nil {
		cx.Logf(ctx.LogLevelPanic, "Resolve Addr error: %s", e)
	}

	listen.udpConn, e = net.ListenUDP(protocol, address)

	if e != nil {
		return nil, e
	}

	listen.writeBuf = bufio.NewWriter(listen.udpConn)
	listen.encoder = gob.NewEncoder(listen.writeBuf)
	listen.decoder = gob.NewDecoder(listen.udpConn)

	go listen.goRunReadFrom()

	return listen, e
}

//
// goRunReadFrom()
// reads incoming packets
// decodes them into Listener listenUDPStruct
// Error checks
// Updates statistics
// extracts ConnFrame and put from on the readChan
//
func (l *listenUDPStruct) goRunReadFrom() {

	for {

		select {
		case <-l.ctx.DoneChan():
			return
		default:
			var header ProtoHeader
			e := l.decoder.Decode(&header)
			if e != nil {
				l.ctx.Logf(ctx.LogLevelPanic, "Decode() returned error %s", e)
			}

			switch header.Payload.(type) {
			case connection.ConnFrame:
				frame := header.Payload.(connection.ConnFrame)
				l.readConnFrameChan <- &frame

			case Ping:
				//ping := header.Payload.(Ping)
				//pong := l.pinger.GotPing()
				//pongHeader, e := NewProtoHeader(udpHeaderSignature, pong)
				//if e != nil {
				//	l.ctx.Logf(ctx.LogLevelPanic, "newHeaderProto() error %s", e)
				//}

				//e = l.send(pongHeader)
				//if e != nil {
				//	l.ctx.Logf(ctx.LogLevelPanic, "send() error %s", e)
				//}

			case Pong:
				pong := header.Payload.(Pong)

				l.ctx.Logf(ctx.LogLevelTrace, "Pong Recieved:%v", pong)

			default:
				l.ctx.Logf(ctx.LogLevelPanic, "Unhandled type Recieved:%v", header.Payload)
			}
		}
	}
}

//
//
//
func (l *listenUDPStruct) Recv() (conFrame *connection.ConnFrame, e error) {

	frame := <-l.readConnFrameChan
	return frame, e
}

//
//
//
func (l *listenUDPStruct) Send(connFrame *connection.ConnFrame) (e error) {
	header, e := NewProtoHeader(udpHeaderSignature, connFrame)
	if e != nil {
		l.ctx.Logf(ctx.LogLevelPanic, " newHeaderFrame error:%s", e)
	}

	e = l.send(header)
	return e
}

//
//
//
func (l *listenUDPStruct) send(header *ProtoHeader) (e error) {
	e = l.encoder.Encode(header)
	return e
}

//
//
//
func (l *listenUDPStruct) State() (c ConnState) {
	return c
}

//
//
//
func (l *listenUDPStruct) Status() (c ConnStatus) {
	return c
}

//
//
//
func (l *listenUDPStruct) Protocol() (p octoconfig.ConnectionProtocol) {
	return l.protocol
}

//
//
//
func (l *listenUDPStruct) HostIP() (h string) {
	return l.ip
}

//
//
//
func (l *listenUDPStruct) Port() (h uint16) {
	return l.port
}

//
//
//
func (l *listenUDPStruct) MTU() (m uint16) {
	return l.mtu
}
