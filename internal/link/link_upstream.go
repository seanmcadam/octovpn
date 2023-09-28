package link

import "github.com/seanmcadam/octovpn/octolib/errors"

func (ls *LinkStateStruct) AddLinkAuthCh(link *LinkStateStruct) (err error) {
	if ls == nil {
		return errors.ErrorNilMethodPointer()
	}
	state := link.GetState()

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkAuthCh,
	}

	return ls.addLink(add)

}

func (ls *LinkStateStruct) AddLinkCloseCh(link *LinkStateStruct) error {
	if ls == nil {
		return errors.ErrorNilMethodPointer()
	}
	state := link.GetState()

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkCloseCh,
	}

	return ls.addLink(add)

}

func (ls *LinkStateStruct) AddLinkChalCh(link *LinkStateStruct) error {
	if ls == nil {
		return errors.ErrorNilMethodPointer()
	}
	state := link.GetState()

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkChalCh,
	}

	return ls.addLink(add)

}

func (ls *LinkStateStruct) AddLinkConnectCh(link *LinkStateStruct) error {
	if ls == nil {
		return errors.ErrorNilMethodPointer()
	}
	state := link.GetState()

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkConnectCh,
	}

	return ls.addLink(add)

}

func (ls *LinkStateStruct) AddLinkDownCh(link *LinkStateStruct) error {
	if ls == nil {
		return errors.ErrorNilMethodPointer()
	}
	state := link.GetState()

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkDownCh,
	}

	return ls.addLink(add)

}

func (ls *LinkStateStruct) AddLinkLatencyCh(link *LinkStateStruct) error {
	if ls == nil {
		return errors.ErrorNilMethodPointer()
	}
	state := link.GetState()

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkLatencyCh,
	}

	return ls.addLink(add)

}

func (ls *LinkStateStruct) AddLinkLinkCh(link *LinkStateStruct) error {
	if ls == nil {
		return errors.ErrorNilMethodPointer()
	}
	state := link.GetState()

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkLinkCh,
	}

	return ls.addLink(add)

}

func (ls *LinkStateStruct) AddLinkListenCh(link *LinkStateStruct) error {
	if ls == nil {
		return errors.ErrorNilMethodPointer()
	}
	state := link.GetState()

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkListenCh,
	}

	return ls.addLink(add)

}

func (ls *LinkStateStruct) AddLinkLossCh(link *LinkStateStruct) error {
	if ls == nil {
		return errors.ErrorNilMethodPointer()
	}
	state := link.GetState()

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkLossCh,
	}

	return ls.addLink(add)

}

func (ls *LinkStateStruct) AddLinkNoLinkCh(link *LinkStateStruct) error {
	if ls == nil {
		return errors.ErrorNilMethodPointer()
	}
	state := link.GetState()

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkNoLinkCh,
	}

	return ls.addLink(add)

}

func (ls *LinkStateStruct) AddLinkNoticeCh(link *LinkStateStruct) error {
	if ls == nil {
		return errors.ErrorNilMethodPointer()
	}
	state := link.GetState()

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkNoticeCh,
	}

	return ls.addLink(add)

}

func (ls *LinkStateStruct) AddLinkNoticeStateCh(link *LinkStateStruct) error {
	if ls == nil {
		return errors.ErrorNilMethodPointer()
	}
	state := link.GetState()

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkNoticeStateCh,
	}

	return ls.addLink(add)

}

func (ls *LinkStateStruct) AddLinkSaturationCh(link *LinkStateStruct) error {
	if ls == nil {
		return errors.ErrorNilMethodPointer()
	}
	state := link.GetState()

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkSaturationCh,
	}

	return ls.addLink(add)

}

func (ls *LinkStateStruct) AddLinkStartCh(link *LinkStateStruct) error {
	if ls == nil {
		return errors.ErrorNilMethodPointer()
	}

	state := link.GetState()

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkStartCh,
	}

	return ls.addLink(add)

}

func (ls *LinkStateStruct) AddLinkStateCh(link *LinkStateStruct) error {
	if ls == nil {
		return errors.ErrorNilMethodPointer()
	}

	state := link.GetState()

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkStateCh,
	}

	return ls.addLink(add)

}

func (ls *LinkStateStruct) AddLinkUpCh(link *LinkStateStruct) error {
	if ls == nil {
		return errors.ErrorNilMethodPointer()
	}
	state := link.GetState()

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkUpCh,
	}

	return ls.addLink(add)

}

func (ls *LinkStateStruct) AddLinkUpDownCh(link *LinkStateStruct) error {
	if ls == nil {
		return errors.ErrorNilMethodPointer()
	}
	state := link.GetState()

	add := &AddLinkStruct{
		State:    state,
		LinkFunc: link.LinkUpDownCh,
	}

	return ls.addLink(add)

}
