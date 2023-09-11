package tcpsrv

import (
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Send()
func (t *TcpServerStruct) Send(co *packet.PacketStruct) (err error) {
	if t == nil {
		return errors.ErrNetNilPointerMethod(log.Errf(""))
	}

	log.Debugf("TCPSrc Send:%v", co)

	if uint16(co.Size()) > uint16(t.config.Mtu) {
		return errors.ErrNetPacketTooBig(log.Errf(""))
	}

	if t.tcpconn != nil {
		return t.tcpconn.Send(co)
	}

	return errors.ErrNetChannelDown(log.Errf(""))

}
