package tcpcli

import (
	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/netlib"
)

// Reset()
// Close current connection, causing a reset
func (t *TcpClientStruct) Reset() error {

	log.Debugf("TCPCli Reset()")

	if t.tcpconn != nil {
		t.tcpconn.Close()
		return nil
	} else {
		return netlib.ErrNetChannelDown
	}
}
