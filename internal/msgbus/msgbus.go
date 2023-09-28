package msgbus

import (
	"github.com/asaskevich/EventBus"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type MsgTarget string
type MsgNotice string
type MsgState string

const (
	CloseNotice     MsgNotice = "Close"
	LossNotice      MsgNotice = "Loss"
	LatencyNotice   MsgNotice = "Latency"
	SaturatedNotice MsgNotice = "Saturated"
	StateNONE       MsgState  = "None"
	StateLISTEN     MsgState  = "Listen"
	StateNOLINK     MsgState  = "NoLink"
	StateSTART      MsgState  = "Start"
	StateLINK       MsgState  = "Link"
	StateChal       MsgState  = "Chal"
	StateAUTH       MsgState  = "Auth"
	StateCONNECTED  MsgState  = "Connected"
	StateERROR      MsgState  = "Error"
)

type MsgBus struct {
	eventbus EventBus.Bus
}

type NoticeStruct struct {
	Notice MsgNotice
	From   MsgTarget
}

type StateStruct struct {
	State MsgState
	From  MsgTarget
}

// -
// New()
// -
func New() (mb *MsgBus) {
	mb = &MsgBus{
		eventbus: EventBus.New(),
	}

	return mb
}

// -
// Send()
// -
func (mb *MsgBus) Send(target MsgTarget, data interface{}) {
	log.Infof("Send() to %s:%T", target, data)
	mb.eventbus.Publish(string(target), data)
}

// -
// ReceiveHandler()
// -
func (mb *MsgBus) ReceiveHandler(source MsgTarget, handler func(...interface{})) (err error) {
	return mb.eventbus.SubscribeAsync(string(source), handler, false)
}

// -
// SendState()
// -
func (mb *MsgBus) SendState(from MsgTarget, to MsgTarget, state MsgState) {
	msg := &StateStruct{
		State: state,
		From:  from,
	}

	mb.Send(to, msg)
}

// -
// SendNotice()
// -
func (mb *MsgBus) SendNotice(from MsgTarget, to MsgTarget, notice MsgNotice) {
	msg := &NoticeStruct{
		Notice: notice,
		From:   from,
	}

	mb.Send(to, msg)
}
