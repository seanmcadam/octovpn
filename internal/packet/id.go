package packet

import (
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/octolib/errors"
)

type IDPacket struct {
	pSize PacketSizeType
	id []string
}

func NewID() (ap *IDPacket, err error) {
	ap = &IDPacket{}
	return ap, err
}

func MakeID(raw []byte) (p *IDPacket, err error) {
	if raw == nil {
		log.ErrorStack("Nil Parameter")
		return nil, errors.ErrPacketBadParameter(log.Errf("Nil raw data"))
	}

	return p, err
}

func (p *IDPacket) Size() PacketSizeType {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return 0
	}

	return p.pSize
}

func (p *IDPacket) ToByte() (raw []byte) {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return raw
	}

	return raw
}
