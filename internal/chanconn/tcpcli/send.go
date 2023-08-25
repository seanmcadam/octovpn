package tcpcli

import (
	"github.com/seanmcadam/octovpn/internal/chanconn"
	"github.com/seanmcadam/octovpn/octolib/netlib"
)

// Send()
func (t *TcpClientStruct) Send(buf []byte) (err error) {

	if len(buf)+chanconn.PacketOverhead > int(t.config.GetMtu()) {
		return netlib.ErrNetPacketTooBig
	}

	if t.tcpconn != nil {
		packet, err := chanconn.NewPacket(chanconn.PACKET_TYPE_TCP, buf)
		if err != nil {
			return err
		}
		return t.tcpconn.Send(packet)
	}

	return netlib.ErrNetChannelDown

}
