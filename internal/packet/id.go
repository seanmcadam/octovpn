package packet

type IDPacket struct {
	pSize PacketSizeType

	id []string
}

func NewID()(ap *IDPacket){
	ap = &IDPacket{}
	return ap
}

func MakeID(raw []byte)(p *IDPacket, err error){
	return p, err
}

func (a *IDPacket) Size() PacketSizeType {
	return a.pSize
}

func (i *IDPacket)ToByte()(raw []byte){
	return raw
}