package netlib

import "fmt"

var ErrNetPacketTooBig error
var ErrNetChannelDown error

func init() {
	ErrNetPacketTooBig = fmt.Errorf("packet too big")
	ErrNetChannelDown = fmt.Errorf("channel down")
}
