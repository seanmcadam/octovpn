package tcpsrv

import (
	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/netlib"
)

// Reset()
func (t *TcpServerStruct) Reset() error {
	log.Debugf("TCPSrv Reset()")

	if t.tcpconn != nil {
		t.tcpconn.Close()
		return nil
	} else {
		return netlib.ErrNetChannelDown
	}
}
