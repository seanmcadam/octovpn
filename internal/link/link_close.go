package link

import "github.com/seanmcadam/octovpn/octolib/log"

func (l *LinkStateStruct) doneChan() <-chan struct{} {
	if l == nil {
		return nil
	}
	return l.cx.DoneChan()
}

func (l *LinkStateStruct) Cancel() {
	if l == nil {
		return
	}
	log.GDebugf("Cancel() %s", l.linkname)
	l.cx.Cancel()
}

func (l *LinkStateStruct) closed() bool {
	if l == nil {
		return true
	}
	return l.cx.Done()
}
