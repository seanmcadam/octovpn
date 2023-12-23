package msg

import (
	"testing"

	"github.com/seanmcadam/octovpn/internal/interfaces"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/instance"
)

func Test_(t *testing.T) {

	var msg interfaces.MsgInterface

	in := instance.NewInstanceName("testing")

	msg = NewState(in, StateNONE)
	msg = NewNotice(in, NoticeCLOSED)
	msg = NewPacket(in, &packet.PacketStruct{})

	_ = msg
}
