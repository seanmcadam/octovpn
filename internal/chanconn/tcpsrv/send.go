package tcpsrv

import (
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

// Send()
func (t *TcpServerStruct) Send(buf []byte) (err error) {

	if len(buf)+packetconn.PacketOverhead > int(t.config.GetMtu()) {
		return errors.ErrNetPacketTooBig
	}

	if t.tcpconn != nil {
		packet, err := packetconn.NewPacket(packetconn.PACKET_TYPE_TCP, buf)
		if err != nil {
			return err
		}
		return t.tcpconn.Send(packet)
	}

	return errors.ErrNetChannelDown

}
