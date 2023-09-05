package packet

type EthStruct struct{
	pSize PacketSizeType
}
func NewEth()(ap *EthStruct){
	ap = &EthStruct{}
	return ap
}

func MakeEth(raw []byte)(p *EthStruct, err error){
	return p, err
}

func (a *EthStruct) Size() PacketSizeType {
	return a.pSize
}

func (e *EthStruct)ToByte()(raw []byte){
	return raw
}