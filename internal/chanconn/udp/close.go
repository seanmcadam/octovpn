package udp

import "github.com/seanmcadam/octovpn/octolib/log"

// Close()
func (u *UdpStruct) Close() {
	log.Debugf("UDP Close() called")
	close(u.Closech)
}

func (u *UdpStruct) closed() bool {
	select {
	case <-u.Closech:
		return true
	default:
		return false
	}
}
