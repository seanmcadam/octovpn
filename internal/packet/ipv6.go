package packet

import (
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/octolib/errors"
)

type IPv6Packet struct {
	pSize PacketSizeType
}

func NewIPv6() (ap *IPv6Packet, err error) {
	ap = &IPv6Packet{}
	return ap, err
}

func MakeIPv6(raw []byte) (p *IPv6Packet, err error) {
	if raw == nil {
		log.ErrorStack("Nil Parameter")
		return nil, errors.ErrPacketBadParameter(log.Errf("Nil raw data"))
	}

	return p, err
}

func (p *IPv6Packet) Size() PacketSizeType {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return 0
	}

	return p.pSize
}

func (p *IPv6Packet) ToByte() (raw []byte) {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return raw
	}

	return raw
}
