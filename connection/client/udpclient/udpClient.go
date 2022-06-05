package udpclient

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"net"

	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/octoconfig"
)

type UDPClientStruct struct {
	ctx      ctx.Ctx
	conn     net.Conn
	writeBuf *bufio.Writer
	decoder  *gob.Decoder
	encoder  *gob.Encoder
	readChan chan interface{}
}

var ErrBadProtocol = errors.New("bad UDP Protocol")

//
// New()
// Create a new UDP client connection
//
func New(cx ctx.Ctx, protocol octoconfig.ConnectionProtocol, host string, port uint16, readchan chan interface{}) (udp *UDPClientStruct, e error) {

	udp = &UDPClientStruct{
		ctx:      cx,
		readChan: readchan,
	}

	switch protocol {
	case octoconfig.UDP:
	case octoconfig.UDP4:
	case octoconfig.UDP6:
	default:
		e = ErrBadProtocol
		return nil, e
	}

	udpaddr, e := net.ResolveUDPAddr(string(protocol), fmt.Sprintf("%s:%d", host, port))
	if e != nil {
		return nil, e
	}

	udp.conn, e = net.DialUDP(string(protocol), nil, udpaddr)
	if e != nil {
		return nil, e
	}

	udp.writeBuf = bufio.NewWriter(udp.conn)
	udp.encoder = gob.NewEncoder(udp.writeBuf)
	udp.decoder = gob.NewDecoder(udp.conn)

	return udp, e
}

//
//
//
func (u *UDPClientStruct) Start() {
	go u.goRecv()
}

//
//
//
func (u *UDPClientStruct) Send(data interface{}) (e error) {
	e = u.encoder.Encode(data)
	if e == nil {
		u.writeBuf.Flush()
	}
	return e
}

func (u *UDPClientStruct) goRecv() {

	for {
		select {
		case <-u.ctx.DoneChan():
			return
		default:
			var data interface{}
			e := u.decoder.Decode(&data)
			if e != nil {
				u.ctx.Logf(ctx.LogLevelPanic, "Decode() returned error %s", e)
			}

			u.readChan <- &data

		}
	}
}

func (u *UDPClientStruct) LocalAddr() net.Addr {
	return u.conn.LocalAddr()
}

func (u *UDPClientStruct) RemoteAddr() net.Addr {
	return u.conn.RemoteAddr()
}
