package tcp

import (
	"github.com/seanmcadam/octovpn/octotypes"
)

type TcpServerStruct struct {
	host string
	port octotypes.NetworkPort
}

// NewTcpServer()
// Returns a TcpServerStruct and error value
func NewTcpServer(host string, port octotypes.NetworkPort) (t *TcpServerStruct, err error) {
	t = &TcpServerStruct{
		host: host,
		port: port,
	}

	return t, err
}

// Close()
func (t *TcpServerStruct) Close() {
}

// Send()
func (t *TcpServerStruct) Send(buf []byte) (err error) {
	return err
}

// Recv()
func (t *TcpServerStruct) Recv() (buf []byte, err error) {
	return nil, err
}

// Reset()
func (t *TcpServerStruct) Reset() (err error) {
	return nil
}
