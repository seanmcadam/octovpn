package udpcli

import (
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

// Send()
func (u *UdpClientStruct) Send(buf []byte) (err error) {

	if len(buf)+packetconn.PacketOverhead > int(u.config.GetMtu()) {
		return errors.ErrNetPacketTooBig
	}

	if u.udpconn != nil {
		packet, err := packetconn.NewPacket(packetconn.PACKET_TYPE_UDP, buf)
		if err != nil {
			return err
		}
		return u.udpconn.Send(packet)
	}

	return errors.ErrNetChannelDown

}
