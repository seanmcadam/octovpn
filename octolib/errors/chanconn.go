package errors

import "fmt"

var ErrConnBadPacket error
var ErrConnShortPacket error
var ErrConnPayloadLength error
var ErrConnPayloadType error

func init() {
	ErrConnBadPacket = fmt.Errorf("bad conn packet type")
	ErrConnShortPacket = fmt.Errorf("short conn packet")
	ErrConnPayloadLength = fmt.Errorf("bad conn payload length")
	ErrConnPayloadType = fmt.Errorf("bad conn payload type")
}
