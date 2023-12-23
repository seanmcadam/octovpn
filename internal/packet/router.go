package packet

import (
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/octolib/errors"
)

type RouterPacket struct {
	pSize PacketSizeType
}

func NewRouter() (ap *RouterPacket) {
	ap = &RouterPacket{}
	return ap
}

func MakeRouter(raw []byte) (p *RouterPacket, err error) {
	if raw == nil {
		log.ErrorStack("Nil Parameter")
		return nil, errors.ErrPacketBadParameter(log.Errf("Nil raw data"))
	}

	return p, err
}

func (p *RouterPacket) Size() PacketSizeType {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return 0
	}

	return p.pSize
}

func (p *RouterPacket) ToByte() (raw []byte) {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return raw
	}
	return raw
}
