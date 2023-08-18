package tcp

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/seanmcadam/octovpn/internal/channel/netconn"
	"github.com/seanmcadam/octovpn/octolib"
	"github.com/seanmcadam/octovpn/octotypes"
)

type TcpClientStruct struct {
	linkUp  bool
	tcptype string
	host    string
	port    octotypes.NetworkPort
	address string
	conn    *netconn.NetConnStruct
	sendch  chan []byte
	recvch  chan []byte
	closech chan interface{}
}

func NewTcpClient(host string, port octotypes.NetworkPort) (t *TcpClientStruct, err error) {
	return new("tcp", host, port)
}
func NewTcp4Client(host string, port octotypes.NetworkPort) (t *TcpClientStruct, err error) {
	return new("tcp4", host, port)
}
func NewTcp6Client(host string, port octotypes.NetworkPort) (t *TcpClientStruct, err error) {
	return new("tcp6", host, port)
}

func new(tcptype string, host string, port octotypes.NetworkPort) (t *TcpClientStruct, err error) {

	// Validate the host string Valid Hostname or IP
	if !(octolib.ValidIP(host) || octolib.ValidHost(host)) {
		return nil, fmt.Errorf("Bad host value:%s", host)
	}

	t = &TcpClientStruct{
		linkUp:  false,
		tcptype: tcptype,
		host:    host,
		port:    port,
		address: fmt.Sprintf("%s:%d", host, port),
		conn:    nil,
		sendch:  make(chan []byte),
		recvch:  make(chan []byte),
		closech: make(chan interface{}),
	}

	return t, err
}

// Close()
func (t *TcpClientStruct) Close() {

}

// Start()
func (t *TcpClientStruct) Start() {
	go t.goRun()
}

// Send()
func (t *TcpClientStruct) Send(buf []byte) (err error) {

//	if !t.linkUp {
//		return fmt.Errorf("Link to %s:%d down", t.host, t.port)
//	}
//
//	go func(buf []byte) {
//		t.sendch <- buf
//	}(buf)
	return err

}

// Recv()
func (t *TcpClientStruct) Recv() (buf []byte, err error) {

//	buf = <-t.recvch

	return buf, err
}

// Reset()
// Close current connection, causing a reset
func (t *TcpClientStruct) Reset() (err error) {
	return err
}

// goRun()
func (t *TcpClientStruct) goRun() {

	for {
		// Establish connection
		conn, err := t.connect()
		if err != nil {
			log.Printf("connection failed %s: %s", t.address, err)
			t.conn = nil
			time.Sleep(1 * time.Second)
			continue
		}

		t.conn = netconn.NewNetConn(conn)

		for {
			if t.conn == nil {
				break
			}

			select {
			case <-t.closech:
				t.linkUp = false
				return

			case send := <-t.sendch:

				if !t.linkUp {
					break
				}
				if t.conn != nil {
					l, err := t.conn.Write(send)
					if err != nil {
						log.Printf("Write Error: %s", err)
						break

					}

					if len(send) != l {
						log.Printf("Write failed  lenths dont match: %d, sent:%d", len(send), l)
						break
					}
				}
			}

			t.conn.Close()
			t.conn = nil
		}
	}

}

func (t *TcpClientStruct) connect() (conn net.Conn, err error) {
	conn, err = net.Dial(t.tcptype, t.address)
	return conn, err
}
