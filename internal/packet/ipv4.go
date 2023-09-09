package packet

type IPv4Packet struct{
	pSize PacketSizeType
}

func NewIPv4()(ap *IPv4Packet){
	ap = &IPv4Packet{}
	return ap
}

func MakeIPv4(raw []byte)(p *IPv4Packet, err error){
	return p, err
}

func (a *IPv4Packet) Size() PacketSizeType {
	return a.pSize
}

func (p *IPv4Packet)ToByte()(raw []byte){
	return raw
}