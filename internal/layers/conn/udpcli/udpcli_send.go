package udpcli

import (
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Send()
func (u *UdpClientStruct) Send(co *packet.PacketStruct) (err error) {
	if u == nil {
		return errors.ErrNetNilPointerMethod(log.Errf(""))
	}


	if uint16(co.Size()) > u.config.GetMtu() {
		return errors.ErrNetPacketTooBig(log.Errf(""))
	}

	if u.udpconn != nil {
		return u.udpconn.Send(co)
	}

	return errors.ErrNetChannelDown(log.Errf(""))

}
