package packetsite

const ConnOverhead int = 4
const sigStart int = 0    // +1
const typeStart int = 1   // +1
const lengthStart int = 2 // +2
const payloadStart int = ConnOverhead

type SiteSig uint8     // 1
type SiteType uint8    // 1
type SiteLength uint16 // 2

const SiteSigVal PacketSig = 0xBB

const (
	SITE_TYPE_RAW    SiteType = 0x00 // []byte
	SITE_TYPE_ETH    SiteType = 0x01 // 
	SITE_TYPE_SITE   SiteType = 0x02 // 
	SITE_TYPE_ERROR  SiteType = 0xFF // []byte
)

type SitePacket struct {
	sSig    SiteSig
	sType   SiteType
	sLength SiteLength
	payload interface{}
}

// NewPacket()
// Packets coming from the low level connections
func NewSite(t SiteType, payload interface{}) (cp *SitePacket, err error) {
	var plen int

	switch payload.(type) {
	case nil:
		plen = 0
	case []byte:
		plen = len(payload.([]byte))
	case counter.Counter64:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(payload.(counter.Counter64)))
		payload = b
		plen = 8
	case *packetchan.ChanPacket:
		plen = payload.(*packetchan.ChanPacket).GetSize()
	default:
		log.Errorf("Bad Payload Type:%t", payload)
		return nil, errors.ErrConnPayloadType
	}

	cp = &ConnPacket{
		pSig:    ConnSigVal,
		pType:   t,
		pLength: PacketLength(plen),
		payload: payload,
	}

	return cp, nil
}

func (cp *ConnPacket) GetType() (t PacketType) {
	return cp.pType
}

func (cp *ConnPacket) GetSize() (l int) {
	return int(cp.pLength) + ConnOverhead
}

func (cp *ConnPacket) GetPayloadLength() (l PacketLength) {
	return cp.pLength
}

func (cp *ConnPacket) GetPayload() (payload interface{}) {
	switch cp.payload.(type) {
	case nil:
		return nil

	case []byte:
		payload = make([]byte, len(cp.payload.([]byte)))
		copy(payload.([]byte), cp.payload.([]byte))

	case *packetchan.ChanPacket:
		payload = cp.payload.(*packetchan.ChanPacket).Copy()

	default:
		log.Fatalf("Unhandled Type:%t", cp.payload)
	}
	return payload
}

func (cp *ConnPacket) Copy() (copy *ConnPacket) {
	copy = &ConnPacket{
		pSig:    ConnSigVal,
		pType:   cp.pType,
		pLength: cp.pLength,
		payload: cp.GetPayload(),
	}
	return copy
}

func MakePacket(data []byte) (cp *ConnPacket, err error) {

	if len(data) < (ConnOverhead) {
		log.Debugf("Short Packet data:%d < %d", len(data), ConnOverhead)
		return nil, errors.ErrConnShortPacket
	}

	if PacketSig(data[sigStart]) != ConnSigVal {
		return nil, errors.ErrChanBadSig
	}

	cp = &ConnPacket{
		pSig:    ConnSigVal,
		pType:   PacketType(data[typeStart]),
		pLength: PacketLength(binary.LittleEndian.Uint16(data[lengthStart : lengthStart+2])),
	}

	var payloadlen PacketLength

	switch cp.pType {
	case SITE_TYPE_RAW:
		cp.payload = data[ConnOverhead:]

	case SITE_TYPE_CHAN:

		ch, err := packetchan.MakePacket(data[ConnOverhead:])
		if err != nil {
			return nil, err
		}

		payloadlen = PacketLength(ch.GetSize())
		cp.payload = ch

	case SITE_TYPE_PING64:
		fallthrough
	case SITE_TYPE_PONG64:
		if payloadlen != 8 {
			log.Fatalf("Bad PING-PONG payload len:%d", payloadlen)
		}
		cp.payload = counter.Counter64(binary.LittleEndian.Uint64(data[ConnOverhead:]))

	case SITE_TYPE_AUTH:
		fallthrough
	default:
		log.Debugf("Bad Packet type:%d", cp.pType)
		return nil, errors.ErrConnBadPacket
	}

	if payloadlen != cp.pLength {
		log.Debugf("Bad Packet length:%d != %d", payloadlen, uint16(cp.pLength))
		return nil, errors.ErrConnPayloadLength
	}

	return cp, nil
}

func (p *ConnPacket) ToByte() (b []byte) {
	// Signature
	b = append(b, byte(ConnSigVal))
	// Type
	b = append(b, byte(p.pType))
	// Length
	len := make([]byte, 2)
	binary.LittleEndian.PutUint16(len, uint16(p.pLength))
	b = append(b, len...)
	// Payload
	switch p.payload.(type) {
	case nil:
	case []byte:
		b = append(b, p.payload.([]byte)...)
	case *packetchan.ChanPacket:
		b = append(b, p.payload.(*packetchan.ChanPacket).ToByte()...)
	default:
		log.Fatalf("Unhandled Type:%t", p.payload)
	}

	return b
}

func (p PacketType) String() string {
	switch p {
	case SITE_TYPE_CHAN:
		return "CHAN"
	case SITE_TYPE_AUTH:
		return "AUTH"
	case SITE_TYPE_PING64:
		return "PING64"
	case SITE_TYPE_PONG64:
		return "PONG64"
	case SITE_TYPE_RAW:
		return "RAW"
	case SITE_TYPE_ERROR:
		return "ERROR"
	default:
		return fmt.Sprintf("UNKNOWN TYPE:%s", p.String())

	}
}
