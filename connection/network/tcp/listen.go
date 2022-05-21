package tcp

import (
	"context"
	"net"

	"github.com/seanmcadam/octovpn/ctx"
)

//
// NewListener
// Accept
// Close
//
//

//
// Opens a listening socket on the port and IP address of the local system
//
func NewTCPListen(c *ctx.Ctx, addr string) (l *ListenTCPStruct, e error) {

	c.Logf(ctx.LogLevelTrace, " called")

	c = c.NewWithCancel()

	listenconfig := net.ListenConfig{}
	listener, e := listenconfig.Listen(c.Context(), "tcp", addr)

	if e != nil {
		return l, e
	}

	accept := make(chan interface{}, TCPAcceptChannelLen)

	l = &ListenTCPStruct{
		ctx:      c,
		listener: listener,
		accept:   accept,
	}

	//
	// Go routine to run accept() on the socket and pass the connection back
	// via the channel.  The accept function can be canceled via the context
	//

	return l, e
}

//
//
//
func (l *ListenTCPStruct) RunTCP() {
	//	contextlib.Logf(l.ctx, contextlib.LevelTrace, lumerinlib.FileLineFunc()+" called")
	go l.goListenAccept()
}

//
// goAccept() go routine to accept connections and return new socket structs to the Accept() function
//
func (l *ListenTCPStruct) goListenAcceptTCP() {

	l.ctx.Logf(ctx.LogLevelTrace, " called")

	defer close(l.accept)

	for !l.Done() {
		// Network Accept()
		conn, e := l.listener.Accept()

		if e != nil {
			l.ctx.Logf(ctx.LogLevelError, " listener.Accept() returned error:%s", e)
			break
		}

		if conn == nil {
			l.ctx.Logf(ctx.LogLevelError, " Accept() returned empty connection")
			break
		}

		newsoc := createNewSocket(l.ctx, conn)

		l.accept <- newsoc
	}

	l.ctx.Logf(ctx.LogLevelTrace, " closing l.accept and exiting")
}

//
//
//
func (l *ListenTCPStruct) GetAcceptChan() <-chan interface{} {

	contextlib.Logf(l.ctx, contextlib.LevelTrace, lumerinlib.FileLineFunc()+" called")

	return l.accept

}

// Opens a listening socket on the port and IP address of the local system
//
// func NewListen(ctx context.Context, network string, addr string) (l interface{}, e error) {
func NewListen(ctx context.Context, network string, addr string) (l *ListenTCPStruct, e error) {

	contextlib.Logf(ctx, contextlib.LevelTrace, lumerinlib.FileLineFunc()+" called")

	//
	// network can be tcp, tcp4 or tcp6
	//
	switch network {
	case "tcp":
	case "tcp4":
	case "tcp6":
	default:
		return l, ErrSocTCPBadNetwork
	}

	ctx, cancel := context.WithCancel(ctx)
	c := cancel

	listenconfig := net.ListenConfig{}
	listener, e := listenconfig.Listen(ctx, network, addr)

	if e != nil {
		return l, e
	}

	// accept := make(chan *SocketTCPStruct, TCPAcceptChannelLen)
	accept := make(chan interface{}, TCPAcceptChannelLen)

	l = &ListenTCPStruct{
		listener: listener,
		ctx:      ctx,
		cancel:   c,
		accept:   accept,
		status: ListenerStatusStruct{
			connectionCount: 0,
		},
	}

	//
	// Go routine to run accept() on the socket and pass the connection back
	// via the channel.  The accept function can be canceled via the context
	//

	return l, e
}

//
//
//
func (l *ListenTCPStruct) Run() {
	//	contextlib.Logf(l.ctx, contextlib.LevelTrace, lumerinlib.FileLineFunc()+" called")
	go l.goListenAccept()
}

//
//
//
func (l *ListenTCPStruct) Ctx() context.Context {
	return l.ctx
}

//
// goAccept() go routine to accept connections and return new socket structs to the Accept() function
//
func (l *ListenTCPStruct) goListenAccept() {

	contextlib.Logf(l.ctx, contextlib.LevelTrace, lumerinlib.FileLineFunc()+" called")

	defer close(l.accept)

	for !l.Done() {
		// Network Accept()
		conn, e := l.listener.Accept()

		if e != nil {
			contextlib.Logf(l.ctx, contextlib.LevelError, lumerinlib.FileLineFunc()+" listener.Accept() returned error:%s", e)
			break
		}

		if conn == nil {
			contextlib.Logf(l.ctx, contextlib.LevelError, lumerinlib.FileLineFunc()+" Accept() returned empty connection")
			break
		}

		newsoc := createNewSocket(l.ctx, conn)

		l.accept <- newsoc
	}

	contextlib.Logf(l.ctx, contextlib.LevelTrace, lumerinlib.FileLineFunc()+" closing l.accept and exiting")
}

//
//
//
func (l *ListenTCPStruct) GetAcceptChan() <-chan interface{} {

	ctx.Logf(l.ctx, contextlib.LevelTrace, lumerinlib.FileLineFunc()+" called")

	return l.accept

}
