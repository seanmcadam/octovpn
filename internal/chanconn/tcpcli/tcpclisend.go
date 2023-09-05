package tcpcli

import (
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/errors"
)

// Send()
func (t *TcpClientStruct) Send(co *packet.PacketStruct) (err error) {

	if uint16(co.Size()) > t.config.GetMtu() {
		return errors.ErrNetPacketTooBig
	}

	if t.tcpconn != nil {
		return t.tcpconn.Send(co)
	}

	return errors.ErrNetChannelDown

}
