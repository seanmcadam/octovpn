package tcpcli

import (
	"fmt"
	"net"
	"time"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/chanconn/tcp"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type TcpClientStruct struct {
	//
	// Link is currently active
	//
	config  *settings.NetworkStruct
	address string
	tcpaddr *net.TCPAddr
	tcpconn *tcp.TcpStruct
	closech chan interface{}
	resetch chan interface{}
}

func New(config *settings.NetworkStruct) (tcpclient interfaces.ChannelInterface, err error) {

	t := &TcpClientStruct{
		config:  config,
		address: fmt.Sprintf("%s:%d", config.GetHost(), config.GetPort()),
		tcpaddr: nil,
		tcpconn: nil,
		closech: make(chan interface{}),
		resetch: make(chan interface{}),
	}

	// Do an initial check and fail if it fails
	// Recheck this each time, the IP could change or rotate
	t.tcpaddr, err = net.ResolveTCPAddr(t.config.Proto, t.address)
	if err != nil {
		return nil, err
	}

	go t.goRun()
	return t, err
}

// goRun()
// Loop on
//
//	Establish Connection
//	Start Send and Recv Goroutines
//	Monitor reset request
func (t *TcpClientStruct) goRun() {

	defer func(t *TcpClientStruct) {
		if t.tcpconn != nil {
			t.tcpconn.Close()
			t.tcpconn = nil
		}
		close(t.closech)
		close(t.resetch)

	}(t)

TCPFOR:
	for {
		var err error
		var conn *net.TCPConn

		// Dial it and keep trying forever
		conn, err = net.DialTCP(t.config.Proto, nil, t.tcpaddr)

		if err != nil {
			log.Warnf("connection failed %s: %s, wait", t.address, err)
			t.tcpconn = nil
			time.Sleep(1 * time.Second)

			if conn != nil {
				select {
				case <-t.closech:
					log.Debug("tcpcli goRun() closed")
					return
				default:
				}
			}
			continue TCPFOR
		}

		if conn == nil {
			log.Fatal("conn == nil")
		}

		log.Info("New TCP Connection")

		t.tcpconn = tcp.NewTCP(conn)
		if t.tcpconn == nil {
			log.Fatal("tcpconn == nil")
		}

		closech := t.tcpconn.Closech

	TCPCLOSE:
		for {
			select {
			case <-t.closech:
				log.Debug("TCPCli Closing Down")
				return
			case <-closech:
				log.Debug("TCPCli Channel Closed")
				t.tcpconn = nil
				break TCPCLOSE
			}
		}
	}
}
