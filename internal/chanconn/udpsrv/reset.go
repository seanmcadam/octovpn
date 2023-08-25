package udpsrv

import (
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Reset()
// Close current connection, causing a reset
func (u *UdpServerStruct) Reset() error {

	log.Debugf("UDPSrv Reset()")

	if u.udpconn != nil {
		u.udpconn.Close()
		return nil
	} else {
		return errors.ErrNetChannelDown
	}
}
