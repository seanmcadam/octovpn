package packet

import "fmt"

type TCPFrame []byte
type TCPPayload []byte

func (f TCPFrame) SourcePort() (port uint16) {
	port = uint16(f[0])
	port = port << 8
	port = port + uint16(f[1])
	return port
}

func (f TCPFrame) DestPort() (port uint16) {
	port = uint16(f[2])
	port = port << 8
	port = port + uint16(f[3])
	return port
}

func (f TCPFrame) SequenceNumber() (seq uint32) {
	seq = uint32(f[4])
	seq = seq << 8
	seq = seq + uint32(f[5])
	seq = seq << 8
	seq = seq + uint32(f[6])
	seq = seq << 8
	seq = seq + uint32(f[7])

	return seq
}

func (f TCPFrame) AckNumber() (ack uint32) {
	ack = uint32(f[8])
	ack = ack << 8
	ack = ack + uint32(f[9])
	ack = ack << 8
	ack = ack + uint32(f[10])
	ack = ack << 8
	ack = ack + uint32(f[11])

	return ack
}

func (f TCPFrame) DataOffset() (offset uint8) {
	offset = f[12]
	offset = offset & 0xf0
	offset = offset >> 4
	if offset < 5 || offset > 15 {
		panic(fmt.Sprintf("DataOffset out of range %d", offset))
	}
	return offset
}

func (f TCPFrame) flagNS() (b bool) {
	b = false
	if (f[12] & 0x01) > 0 {
		b = true
	}
	return b
}

func (f TCPFrame) flagCWR() (b bool) {
	b = false
	if (f[13] & 0x80) > 0 {
		b = true
	}
	return b
}

func (f TCPFrame) flagECE() (b bool) {
	b = false
	if (f[13] & 0x40) > 0 {
		b = true
	}
	return b
}

func (f TCPFrame) flagURG() (b bool) {
	b = false
	if (f[13] & 0x20) > 0 {
		b = true
	}
	return b
}

func (f TCPFrame) flagACK() (b bool) {
	b = false
	if (f[13] & 0x10) > 0 {
		b = true
	}
	return b
}

func (f TCPFrame) flagPSH() (b bool) {
	b = false
	if (f[13] & 0x08) > 0 {
		b = true
	}
	return b
}

func (f TCPFrame) flagRST() (b bool) {
	b = false
	if (f[13] & 0x04) > 0 {
		b = true
	}
	return b
}

func (f TCPFrame) flagSYN() (b bool) {
	b = false
	if (f[13] & 0x02) > 0 {
		b = true
	}
	return b
}

func (f TCPFrame) flagFIN() (b bool) {
	b = false
	if (f[13] & 0x01) > 0 {
		b = true
	}
	return b
}

func (f TCPFrame) WindowSize() (ws uint16) {
	ws = uint16(f[14])
	ws = ws << 8
	ws = ws + uint16(f[15])
	return ws
}

func (f TCPFrame) CheckSum() (cs uint16) {
	cs = uint16(f[16])
	cs = cs << 8
	cs = cs + uint16(f[17])
	return cs
}

func (f TCPFrame) UrgentPointer() (up uint16) {
	up = uint16(f[18])
	up = up << 8
	up = up + uint16(f[19])
	return up
}

func (f TCPFrame) Payload() (payload TCPPayload) {
	do := f.DataOffset()
	b := do * 4
	payload = TCPPayload(f[b:])

	return payload
}
