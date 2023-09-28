package msgbus

import (
	"github.com/asaskevich/EventBus"
)

type MsgTarget string
type MsgTopics string

type MsgBus struct {
	eventbus EventBus.Bus
}

const (
	MsgPacket MsgTopics = "PACKET"
	MsgNotice MsgTopics = "NOTICE"
	MsgState  MsgTopics = "STATE"
)

// -
// New()
// -
func New() (mb *MsgBus) {
	mb = &MsgBus{
		eventbus: EventBus.New(),
	}

	return mb
}
