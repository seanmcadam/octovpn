package tcpsrv

import (
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Reset()
func (t *TcpServerStruct) Reset() error {
	log.Debugf("TCPSrv Reset()")

	if t.tcpconn != nil {
		t.tcpconn.Close()
		return nil
	} else {
		return errors.ErrNetChannelDown
	}
}
