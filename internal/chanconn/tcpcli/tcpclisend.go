package tcpcli

import (
	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/octolib/errors"
)

// Send()
func (t *TcpClientStruct) Send(co interfaces.PacketInterface) (err error) {

	if uint16(co.Size()) > t.config.GetMtu() {
		return errors.ErrNetPacketTooBig
	}

	if t.tcpconn != nil {
		return t.tcpconn.Send(co)
	}

	return errors.ErrNetChannelDown

}
