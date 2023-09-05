package packet

type RouterStruct struct{
	pSize PacketSizeType
}

func NewRouter()(ap *RouterStruct){
	ap = &RouterStruct{}
	return ap
}

func MakeRouter(raw []byte)(p *RouterStruct, err error){
	return p, err
}

func (a *RouterStruct) Size() PacketSizeType {
	return a.pSize
}

func (p *RouterStruct)ToByte()(raw []byte){
	return raw
}