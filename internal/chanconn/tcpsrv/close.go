package tcpsrv

import "github.com/seanmcadam/octovpn/octolib/log"

// Close()
func (t *TcpServerStruct) Close() {
	log.Debugf("TCPSrv Close()")
	close(t.closech)
}
