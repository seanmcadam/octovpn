package link

import (
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

const DefaultChanDepth = 64
const LinkChanDepth = 1 // Always allow 1 message to be send, no waiting...

func newLinkChan(name string) (lc *LinkChan) {
	lc = &LinkChan{
		name:       name,
		listenChan: make(chan LinkNoticeStateListenCh, DefaultChanDepth),
	}

	return lc
}

func (lc *LinkChan) LinkCh() (newch LinkNoticeStateCh) {
	if lc == nil {
		log.FatalfStack("%s", errors.ErrorNilMethodPointer())
	}

	ch := make(chan *LinkNoticeStateType, LinkChanDepth)
	select {
	case lc.listenChan <- ch:
		return ch
	default:
		log.FatalfStack("Channel Full - Link:%s chan len:%d", lc.name, len(lc.listenChan))
	}

	return nil
}

func (lc *LinkChan) send(ns *LinkNoticeStateType) {
	if lc == nil {
		log.ErrorStack("send() nil method pointer")
		return
	}

	if lc.listenChan == nil {
		log.Warnf("send() listenChan[%s] is nil", lc.name)
		return
	}

	//log.GDebugf("LinkChan send()[%s] count %d", lc.name, len(lc.listenChan))

	if len(lc.listenChan) > 0 {
		lc.sendmx.Lock()
		defer lc.sendmx.Unlock()

		// log.GDebugf("LinkChan send() listenChan[%s] count:%d, Msg:%s", lc.name, len(lc.listenChan), ns)
		length := len(lc.listenChan)
		for i := 0; i < length; i++ {
			var l LinkNoticeStateListenCh
			select {
			case l = <-lc.listenChan:
			default:
				log.ErrorfStack("LinkChan[%s] unable unable to get next chan %s", lc.name, ns)
			}

			select {
			case l <- ns:
			default:
				log.ErrorfStack("LinkChan[%s] unable unable to send %s", lc.name, ns)
			}
			close(l)
		}
	}
}
