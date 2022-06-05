package tcpclient

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"net"

	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/octoconfig"
)

type TCPClientStruct struct {
	ctx      ctx.Ctx
	conn     net.Conn
	writeBuf *bufio.Writer
	decoder  *gob.Decoder
	encoder  *gob.Encoder
	readChan chan interface{}
}

var ErrBadProtocol = errors.New("bad TCP Protocol")

//
// New()
// Create a new TCP client connection
//
func New(cx ctx.Ctx, protocol octoconfig.ConnectionProtocol, host string, port uint16, readchan chan interface{}) (tcp *TCPClientStruct, e error) {

	tcp = &TCPClientStruct{
		ctx:      cx,
		readChan: readchan,
	}

	switch protocol {
	case octoconfig.TCP:
	case octoconfig.TCP4:
	case octoconfig.TCP6:
	default:
		e = ErrBadProtocol
		return nil, e
	}

	tcpaddr, e := net.ResolveTCPAddr(string(protocol), fmt.Sprintf("%s:%d", host, port))
	if e != nil {
		return nil, e
	}

	tcp.conn, e = net.DialTCP(string(protocol), nil, tcpaddr)
	if e != nil {
		return nil, e
	}

	tcp.writeBuf = bufio.NewWriter(tcp.conn)

	return tcp, e
}

func (t *TCPClientStruct) Start() {
	go t.goRecv()
}

func (t *TCPClientStruct) Send(data interface{}) (e error) {
	e = t.encoder.Encode(data)
	if e == nil {
		t.writeBuf.Flush()
	}
	return e
}

func (t *TCPClientStruct) goRecv() {

	for {
		select {
		case <-t.ctx.DoneChan():
			return
		default:
			var data interface{}
			e := t.decoder.Decode(&data)
			if e != nil {
				t.ctx.Logf(ctx.LogLevelPanic, "Decode() returned error %s", e)
			}

			t.readChan <- &data

		}
	}
}

func (t *TCPClientStruct) LocalAddr() net.Addr {
	return t.conn.LocalAddr()
}

func (t *TCPClientStruct) RemoteAddr() net.Addr {
	return t.conn.RemoteAddr()
}
