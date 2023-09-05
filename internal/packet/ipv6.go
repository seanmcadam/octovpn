package packet

type IPv6Struct struct{
	pSize PacketSizeType
}

func NewIPv6()(ap *IPv6Struct){
	ap = &IPv6Struct{}
	return ap
}

func MakeIPv6(raw []byte)(p *IPv6Struct, err error){
	return p, err
}

func (a *IPv6Struct) Size() PacketSizeType {
	return a.pSize
}

func (p *IPv6Struct)ToByte()(raw []byte){
	return raw
}