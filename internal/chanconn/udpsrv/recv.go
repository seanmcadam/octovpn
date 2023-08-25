package udpsrv

import (
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

// Recv()
func (u *UdpServerStruct) Recv() (buf []byte, err error) {
	var packet *packetconn.ConnPacket

	if u.udpconn == nil {
		err = errors.ErrNetChannelDown
	}

	packet, err = u.udpconn.Recv()
	if err != nil {
		return nil, err
	}

	if (int(packet.GetLength()) + packetconn.PacketOverhead) > int(u.config.GetMtu()) {
		log.Warnf("TCPSrv recv large packet %d > %d", len(buf), u.config.GetMtu())
	}

	return packet.GetPayload(), err
}
