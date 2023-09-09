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
	if len(ls.listenChan) > 0 {
		log.Debugf("LinkChan[%s] Send Message:%s", ls.name, ns)
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
