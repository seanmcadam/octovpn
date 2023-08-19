package netconn

import "fmt"

const MaxNetworkPacket = 1520

type NetPacketStruct struct {
	length uint16 
	count uint32 
	packet []byte
}

func NewNetworkPacket(c uint32, buf []byte)(p *NetPacketStruct, err error){
	l := len(buf)
	if l == 0 {
		return nil, fmt.Errorf("BufferLength is 0")
	}

	p = &NetPacketStruct{
		length: uint16(len(buf)+4),
		count:c,
		packet:buf,
	}
return p, err
}

func (np *NetPacketStruct)Packet()(p []byte){
	return np.packet
}
func (np *NetPacketStruct)Length()(l uint16){
	return(np.length - uint16(4))
}

