package tcpsrv

import (
	"github.com/seanmcadam/octovpn/octolib/netlib"
)

// Send()
func (t *TcpServerStruct) Send(buf []byte) (err error) {

	if len(buf) > int(t.config.GetMtu()) {
		return netlib.ErrNetPacketTooBig
	}

	if t.tcpconn != nil {
		return t.tcpconn.Send(buf)
	}

	return netlib.ErrNetChannelDown

}
