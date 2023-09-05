package packet

type IPv4Struct struct{
	pSize PacketSizeType
}

func NewIPv4()(ap *IPv4Struct){
	ap = &IPv4Struct{}
	return ap
}

func MakeIPv4(raw []byte)(p *IPv4Struct, err error){
	return p, err
}

func (a *IPv4Struct) Size() PacketSizeType {
	return a.pSize
}

func (p *IPv4Struct)ToByte()(raw []byte){
	return raw
}