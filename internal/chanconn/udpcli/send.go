package udpcli

import (
	"github.com/seanmcadam/octovpn/internal/chanconn"
	"github.com/seanmcadam/octovpn/octolib/netlib"
)

// Send()
func (u *UdpClientStruct) Send(buf []byte) (err error) {

	if len(buf) + chanconn.PacketOverhead > int(u.config.GetMtu()) {
		return netlib.ErrNetPacketTooBig
	}

	if u.udpconn != nil {
		packet, err := chanconn.NewPacket(chanconn.PACKET_TYPE_UDP, buf)
		if err != nil {
			return err
		}
		return u.udpconn.Send(packet)
	}

	return netlib.ErrNetChannelDown

}
