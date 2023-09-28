package link

import "github.com/seanmcadam/octovpn/octolib/log"

func (ls *LinkStateStruct) goRecv() {
	if ls == nil {
		return
	}
	defer ls.Cancel()

	log.Debugf("Starting:%s", ls.linkname)

	for {
		select {
		case <-ls.processCh:
			ls.sendState()

		case <-ls.doneChan():
			return
		}
	}
}
