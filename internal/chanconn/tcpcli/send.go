package tcpcli

import (
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/packet/packetchan"
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

// Send()
func (t *TcpClientStruct) Send(cp *packetchan.ChanPacket) (err error) {

	if cp.GetSize() > int(t.config.GetMtu()) {
		return errors.ErrNetPacketTooBig
	}

	if t.tcpconn != nil {
		packet, err := packetconn.NewPacket(packetconn.PACKET_TYPE_TCP, cp)
		if err != nil {
			return err
		}
		return t.tcpconn.Send(packet)
	}

	return errors.ErrNetChannelDown

}
