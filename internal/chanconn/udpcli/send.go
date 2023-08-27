package udpcli

import (
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

// Send()
func (u *UdpClientStruct) Send(co *packetconn.ConnPacket) (err error) {

	if co.GetSize() > int(u.config.GetMtu()) {
		return errors.ErrNetPacketTooBig
	}

	if u.udpconn != nil {
		return u.udpconn.Send(co)
	}

	return errors.ErrNetChannelDown

}
