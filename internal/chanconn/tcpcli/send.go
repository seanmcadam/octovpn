package tcpcli

import (
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

// Send()
func (t *TcpClientStruct) Send(co *packetconn.ConnPacket) (err error) {

	if co.GetSize() > int(t.config.GetMtu()) {
		return errors.ErrNetPacketTooBig
	}

	if t.tcpconn != nil {
		return t.tcpconn.Send(co)
	}

	return errors.ErrNetChannelDown

}
