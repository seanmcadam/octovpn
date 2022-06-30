package connection

import (
	"fmt"
	"net"

	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/octoconfig"
)

type TCPListenStruct struct {
	ctx        *ctx.Ctx
	listener   net.Listener
	config     *octoconfig.ConfigListen
	acceptChan chan interface{}
	//readProtoChan chan *packet.ProtoHeader
}

//
//
//
func NewTCPListen(cx *ctx.Ctx, t *octoconfig.ConfigListen) (listen *TCPListenStruct, e error) {

	cx = cx.NewWithCancel()

	listen = &TCPListenStruct{
		ctx:        cx,
		config:     t,
		acceptChan: make(chan interface{}),
	}

	protocol := t.Protocol
	ip := t.IP
	port := t.Port

	switch protocol {
	case octoconfig.TCP:
	case octoconfig.TCP4:
	case octoconfig.TCP6:
	default:
		e = ErrBadProtocol
		return nil, e
	}

	tcpaddr, e := net.ResolveTCPAddr(string(protocol), fmt.Sprintf("%s:%d", ip, port))
	if e != nil {
		return nil, e
	}

	listen.listener, e = net.ListenTCP(string(protocol), tcpaddr)
	if e != nil {
		return nil, e
	}

	return listen, e
}

func (tl *TCPListenStruct) Start() {
	go tl.goRunListener()
}

func (tl *TCPListenStruct) Stop() {
	tl.listener.Close()
	tl.ctx.Cancel()
}

func (tl *TCPListenStruct) LocalAddr() (addr net.Addr) {
	return tl.listener.Addr()
}

func (tl *TCPListenStruct) AcceptChan() chan interface{} {
	return tl.acceptChan
}

//
// goRunListener()
// Returns a TCPStruct to the connection channel
//
func (tl *TCPListenStruct) goRunListener() {

	for {

		select {
		case <-tl.ctx.DoneChan():
			return
		default:
		}

		netConn, e := tl.listener.Accept()
		if e != nil {
			tl.ctx.Logf(ctx.LogLevelPanic, " error %s", e)
		}

		cx := tl.ctx.NewWithCancel()
		connStruct := newStruct(cx, netConn, nil)

		tl.acceptChan <- connStruct
	}
}
