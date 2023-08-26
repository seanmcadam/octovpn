package udpsrv

import (
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/packet/packetchan"
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

// Send()
func (u *UdpServerStruct) Send(cp *packetchan.ChanPacket) (err error) {

	if cp.GetSize() > int(u.config.GetMtu()) {
		return errors.ErrNetPacketTooBig
	}

	if u.udpconn != nil {
		packet, err := packetconn.NewPacket(packetconn.PACKET_TYPE_UDP, cp)
		if err != nil {
			return err
		}
		return u.udpconn.Send(packet)
	}

	return errors.ErrNetChannelDown

}
