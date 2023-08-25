package tcpcli

import "github.com/seanmcadam/octovpn/octolib/log"

// Close()
func (t *TcpClientStruct) Close() {
	log.Debugf("TCPCli Close()")
	t.cx.Done()
	// close(t.closech)
}
