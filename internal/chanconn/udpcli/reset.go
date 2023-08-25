package udpcli

import (
	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/netlib"
)

// Reset()
// Close current connection, causing a reset
func (u *UdpClientStruct) Reset() error {

	log.Debugf("UDPCli Reset()")

	if u.udpconn != nil {
		u.udpconn.Close()
		return nil
	} else {
		return netlib.ErrNetChannelDown
	}
}
