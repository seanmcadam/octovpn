package tcpcli

import (
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Reset()
// Close current connection, causing a reset
func (t *TcpClientStruct) Reset() error {
	if t == nil {
		return errors.ErrNetNilMethodPointer(log.Errf(""))
	}

	log.Debugf("TCPCli Reset()")

	if t.tcpconn != nil {
		t.tcpconn.Cancel()
		return nil
	} else {
		return errors.ErrNetChannelDown(log.Errf(""))
	}
}
