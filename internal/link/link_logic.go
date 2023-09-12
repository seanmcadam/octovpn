package link

import "github.com/seanmcadam/octovpn/octolib/log"

// Msg       UP     DOWN
// Cond   -----------------
// ANY    | Send    Send
// UP AND | ALL UP  Send   // All are UP to be UP - CONN Auth/Connection Example
// DN AND | Send    ALL DN // ALL are DN to be DN - Site/Channel example
// UP OR  | Send    ALL DN // ANY are UP to be UP
// DN OR  | All UP  Send   // ANY are DN to be DN
//
// PassALL - Pass All messages
// UpAND   - All DOWN message are sent - All UP to send UP
// DownOR  - same
// UpOR    - ALL UP Messages are sent  - All DOWN to send DOWN
// DownAND - same
//

// State has already assumed to have changed
func (ls *LinkStateStruct) processMessage(ns LinkNoticeStateType) {
	if ls == nil{
		return 
	}

	statechange := (ns.State() != LinkStateNONE)
	noticechange := (ns.Notice() != LinkNoticeNONE)

	if noticechange {
		log.Debugf("Link[%d] Notice:%s", ls.instance, ns)
		ls.processNotice(noticeState(ns.Notice(), LinkStateNONE))
	}
	if statechange {
		ls.processStateChange(noticeState(LinkNoticeNONE, ns.State()))
	}

}

func (ls *LinkStateStruct) processNotice(ns LinkNoticeStateType) {
	if ls == nil{
		return 
	}

	ls.linkNotice.send(ns)
	ls.linkNoticeState.send(ns)

	switch ns.Notice() {
	case LinkNoticeCLOSED:
		ls.linkClose.send(ns)
	case LinkNoticeLOSS:
		ls.linkLoss.send(ns)
	case LinkNoticeLATENCY:
		ls.linkLatency.send(ns)
	case LinkNoticeSATURATED:
		ls.linkSaturation.send(ns)
	default:
	}
}

func (ls *LinkStateStruct) processStateChange(ns LinkNoticeStateType) {
	if ls == nil{
		return 
	}

	currentState := ls.state
	newState := ns.State()
	if currentState == newState {
		return
	}

	var sendTransition = false

	switch ls.mode {
	case LinkModePassALL:
		//
		// Send all packets Always
		//
		sendTransition = true

	case LinkModeUpAND:
		fallthrough
	case LinkModeDownOR:
		//
		// Msg is Down so Send packet
		//
		if (ns.State() & LinkStateDownMASK) > 0 {
			sendTransition = true
		} else {
			//
			// Send Msg UP All links UP
			//
			allup := true
			for _, v := range ls.recvstate {
				if (v & LinkStateUpMASK) == 0 {
					allup = false
					break
				}
			}

			if allup {
				sendTransition = true
			}
		}

	case LinkModeUpOR:
		fallthrough
	case LinkModeDownAND:
		//
		// Msg is UP so Send packet
		//
		if (ns.State() & LinkStateUpMASK) > 0 {
			sendTransition = true
		} else {
			//
			// Send if All Links DOWN
			//
			alldown := true
			for _, v := range ls.recvstate {
				if (v & LinkStateDownMASK) == 0 {
					alldown = false
					break
				}
			}

			if alldown {
				sendTransition = true
			}
		}

	}

	if sendTransition {

		ls.state = newState

		ls.linkNoticeState.send(ns)
		ls.linkState.send(ns)
		switch ns.State() {
		case LinkStateCONNECTED:
			ls.linkConnected.send(ns)
		case LinkStateAUTH:
			ls.linkAuth.send(ns)
		case LinkStateCHAL:
			ls.linkChal.send(ns)
		case LinkStateLINK:
			ls.linkLink.send(ns)
		case LinkStateNOLINK:
			ls.linkNoLink.send(ns)
		}

		if (currentState&LinkStateUpMASK) > 0 && (newState&LinkStateDownMASK) > 0 {
			// UP -> DOWN
			ls.linkDown.send(ns)
		} else if (currentState&LinkStateDownMASK) > 0 && (newState&LinkStateUpMASK) > 0 {
			// DOWN -> UP
			ls.linkUp.send(ns)
		}

	}
}

func (ls *LinkStateStruct) LinkStateCh() (newch LinkNoticeStateCh) {
	if ls == nil {
		return nil
	}

	return ls.linkState.LinkCh()
}

func (ls *LinkStateStruct) LinkNoticeCh() (newch LinkNoticeStateCh) {
	if ls == nil {
		return nil
	}

	return ls.linkNotice.LinkCh()
}

func (ls *LinkStateStruct) LinkNoticeStateCh() (newch LinkNoticeStateCh) {
	if ls == nil {
		return nil
	}

	return ls.linkNoticeState.LinkCh()
}

func (ls *LinkStateStruct) LinkChalCh() (newch LinkNoticeStateCh) {
	if ls == nil {
		return nil
	}

	return ls.linkChal.LinkCh()
}

func (ls *LinkStateStruct) LinkAuthCh() (newch LinkNoticeStateCh) {
	if ls == nil {
		return nil
	}

	return ls.linkAuth.LinkCh()
}

func (ls *LinkStateStruct) LinkLinkCh() (newch LinkNoticeStateCh) {
	if ls == nil {
		return nil
	}

	return ls.linkLink.LinkCh()
}

func (ls *LinkStateStruct) LinkConnectCh() (newch LinkNoticeStateCh) {
	if ls == nil {
		return nil
	}

	return ls.linkConnected.LinkCh()
}

func (ls *LinkStateStruct) LinkUpDownCh() (newch LinkNoticeStateCh) {
	if ls == nil {
		return nil
	}

	return ls.linkUp.LinkCh()
}

func (ls *LinkStateStruct) LinkUpCh() (newch LinkNoticeStateCh) {
	if ls == nil {
		return nil
	}

	return ls.linkUp.LinkCh()
}

func (ls *LinkStateStruct) LinkDownCh() (newch LinkNoticeStateCh) {
	if ls == nil {
		return nil
	}

	return ls.linkDown.LinkCh()
}

func (ls *LinkStateStruct) LinkCloseCh() (newch LinkNoticeStateCh) {
	if ls == nil {
		return nil
	}

	return ls.linkClose.LinkCh()
}

func (ls *LinkStateStruct) LinkLossCh() (newch LinkNoticeStateCh) {
	if ls == nil {
		return nil
	}

	return ls.linkLoss.LinkCh()
}

func (ls *LinkStateStruct) LinkLatencyCh() (newch LinkNoticeStateCh) {
	if ls == nil {
		return nil
	}

	return ls.linkLatency.LinkCh()
}

func (ls *LinkStateStruct) LinkSaturationCh() (newch LinkNoticeStateCh) {
	if ls == nil {
		return nil
	}

	return ls.linkSaturation.LinkCh()
}

func (ls *LinkStateStruct) refreshRecvLinks() {
	if ls == nil {
		return
	}
	ls.recvnew <- noticeState(LinkNoticeNONE, LinkStateNONE)
}

func (ls *LinkStateStruct) AddLinkStateCh(link *LinkStateStruct) {
	if ls == nil {
		return
	}
	state := link.GetState()

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkStateCh,
	}

	ls.addlinkch <- add
	ls.refreshRecvLinks()

}
func (ls *LinkStateStruct) AddLinkNoticeCh(link *LinkStateStruct) {
	if ls == nil {
		return
	}
	state := LinkStateNONE

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkNoticeCh,
	}

	ls.addlinkch <- add
	ls.refreshRecvLinks()

}
func (ls *LinkStateStruct) AddLinkNoticeStateCh(link *LinkStateStruct) {
	if ls == nil {
		return
	}
	state := link.GetState()

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkNoticeStateCh,
	}

	ls.addlinkch <- add
	ls.refreshRecvLinks()

}
func (ls *LinkStateStruct) AddLinkLinkCh(link *LinkStateStruct) {
	if ls == nil {
		return
	}
	state := LinkStateNONE

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkLinkCh,
	}

	ls.addlinkch <- add
	ls.refreshRecvLinks()

}
func (ls *LinkStateStruct) AddLinkUpDownCh(link *LinkStateStruct) {
	if ls == nil {
		return
	}
	state := link.GetState()

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkUpDownCh,
	}

	ls.addlinkch <- add
	ls.refreshRecvLinks()

}
func (ls *LinkStateStruct) AddLinkUpCh(link *LinkStateStruct) {
	if ls == nil {
		return
	}
	state := LinkStateNONE

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkUpCh,
	}

	ls.addlinkch <- add
	ls.refreshRecvLinks()

}
func (ls *LinkStateStruct) AddLinkConnectCh(link *LinkStateStruct) {
	if ls == nil {
		return
	}
	state := LinkStateNONE

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkConnectCh,
	}

	ls.addlinkch <- add
	ls.refreshRecvLinks()

}
func (ls *LinkStateStruct) AddLinkDownCh(link *LinkStateStruct) {
	if ls == nil {
		return
	}
	state := LinkStateNONE

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkDownCh,
	}

	ls.addlinkch <- add
	ls.refreshRecvLinks()

}
func (ls *LinkStateStruct) AddLinkCloseCh(link *LinkStateStruct) {
	if ls == nil {
		return
	}
	state := LinkStateNONE

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkCloseCh,
	}

	ls.addlinkch <- add
	ls.refreshRecvLinks()

}
func (ls *LinkStateStruct) AddLinkLossCh(link *LinkStateStruct) {
	if ls == nil {
		return
	}
	state := LinkStateNONE

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkLossCh,
	}

	ls.addlinkch <- add
	ls.refreshRecvLinks()

}
func (ls *LinkStateStruct) AddLinkLatencyCh(link *LinkStateStruct) {
	if ls == nil {
		return
	}
	state := LinkStateNONE

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkLatencyCh,
	}

	ls.addlinkch <- add
	ls.refreshRecvLinks()

}
func (ls *LinkStateStruct) AddLinkSaturationCh(link *LinkStateStruct) {
	if ls == nil {
		return
	}
	state := LinkStateNONE

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkSaturationCh,
	}

	ls.addlinkch <- add
	ls.refreshRecvLinks()

}
