package tcpsrv

import (
	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/netlib"
)

// Recv()
func (t *TcpServerStruct) Recv() (buf []byte, err error) {

	if t.tcpconn != nil {
		buf, err = t.tcpconn.Recv()

		if len(buf) > int(t.config.GetMtu()) {
			log.Warnf("TCPCli recv large packet %d > %d", len(buf), t.config.GetMtu())
		}
	} else {
		err = netlib.ErrNetChannelDown
	}

	return buf, err
}
