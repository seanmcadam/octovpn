package udpsrv

import (
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/octolib/errors"
)

// Reset()
// Close current connection, causing a reset
func (u *UdpServerStruct) Reset() error {
	if u == nil {
		return errors.ErrNetNilMethodPointer(log.Errf(""))
	}

	log.Debugf("UDPSrv Reset()")

	if u.udpconn != nil {
		u.udpconn.Cancel()
		return nil
	} else {
		return errors.ErrNetChannelDown(log.Errf(""))
	}
}
