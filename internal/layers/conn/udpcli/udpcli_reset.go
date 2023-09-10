package udpcli

import (
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Reset()
// Close current connection, causing a reset
func (u *UdpClientStruct) Reset() error {
	if u == nil {
		return errors.ErrNetNilPointerMethod(log.Errf(""))
	}


	log.Debugf("UDPCli Reset()")

	if u.udpconn != nil {
		u.udpconn.Cancel()
		return nil
	} else {
		return errors.ErrNetChannelDown(log.Errf(""))
	}
}
