package msgbus

import "fmt"

type StateMsg uint8

const (
	StateNONE      StateMsg = 0x00
	StateLISTEN    StateMsg = 0x01
	StateNOLINK    StateMsg = 0x02
	StateSTART     StateMsg = 0x03
	StateLINK      StateMsg = 0x10
	StateCHAL      StateMsg = 0x20
	StateAUTH      StateMsg = 0x30
	StateCONNECTED StateMsg = 0x80
	StateERROR     StateMsg = 0xFF
	StateUpMASK    StateMsg = 0xF0
	StateDownMASK  StateMsg = 0x0F
)

type StateMsgStruct struct {
	State StateMsg
}

// -
// state()
// -
func (mb *MsgBus) State(target MsgTarget, state *StateMsg) {
	topic := fmt.Sprintf("%s:%s", string(target), string(MsgState))
	mb.eventbus.Publish(topic, state)
}

//func (mb *MsgBus) SetState(state StateMsg) {
//	ms := &StateMsgStruct{
//		State: state,
//	}
//	mb.State(ms)
//}

// -
// Register State Handler
// -
func (mb *MsgBus) StateHandler(handler func(...interface{})) (err error) {
	return mb.eventbus.Subscribe(string(MsgState), handler)
}
