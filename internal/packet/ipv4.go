package packet

import (
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/octolib/errors"
)

type IPv4Packet struct {
	pSize PacketSizeType
}

func NewIPv4() (ap *IPv4Packet, err error) {
	ap = &IPv4Packet{}
	return ap, err
}

func MakeIPv4(raw []byte) (p *IPv4Packet, err error) {
	if raw == nil {
		log.ErrorStack("Nil Parameter")
		return nil, errors.ErrPacketBadParameter(log.Errf("Nil raw data"))
	}

	return p, err
}

func (p *IPv4Packet) Size() PacketSizeType {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return 0
	}

	return p.pSize
}

func (p *IPv4Packet) ToByte() (raw []byte) {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return raw
	}

	return raw
}
