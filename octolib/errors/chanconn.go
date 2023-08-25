package errors

import "fmt"

var ErrChanConnBadPacket error
var ErrChanConnShortPacket error
var ErrChanConnPayloadLength error

func init() {
	ErrChanConnBadPacket = fmt.Errorf("bad conn packet type")
	ErrChanConnShortPacket = fmt.Errorf("short conn packet")
	ErrChanConnPayloadLength = fmt.Errorf("bad conn payload length")
}
