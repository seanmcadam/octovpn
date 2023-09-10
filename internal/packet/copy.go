package packet

import (
	"fmt"

	"github.com/seanmcadam/octovpn/internal/pinger"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (p *PacketStruct) Copy() (copy *PacketStruct, err error) {

	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return nil, errors.ErrPacketBadParameter(log.Errf("Nil Methos Pointer"))
	}

	log.Debug("Do we really need to copy all of this data?")

	var fields []interface{}

	if p.counter != nil {
		fields = append(fields, p.counter)
	}

	if p.ping != nil {
		fields = append(fields, p.ping)
	} else if p.pong != nil {
		fields = append(fields, p.pong)
	}

	if p.packet != nil {
		fields = append(fields, p.packet)
	} else if p.ipv4 != nil {
		fields = append(fields, p.ipv4)
	} else if p.ipv6 != nil {
		fields = append(fields, p.ipv6)
	} else if p.eth != nil {
		fields = append(fields, p.eth)
	} else if p.router != nil {
		fields = append(fields, p.router)
	} else if p.id != nil {
		fields = append(fields, p.id)
	} else if p.auth != nil {
		fields = append(fields, p.auth)
	} else if p.raw != nil {
		fields = append(fields, p.raw)
	}

	copy, err = NewPacket(p.Sig(), fields...)
	if err != nil {
		return nil, err
	}

	return copy, nil
}

func (p *PacketStruct) CopyAck() (copy *PacketStruct, err error) {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return nil, errors.ErrPacketBadParameter(log.Errf("Nil Methos Pointer"))
	}

	var sig PacketSigType

	if p.pSig.Width32() {
		sig = SIG_CHAN_32_ACK
	} else if p.pSig.Width64() {
		sig = SIG_CHAN_64_ACK
	} else {
		log.FatalStack("No Size")
	}

	copy, err = NewPacket(sig, p.Counter())
	if err != nil {
		return nil, err
	}

	return copy, nil
}

func (p *PacketStruct) CopyPong() (ppong *PacketStruct, err error) {
	if p == nil {
		log.ErrorStack("Nil Method Pointer")
		return nil, errors.ErrPacketBadParameter(log.Errf("Nil Methos Pointer"))
	}

	if !p.pSig.Ping() {
		err = fmt.Errorf("CopyPong() Not a ping packet")
	}

	ppong = &PacketStruct{
		pSig: (p.pSig & (^FIELD_PING)) | FIELD_PONG,
		pong: pinger.Pong(p.ping.Copy()),
	}

	if p.pSig.Counter() {
		ppong.counter = p.counter.Copy()
	}

	return ppong, err
}
