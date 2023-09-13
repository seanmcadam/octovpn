package tcpsrv

import (
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Reset()
func (t *TcpServerStruct) Reset() error {
	if t == nil {
		return errors.ErrNetNilMethodPointer(log.Errf(""))
	}

	log.Debugf("TCPSrv Reset()")

	if t.tcpconn != nil {
		t.tcpconn.Cancel()
		return nil
	} else {
		return errors.ErrNetChannelDown(log.Errf(""))
	}
}
