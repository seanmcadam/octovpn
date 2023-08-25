package udpsrv

import "github.com/seanmcadam/octovpn/octolib/log"

// Close()
func (u *UdpServerStruct) Close() {
	log.Debugf("UDPSrv Close()")
	close(u.closech)
}