package packet

type IPv6Packet struct{
	pSize PacketSizeType
}

func NewIPv6()(ap *IPv6Packet){
	ap = &IPv6Packet{}
	return ap
}

func MakeIPv6(raw []byte)(p *IPv6Packet, err error){
	return p, err
}

func (a *IPv6Packet) Size() PacketSizeType {
	return a.pSize
}

func (p *IPv6Packet)ToByte()(raw []byte){
	return raw
}