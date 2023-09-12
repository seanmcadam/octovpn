package link

import "github.com/seanmcadam/octovpn/octolib/log"

const DefaultChanDepth = 64

func newLinkChan(name string) (lc *LinkChan) {
	lc = &LinkChan{
		name:       name,
		listenChan: make(chan LinkNoticeStateListenCh, DefaultChanDepth),
	}

	return lc
}

func (ls *LinkChan) LinkCh() (newch LinkNoticeStateCh) {
	if ls == nil {
		log.ErrorStack("send() nil method pointer")
		return
	}

	ch := make(chan LinkNoticeStateType, 1) // Always allow 1 message to be send, no waiting...
	select {
	case ls.listenChan <- ch:
		return ch
	default:
		log.Errorf("Channel Full:%v", ls.listenChan)
	}

	log.Debug("Link:%s chan len:%d", ls.name, len(ls.listenChan))

	return nil
}

func (lc *LinkChan) send(ns LinkNoticeStateType) {
	if lc == nil {
		log.ErrorStack("send() nil method pointer")
		return
	}

	if lc.listenChan == nil {
		log.Warnf("send() listenChan[%s] is nil", lc.name)
		return
	}

	if len(lc.listenChan) > 0 {
		log.GDebugf("LinkChan send() listenChan[%s] count:%d, Msg:%s", lc.name, len(lc.listenChan), ns)
		length := len(lc.listenChan)
		for i := 0; i < length; i++ {
			l := <-lc.listenChan
			select {
			case l <- ns:
			default:
				log.Warn("LinkChan[%s] unable to send", lc.name)
			}
			close(l)
		}
	}
}
