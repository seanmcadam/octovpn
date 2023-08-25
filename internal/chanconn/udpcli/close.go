package udpcli

import "github.com/seanmcadam/octovpn/octolib/log"

// Close()
func (u *UdpClientStruct) Close() {
	log.Debugf("UDPCli Close()")
	close(u.closech)
}