package tcpsrv

import (
	"fmt"
	"net"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/chanconn/tcp"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type TcpServerStruct struct {
	config      *settings.NetworkStruct
	address     string
	tcplistener *net.TCPListener
	tcpaddr     *net.TCPAddr
	tcpconn     *tcp.TcpStruct
	tcpconnch   chan *tcp.TcpStruct
	closech     chan interface{}
	resetch     chan interface{}
}

// NewTcpServer()
// Returns a TcpServerStruct and error value
func New(config *settings.NetworkStruct) (tcpserver interfaces.ChannelInterface, err error) {

	t := &TcpServerStruct{
		config:      config,
		address:     fmt.Sprintf("%s:%d", config.GetHost(), config.GetPort()),
		tcplistener: nil,
		tcpaddr:     nil,
		tcpconn:     nil,
		closech:     make(chan interface{}),
		resetch:     make(chan interface{}),
		tcpconnch:   make(chan *tcp.TcpStruct),
	}

	// Recheck this each time, the IP could change or rotate
	t.tcpaddr, err = net.ResolveTCPAddr(t.config.Proto, t.address)
	if err != nil {
		return nil, fmt.Errorf("ResolveTCPAddr Failed:%s", err)
	}

	t.tcplistener, err = net.ListenTCP(t.config.Proto, t.tcpaddr)
	if err != nil {
		return nil, fmt.Errorf("ListenTCP Failed:%s", err)
	}

	go t.goListen()
	go t.goRun()

	return t, err

}

// goRun()
// Loop on
// 	Establish Connection
// 	Start Send and Recv Goroutines
// 	Monitor reset request
//

func (t *TcpServerStruct) goRun() {

	defer func(t *TcpServerStruct) {
		if t.tcplistener != nil {
			t.tcplistener.Close()
			t.tcplistener = nil
		}
		if t.tcpconn != nil {
			t.tcpconn.Close()
			t.emptyconn()
			t.tcpconn = nil
		}
		close(t.resetch)
		close(t.tcpconnch)
	}(t)

	for {
		var tcpconnclosech chan interface{}

		select {
		case conn := <-t.tcpconnch:
			log.Info("New incoming TCP Connection")

			// Authenticate

			// Terminate other connection

			if t.tcpconn != nil {
				log.Debug("Shutdown Previous connection")
				// Existing Go Routines will close out and shutdown.
				t.tcpconn.Close()
				t.tcpconn = nil
			}

			// Start New Connection
			go t.goNewConn(conn)

		case <-tcpconnclosech:
			continue

		case <-t.closech:
			return

		case <-t.resetch:
			if t.tcpconn != nil {
				t.tcpconn.Close()
			}
		}
	}
}

func (t *TcpServerStruct) emptyconn() {
	for {
		select {
		case <-t.tcpconnch:
		default:
			return
		}
	}
}
