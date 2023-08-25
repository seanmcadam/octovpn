package tcpcli

import (
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

// Recv()
func (t *TcpClientStruct) Recv() (buf []byte, err error) {
	var packet *packetconn.ConnPacket

	if t.tcpconn == nil {
		err = errors.ErrNetChannelDown
		return nil, err
	}

	packet, err = t.tcpconn.Recv()
	if err != nil {
		return nil, err
	}

	if int(packet.GetLength())+packetconn.PacketOverhead > int(t.config.GetMtu()) {
		log.Warnf("TCPCli recv large packet %d > %d", packet.GetLength(), t.config.GetMtu())
	}

	return packet.GetPayload(), err
}
