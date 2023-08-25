package chanconn

import "fmt"

var ErrChanConnBadPacket error
var ErrChanConnShortPacket error
var ErrChanConnPayloadLength error

func init() {
	ErrChanConnBadPacket = fmt.Errorf("bad packet type")
	ErrChanConnShortPacket = fmt.Errorf("short packet")
	ErrChanConnPayloadLength = fmt.Errorf("bad payload length")
}
