package packet

type AuthStruct struct {
	pSize PacketSizeType
}

func NewAuth()(ap *AuthStruct){
	ap = &AuthStruct{}
	return ap
}

func MakeAuth(raw []byte) (p *AuthStruct, err error) {
	return p, err
}

func (a *AuthStruct) Size() PacketSizeType {
	return a.pSize
}

func (a *AuthStruct) ToByte() (raw []byte) {
	return raw
}
