package msg

import (
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/octolib/instance"
)

type MsgState string

const (
	StateNONE      MsgState = "None"
	StateLISTEN    MsgState = "Listen"
	StateNOLINK    MsgState = "NoLink"
	StateSTART     MsgState = "Start"
	StateLINK      MsgState = "Link"
	StateCHAL      MsgState = "Chal"
	StateAUTH      MsgState = "Auth"
	StateCONNECTED MsgState = "Connected"
	StateERROR     MsgState = "Error"
)

type StateHandleTable struct {
	handler map[MsgState]func(*StateStruct)
}

func NewStateHandler() (sht *StateHandleTable) {
	sht = &StateHandleTable{
		handler: map[MsgState]func(*StateStruct){},
	}
	sht.handler[StateNONE] = emptyStateHandle
	sht.handler[StateLISTEN] = emptyStateHandle
	sht.handler[StateNOLINK] = emptyStateHandle
	sht.handler[StateSTART] = emptyStateHandle
	sht.handler[StateLINK] = emptyStateHandle
	sht.handler[StateCHAL] = emptyStateHandle
	sht.handler[StateAUTH] = emptyStateHandle
	sht.handler[StateCONNECTED] = emptyStateHandle
	sht.handler[StateERROR] = emptyStateHandle
	return sht
}

func (nht *StateHandleTable) Run(ss *StateStruct) {
	if handlerFn, ok := nht.handler[ss.State]; ok {
		handlerFn(ss)
	} else {
		log.FatalfStack("Unknown Run Notice Type %v", ss)
	}
}

func (sht *StateHandleTable) AddHandle(ms MsgState, fn func(*StateStruct)) {
	sht.handler[ms] = fn
}

func (sht *StateHandleTable) CallHandle(ss *StateStruct) {
	sht.handler[ss.State](ss)
}

type StateStruct struct {
	State MsgState
	From  *instance.InstanceName
}

func NewState(from *instance.InstanceName, state MsgState) (ss *StateStruct) {
	ss = &StateStruct{
		State: state,
		From:  from,
	}
	return ss
}

func (s *StateStruct) FromName() *instance.InstanceName {
	return s.From
}

func (n *StateStruct) Data() interface{} {
	return &n.State
}

func emptyStateHandle(ss *StateStruct) {
	log.ErrorfStack("StateEmptyHandler From:%s State:%s", ss.From, ss.State)
}
