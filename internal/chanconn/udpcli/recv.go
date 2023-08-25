package udpcli

import (
	"github.com/seanmcadam/octovpn/internal/chanconn"
	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/netlib"
)

// Recv()
func (u *UdpClientStruct) Recv() (buf []byte, err error) {
	var packet *chanconn.ConnPacket

	if u.udpconn != nil {
		packet, err = u.udpconn.Recv()
		if err != nil {
			return nil, err
		}

		if (int(packet.GetLength()) + chanconn.PacketOverhead) > int(u.config.GetMtu()) {
			log.Warnf("TCPCli recv large packet %d > %d", len(buf), u.config.GetMtu())
		}
	} else {
		err = netlib.ErrNetChannelDown
	}

	return packet.GetPayload(), err
}
