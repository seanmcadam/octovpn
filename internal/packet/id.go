package packet

type IDStruct struct {
	pSize PacketSizeType

	id []string
}

func NewID()(ap *IDStruct){
	ap = &IDStruct{}
	return ap
}

func MakeID(raw []byte)(p *IDStruct, err error){
	return p, err
}

func (a *IDStruct) Size() PacketSizeType {
	return a.pSize
}

func (i *IDStruct)ToByte()(raw []byte){
	return raw
}