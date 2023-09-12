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

// func (ls *LinkChan) LinkCh() (newch chan LinkNoticeStateType) {
func (ls *LinkChan) LinkCh() (newch LinkNoticeStateCh) {
	ch := make(chan LinkNoticeStateType, 1) // Always allow 1 message to be send, no waiting...
	select {
	case ls.listenChan <- ch:
		return ch
	default:
		log.Errorf("Channel Full:%v", ls.listenChan)
	}
	return nil
}

func (ls *LinkChan) send(ns LinkNoticeStateType) {
	if ls == nil {
		log.Error("send() nil method pointer")
		return
	}

	if ls.listenChan == nil {
		log.Warnf("send() listenChan[%s] is nil", ls.name)
		return
	}

	if len(ls.listenChan) > 0 {
		log.GDebugf("send() listenChan[%s] count:%d, Msg:%s", ls.name, len(ls.listenChan), ns)
		length := len(ls.listenChan)
		for i := 0; i < length; i++ {
			l := <-ls.listenChan
			select {
			case l <- ns:
			default:
				log.Warn("LinkChan[%s] unable to send", ls.name)
			}
			close(l)
		}
	}
}
