package tcp

import (
	"errors"
	"net"

	"github.com/seanmcadam/octovpn/ctx"
)

const TCPAcceptChannelLen int = 2
const TCPReadChannelLen int = 10

const TCPReadBufferSize int = 1600

var ErrSocTCPListerClosed = errors.New("socket TCP: Listener closed")
var ErrSocTCPClosed = errors.New("socket TCP: socket closed")
var ErrSocTCPBadNetwork = errors.New("socket TCP: bad network protocol")
var ErrSocTCPEmtpyWriteBuf = errors.New("socket TCP: Emtpy Write Buffer")

type ListenTCPStruct struct {
	listener net.Listener
	ctx      *ctx.Ctx
	accept   chan interface{}
}

type SocketTCPStruct struct {
	socket   net.Conn
	ctx      *ctx.Ctx
	readchan chan readStruct
}

type readStruct struct {
	err   error
	count int
	buf   []byte
}
