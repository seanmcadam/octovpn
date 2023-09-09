package packet

type EthPacket struct{
	pSize PacketSizeType
}
func NewEth()(ap *EthPacket){
	ap = &EthPacket{}
	return ap
}

func MakeEth(raw []byte)(p *EthPacket, err error){
	return p, err
}

func (a *EthPacket) Size() PacketSizeType {
	return a.pSize
}

func (e *EthPacket)ToByte()(raw []byte){
	return raw
}