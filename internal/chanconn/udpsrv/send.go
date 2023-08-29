package udpsrv

import (
	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/octolib/errors"
)

// Send()
func (u *UdpServerStruct) Send(co interfaces.PacketInterface) (err error) {

	if uint16(co.Size()) > u.config.GetMtu() {
		return errors.ErrNetPacketTooBig
	}

	if u.udpconn != nil {
		return u.udpconn.Send(co)
	}

	return errors.ErrNetChannelDown

}
