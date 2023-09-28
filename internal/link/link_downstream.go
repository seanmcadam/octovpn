package link

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
	return ls.linkUpDown.LinkCh()
}

func (ls *LinkStateStruct) LinkListenCh() (newch LinkNoticeStateCh) {
	if ls == nil {
		return nil
	}
	return ls.linkListen.LinkCh()
}

func (ls *LinkStateStruct) LinkNoLinkCh() (newch LinkNoticeStateCh) {
	if ls == nil {
		return nil
	}
	return ls.linkNoLink.LinkCh()
}

func (ls *LinkStateStruct) LinkUpCh() (newch LinkNoticeStateCh) {
	if ls == nil {
		return nil
	}
	return ls.linkUp.LinkCh()
}

func (ls *LinkStateStruct) LinkStartCh() (newch LinkNoticeStateCh) {
	if ls == nil {
		return nil
	}
	return ls.linkStart.LinkCh()
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
