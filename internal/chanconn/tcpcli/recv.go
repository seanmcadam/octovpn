package tcpcli

import (
	"github.com/seanmcadam/octovpn/internal/chanconn"
	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/netlib"
)

// Recv()
func (t *TcpClientStruct) Recv() (buf []byte, err error) {
	var packet *chanconn.ConnPacket

	if t.tcpconn != nil {
		packet, err = t.tcpconn.Recv()
		if err != nil {
			return nil, err
		}

		if int(packet.GetLength())+chanconn.PacketOverhead > int(t.config.GetMtu()) {
			log.Warnf("TCPCli recv large packet %d > %d", packet.GetLength(), t.config.GetMtu())
		}
	} else {
		err = netlib.ErrNetChannelDown
	}

	return packet.GetPayload(), err
}
