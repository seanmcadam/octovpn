package tcp

import "github.com/seanmcadam/octovpn/octolib/log"

// Close()
func (t *TcpStruct) Close() {
	log.Debugf("TCP Close() called")
	close(t.Closech)
}
