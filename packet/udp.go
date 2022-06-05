package packet

import "fmt"

type UDPFrame []byte
type UDPPayload []byte

func (f UDPFrame) SourcePort() (port uint16) {
	port = uint16(f[0])
	port = port << 8
	port = port + uint16(f[1])
	return port
}

func (f UDPFrame) DestPort() (port uint16) {
	port = uint16(f[2])
	port = port << 8
	port = port + uint16(f[3])
	return port
}

func (f UDPFrame) Length() (len uint16) {
	len = uint16(f[4])
	len = len << 8
	len = len + uint16(f[5])
	return len
}

func (f UDPFrame) CheckSum() (cs uint16) {
	cs = uint16(f[4])
	cs = cs << 8
	cs = cs + uint16(f[7])
	return cs
}

func (f UDPFrame) Payload() (payload UDPPayload) {
	payload = UDPPayload(f[8:])

	if int(f.Length()) != len(payload) {
		panic(fmt.Sprintf("UDP Payload size does not match %d != %d", f.Length(), len(payload)))
	}
	return payload
}
