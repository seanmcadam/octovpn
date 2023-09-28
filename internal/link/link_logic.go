package link

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
func (ls *LinkStateStruct) processIncomingMessage(ns *LinkNoticeStateType) {
	if ls == nil {
		return
	}

	statechange := (ns.State() != LinkStateNONE)
	noticechange := (ns.Notice() != LinkNoticeNONE)

	if noticechange {
		// log.Debugf("[%s] Notice:%s", ls.linkname, ns)
		ls.processNotice(noticeState(ns.Notice(), LinkStateNONE))
	}

	if statechange {
		ls.processStateChange(noticeState(LinkNoticeNONE, ns.State()))
	}

}

// processNotice()
// Just send it
func (ls *LinkStateStruct) processNotice(ns *LinkNoticeStateType) {
	if ls == nil {
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

//
// processStateChange()
//
func (ls *LinkStateStruct) processStateChange(ns *LinkNoticeStateType) {
	if ls == nil {
		return
	}

	newState := ns.State()

	if newState == ls.state {
		return
	}

	switch ls.mode {
	case LinkModePassALL:
		//
		// Send all packets Always
		//
		ls.setState(newState)

	case LinkModeUpAND:
		fallthrough
	case LinkModeDownOR:
		//
		// Msg is Down so Send packet
		//
		if (ns.State() & LinkStateDownMASK) > 0 {
			ls.setState(newState)
		} else {
			//
			// Send Msg UP All links UP
			//
			allup := true
			for _, v := range ls.getStates() {
				if (v & LinkStateUpMASK) == 0 {
					allup = false
					break
				}
			}

			if allup {
				ls.setState(newState)
			}
		}

	case LinkModeUpOR:
		fallthrough
	case LinkModeDownAND:
		//
		// Msg is UP so Send packet
		//
		if (ns.State() & LinkStateUpMASK) > 0 {
			ls.setState(newState)
		} else {
			//
			// Send if All Links DOWN
			//
			alldown := true
			for _, v := range ls.getStates() {
				if (v & LinkStateDownMASK) == 0 {
					alldown = false
					break
				}
			}

			if alldown {
				ls.setState(newState)
			}
		}
	}
}
