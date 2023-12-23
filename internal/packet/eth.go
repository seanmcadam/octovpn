package packet

import (
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/octolib/errors"
)

type EthPacket struct {
	pSize PacketSizeType
}

func NewEth() (ap *EthPacket, err error) {

	ap = &EthPacket{}
	return ap, nil
}

func MakeEth(raw []byte) (p *EthPacket, err error) {
	if raw == nil {
		log.ErrorStack("Nil Parameter")
		return nil, errors.ErrPacketBadParameter(log.Errf("Nil raw data"))
	}

	return p, err
}

func (e *EthPacket) Size() PacketSizeType {
	if e == nil {
		log.ErrorStack("Nil Method Pointer")
		return 0
	}

	return e.pSize
}

func (e *EthPacket) ToByte() (raw []byte) {
	if e == nil {
		log.ErrorStack("Nil Method Pointer")
		return nil
	}

	return raw
}
