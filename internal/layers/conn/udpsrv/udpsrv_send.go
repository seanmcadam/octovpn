package udpsrv

import (
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/errors"
)

// Send()
func (u *UdpServerStruct) Send(co *packet.PacketStruct) (err error) {
	if u == nil {
		return errors.ErrNetNilMethodPointer(log.Errf(""))
	}

	if co == nil {
		log.ErrorStack("Nil Packet in Send()")
		return log.Errf("Nil Packet in Send()")
	}

	if uint16(co.Size()) > uint16(u.config.Mtu) {
		return errors.ErrNetPacketTooBig(log.Errf(""))
	}

	if u.udpconn != nil {
		return u.udpconn.Send(co)
	}

	return errors.ErrNetChannelDown(log.Errf(""))

}
