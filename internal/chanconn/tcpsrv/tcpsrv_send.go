package tcpsrv

import (
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Send()
func (t *TcpServerStruct) Send(co *packet.PacketStruct) (err error) {
	log.Debugf("TCPSrc Send:%v", co)

	if uint16(co.Size()) > t.config.GetMtu() {
		return errors.ErrNetPacketTooBig
	}

	if t.tcpconn != nil {
		return t.tcpconn.Send(co)
	}

	return errors.ErrNetChannelDown

}
