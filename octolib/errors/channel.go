package errors

import "fmt"

var ErrChanBadSig error
var ErrChanBadPacket error
var ErrChanShortPacket error
var ErrChanPayloadLength error

func init() {
	ErrChanBadSig = fmt.Errorf("bad chan signature type")
	ErrChanBadPacket = fmt.Errorf("bad chan packet type")
	ErrChanShortPacket = fmt.Errorf("short chan packet")
	ErrChanPayloadLength = fmt.Errorf("bad chan payload length")
}
