package link

import (
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (ls *LinkStateStruct) Auth() {
	if ls == nil {
		return
	}
	ls.setState(LinkStateAUTH)
}

func (ls *LinkStateStruct) Chal() {
	if ls == nil {
		return
	}
	ls.setState(LinkStateCHAL)
}

func (ls *LinkStateStruct) Close() {
	if ls == nil {
		return
	}
	ls.SendNotice(LinkNoticeCLOSED)
}

func (ls *LinkStateStruct) Connected() {
	if ls == nil {
		return
	}
	ls.setState(LinkStateCONNECTED)
}

func (ls *LinkStateStruct) Link() {
	if ls == nil {
		return
	}
	ls.setState(LinkStateLINK)
}

func (ls *LinkStateStruct) Listen() {
	if ls == nil {
		return
	}
	ls.setState(LinkStateLISTEN)
}

func (ls *LinkStateStruct) NoLink() {
	if ls == nil {
		return
	}
	ls.setState(LinkStateNOLINK)
}

func (ls *LinkStateStruct) Start() {
	if ls == nil {
		return
	}
	ls.setState(LinkStateSTART)
}

func (ls *LinkStateStruct) setState(s LinkStateType) {
	if ls == nil {
		return
	}

	if s == ls.state {
		log.GDebugf("No Link State Change[%s]:%s -> %s", ls.linkname, ls.state, s)
		return
	}

	ls.laststate = ls.state
	ls.state = s

	log.GDebugf(" Link State Change[%s]:%s -> %s", ls.linkname, ls.state, s)
	ls.processCh <- noticeState(LinkNoticeNONE, s)
}

func (ls *LinkStateStruct) IsUp() bool {
	s := ls.state & LinkStateUpMASK
	return s > 0
	//return ls.state & LinkStateUpMASK > 0
}

func (ls *LinkStateStruct) IsDown() bool {
	return ls.state&LinkStateDownMASK > 0
}

func (ls *LinkStateStruct) GetState() LinkStateType {
	if ls == nil {
		return 0
	}
	return ls.state
}

func (ls *LinkStateStruct) sendState() {

	log.GDebugf("Send State[%s]: %s", ls.linkname, ls.state)

	ns := noticeState(LinkNoticeNONE, ls.state)

	ls.linkNoticeState.send(ns)
	ls.linkState.send(ns)

	switch ns.State() {
	case LinkStateCONNECTED:
		ls.linkConnected.send(ns)
		ls.linkUp.send(ns)
		ls.linkUpDown.send(ns)
	case LinkStateAUTH:
		ls.linkAuth.send(ns)
		ls.linkUp.send(ns)
		ls.linkUpDown.send(ns)
	case LinkStateCHAL:
		ls.linkChal.send(ns)
		ls.linkUp.send(ns)
		ls.linkUpDown.send(ns)
	case LinkStateLINK:
		ls.linkLink.send(ns)
		ls.linkUp.send(ns)
		ls.linkUpDown.send(ns)
	case LinkStateLISTEN:
		ls.linkListen.send(ns)
		ls.linkDown.send(ns)
		ls.linkUpDown.send(ns)
	case LinkStateSTART:
		ls.linkStart.send(ns)
		ls.linkDown.send(ns)
		ls.linkUpDown.send(ns)
	case LinkStateNOLINK:
		ls.linkNoLink.send(ns)
		ls.linkDown.send(ns)
		ls.linkUpDown.send(ns)
	}

	if (ls.laststate&LinkStateUpMASK) > 0 && (ls.state&LinkStateDownMASK) > 0 {
		// UP -> DOWN
		ls.linkDown.send(ns)
	} else if (ls.laststate&LinkStateDownMASK) > 0 && (ls.state&LinkStateUpMASK) > 0 {
		// DOWN -> UP
		ls.linkUp.send(ns)
	}
}
