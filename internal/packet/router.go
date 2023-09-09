package packet

type RouterPacket struct{
	pSize PacketSizeType
}

func NewRouter()(ap *RouterPacket){
	ap = &RouterPacket{}
	return ap
}

func MakeRouter(raw []byte)(p *RouterPacket, err error){
	return p, err
}

func (a *RouterPacket) Size() PacketSizeType {
	return a.pSize
}

func (p *RouterPacket)ToByte()(raw []byte){
	return raw
}