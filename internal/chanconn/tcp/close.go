package tcp

import "github.com/seanmcadam/octovpn/octolib/log"

// Close()
func (t *TcpStruct) Close() {
	log.Debugf("TCP Close() called")
	close(t.Closech)
}

func (u *TcpStruct) closed() bool {
	select {
	case <-u.Closech:
		return true
	default:
		return false
	}
}
