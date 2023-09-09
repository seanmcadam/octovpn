package link

import "github.com/seanmcadam/octovpn/octolib/log"

// Send      UP     DOWN    OTHERS
// Cond   --------------------------
// ANY    | Send    Send    Send
// UP AND | ALL UP  Send    Send    // All are UP to be UP - CONN Auth/Connection Example
// DN AND | Send    ALL DN  Send    // ALL are DN to be DN - Site/Channel example
// UP OR  | Send    ALL!UP  ALL!UP  // ANY are UP to be UP
// DN OR  | All!DN  Send    ALL!DN  // ANY are DN to be DN
//
// PassALL - Pass All messages
// UpAND - Any non-UP message are sent - Send UP is all others are up   All UP to be UP
// UpOR -  Any UP Messages are sent - All others must not be UP to send Any UP to be UP
// DownAND - All connections must be DOWN for an DOWN - All other messages others are sent
// DownOR - Any DOWN Messages are sent - All others must not be DOWN to send
//

// State has already assumed to have changed
func (ls *LinkStateStruct) processMessage(ns LinkNoticeStateType) {
	statechange := (ns.State() != LinkStateNONE)
	noticechange := (ns.Notice() != LinkNoticeNONE)

	if noticechange {
		log.Debugf("Link[%d] Notice:%s",ls.instance, ns)
		ls.processNotice(noticeState(ns.Notice(), LinkStateNONE))
	}
	if statechange {
		ls.processStateChange(noticeState(LinkNoticeNONE, ns.State()))
	}

}

func (ls *LinkStateStruct) processNotice(ns LinkNoticeStateType) {

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

func (ls *LinkStateStruct) processState(ns LinkNoticeStateType) {
	switch ns.State() {
	case LinkStateUP:
		ls.linkUp.send(ns)
	case LinkStateLINK:
		ls.linkLink.send(ns)
	case LinkStateDOWN:
		ls.linkDown.send(ns)
	}
}

func (ls *LinkStateStruct) processStateChange(ns LinkNoticeStateType) {

	switch ls.mode {
	// Send all packets Always
	case LinkModePassALL:
		ls.linkNoticeState.send(ns)
		ls.linkState.send(ns)
		ls.processState(ns)

		// Always Send UP
		//
	case LinkModeUpAND:
		if ns.State() != LinkStateUP {
			ls.linkNoticeState.send(ns)
			ls.linkState.send(ns)
			ls.processState(ns)
		} else {
			// State UP, send if all UP
			allup := true
			for _, v := range ls.recvstate {
				if v != LinkStateUP {
					allup = false
					break
				}
			}
			if allup {
				ls.linkNoticeState.send(ns)
				ls.linkState.send(ns)
				ls.processState(ns)
			}
		}

	case LinkModeDownAND:
		if ns.State() != LinkStateDOWN {
			ls.linkNoticeState.send(ns)
			ls.linkState.send(ns)
			ls.processState(ns)
		} else {
			// State UP, send if all UP
			alldown := true
			for _, v := range ls.recvstate {
				if v != LinkStateDOWN {
					alldown = false
					break
				}
			}
			if alldown {
				ls.linkNoticeState.send(ns)
				ls.linkState.send(ns)
				ls.processState(ns)
			}
		}

	case LinkModeUpOR:
		if ns.State() == LinkStateUP {
			ls.linkNoticeState.send(ns)
			ls.linkState.send(ns)
			ls.processState(ns)
		} else {
			// State UP, send if all Not UP
			allnotup := true
			for _, v := range ls.recvstate {
				if v == LinkStateUP {
					allnotup = false
					break
				}
			}
			if allnotup {
				ls.linkNoticeState.send(ns)
				ls.linkState.send(ns)
				ls.processState(ns)
			}

		}
	case LinkModeDownOR:
		if ns.State() == LinkStateDOWN {
			ls.linkNoticeState.send(ns)
			ls.linkState.send(ns)
			ls.processState(ns)
		} else {
			// State Down, send if all not Down
			allnotdown := true
			for _, v := range ls.recvstate {
				if v == LinkStateDOWN {
					allnotdown = false
					break
				}
			}
			if allnotdown {
				ls.linkNoticeState.send(ns)
				ls.linkState.send(ns)
				ls.processState(ns)
			}
		}
	}

}

func (ls *LinkStateStruct) LinkStateCh() (newch LinkNoticeStateCh) {
	return ls.linkState.LinkCh()
}

func (ls *LinkStateStruct) LinkNoticeCh() (newch LinkNoticeStateCh) {
	return ls.linkNotice.LinkCh()
}

func (ls *LinkStateStruct) LinkNoticeStateCh() (newch LinkNoticeStateCh) {
	return ls.linkNoticeState.LinkCh()
}

func (ls *LinkStateStruct) LinkLinkCh() (newch LinkNoticeStateCh) {
	return ls.linkLink.LinkCh()
}

func (ls *LinkStateStruct) LinkUpCh() (newch LinkNoticeStateCh) {
	return ls.linkUp.LinkCh()
}

func (ls *LinkStateStruct) LinkDownCh() (newch LinkNoticeStateCh) {
	return ls.linkDown.LinkCh()
}

func (ls *LinkStateStruct) LinkCloseCh() (newch LinkNoticeStateCh) {
	return ls.linkClose.LinkCh()
}

func (ls *LinkStateStruct) LinkLossCh() (newch LinkNoticeStateCh) {
	return ls.linkLoss.LinkCh()
}

func (ls *LinkStateStruct) LinkLatencyCh() (newch LinkNoticeStateCh) {
	return ls.linkLatency.LinkCh()
}

func (ls *LinkStateStruct) LinkSaturationCh() (newch LinkNoticeStateCh) {
	return ls.linkSaturation.LinkCh()
}
