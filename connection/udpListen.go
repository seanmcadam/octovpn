package connection

import (
	"errors"
	"fmt"
	"net"

	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/octoconfig"
)

var ErrBadProtocol = errors.New("bad UDP Protocol")

type UDPListenStruct struct {
	ctx      *ctx.Ctx
	listener *net.UDPConn
	config   *octoconfig.ConfigListen
	//readProtoChan chan *packet.ProtoHeader
	acceptChan chan interface{}
}

func NewUDPListen(cx *ctx.Ctx, t *octoconfig.ConfigListen) (listen *UDPListenStruct, e error) {

	cx = cx.NewWithCancel()

	listen = &UDPListenStruct{
		ctx:    cx,
		config: t,
		//readProtoChan: ReadProtoChan,
		acceptChan: make(chan interface{}),
	}

	protocol := t.Protocol
	ip := t.IP
	port := t.Port

	switch protocol {
	case octoconfig.UDP:
	case octoconfig.UDP4:
	case octoconfig.UDP6:
	default:
		e = ErrBadProtocol
		return nil, e
	}

	udpaddr, e := net.ResolveUDPAddr(string(protocol), fmt.Sprintf("%s:%d", ip, port))
	if e != nil {
		return nil, e
	}

	listen.listener, e = net.ListenUDP(string(protocol), udpaddr)
	if e != nil {
		return nil, e
	}

	return listen, e
}

func (ul *UDPListenStruct) Start() {
	go ul.goRunReader()
}

func (ul *UDPListenStruct) Stop() {
	ul.listener.Close()
	ul.ctx.Cancel()
}

func (ul *UDPListenStruct) LocalAddr() (addr net.Addr) {
	return ul.listener.LocalAddr()
}

func (ul *UDPListenStruct) RemoteAddr() (addr net.Addr) {
	return ul.listener.RemoteAddr()
}

func (ul *UDPListenStruct) AcceptChan() chan interface{} {
	return ul.acceptChan
}

//
// goRunReader()
// Returns a UDPStruct to the connection channel
//
func (ul *UDPListenStruct) goRunReader() {

	//	for {
	//
	//		select {
	//		case <-ul.ctx.DoneChan():
	//			return
	//		default:
	//		}
	//
	//		var buf []byte
	//		_, addr, e := ul.listener.ReadFromUDP(buf)
	//		if e != nil {
	//			ul.ctx.Logf(ctx.LogLevelPanic, " error %s", e)
	//		}
	//
	//		i, ok := ul.connAddr[addr]
	//		if ok {
	//			i.bufChan <- buf
	//		} else {
	//
	//			ctx := ul.ctx.NewWithCancel()
	//			udpConn := &UDPServerStruct{
	//				ctx:        ctx,
	//				udpConn:    ul.listener,
	//				remoteAddr: addr,
	//				readChan:   ul.readChan,
	//				bufChan:    make(chan []byte),
	//			}
	//
	//			udpConn.writeBuf = bufio.NewWriter(udpConn.udpConn)
	//			udpConn.encoder = gob.NewEncoder(udpConn.writeBuf)
	//
	//			ul.newConnChan <- udpConn
	//		}
	//	}
}
