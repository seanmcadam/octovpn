package link

import (
	"fmt"

	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func newLinkState(ctx *ctx.Ctx, name string, m ...LinkModeType) (ls *LinkStateStruct) {
	ls = NewLinkState(ctx, m...)
	ls.linkname = ls.linkname + name
	return ls
}

func NewNameLinkState(ctx *ctx.Ctx, name string, m ...LinkModeType) (ls *LinkStateStruct) {
	ls = NewLinkState(ctx, m...)
	ls.linkname = ls.linkname + name
	return ls
}

func NewLinkState(ctx *ctx.Ctx, m ...LinkModeType) (ls *LinkStateStruct) {
	count := <-instanceCounter.GetCountCh()
	mode := LinkModePassALL // Default Pass All

	if len(m) > 0 {
		mode = m[0]
	}

	name := fmt.Sprintf("link#%d ", count.Uint().(uint32))

	ls = &LinkStateStruct{
		cx:              ctx,
		linkname:        name,
		mode:            mode,
		laststate:       LinkStateNONE,
		state:           LinkStateNONE,
		linkNoticeState: newLinkChan("NoticeState"),
		linkState:       newLinkChan("State"),
		linkNotice:      newLinkChan("Notice"),
		linkUpDown:      newLinkChan("UpDown"),
		linkUp:          newLinkChan("Up"),
		linkDown:        newLinkChan("Down"),
		linkLink:        newLinkChan("Link"),
		linkAuth:        newLinkChan("Auth"),
		linkChal:        newLinkChan("Chal"),
		linkNoLink:      newLinkChan("NoLink"),
		linkListen:      newLinkChan("Listen"),
		linkConnected:   newLinkChan("Connected"),
		linkLoss:        newLinkChan("Loss"),
		linkLatency:     newLinkChan("Latency"),
		linkSaturation:  newLinkChan("Saturation"),
		linkClose:       newLinkChan("Close"),
		linkStart:       newLinkChan("Start"),
		processCh:       make(chan *LinkNoticeStateType),
		recvcounter:     counter.NewCounter64(ctx),
		recvmap:         map[uint64]*AddLinkStruct{},
	}

	// log.Debugf("[%s] Starting", ls.linkname)
	go ls.goRecv()
	return ls
}

func (ls *LinkStateStruct) SendNotice(n LinkNoticeType) {
	if n == LinkNoticeNONE {
		return
	}
	ls.sendNotice(noticeState(n, LinkStateNONE))
}

// -
//
// -
func (l *LinkStateStruct) addLink(add *AddLinkStruct) (err error) {

	if l == nil {
		return errors.ErrorNilMethodPointer()
	}

	if add == nil {
		log.Fatal("Nil Add Pointer")
	}

	go func(l *LinkStateStruct, add *AddLinkStruct) {

		c := l.recvcounter.Next().Uint().(uint64)
		l.recvmx.Lock()
		l.recvmap[c] = add
		l.recvmx.Unlock()
		defer l.deleteLink(c)

		log.Debugf("[%s]Add Channel: %d", l.linkname, c)

		done := l.doneChan()
		var ch *LinkNoticeStateCh
		for {
			ch = func() (ret *LinkNoticeStateCh) {
				l.recvmx.Lock()
				defer l.recvmx.Unlock()

				log.Debugf("[%s]addLink() Locked Channel: %d", l.linkname, c)

				if add, isset := l.recvmap[c]; !isset {
					log.FatalfStack("Bad Map - Count:%d Recvmap:%v", c, l.recvmap)
				} else {
					if f := add.LinkFunc(); f == nil {
						log.FatalfStack("LinkFunc() Count:%d Recvmap:%v", c, l.recvmap)
					} else {
						return &f
					}
				}
				return nil
			}()
			log.Debugf("[%s]addLink() Channel UnLocked: %d", l.linkname, c)

			if ch == nil {
				return
			}

			var state *LinkNoticeStateType
			select {
			case <-done:
				return

			case state = <-*ch:
				if state == nil {
					return
				}

				l.processIncomingMessage(state)
			}
		}
	}(l, add)

	return nil
}

// -
//
// -
func (l *LinkStateStruct) deleteLink(c uint64) {

	if l == nil {
		log.Fatal("Nil Method Pointer")
	}

	l.recvmx.Lock()
	defer l.recvmx.Unlock()

	log.Debugf("[%s]Delete Channel: %d", l.linkname, c)

	if _, isset := l.recvmap[c]; !isset {
		log.Fatalf("[%s]Delete Bad Channel: %d", l.linkname, c)
		return
	}

	delete(l.recvmap, c)
}

// -
//
// -
func (l *LinkStateStruct) getStates() (arr []LinkStateType) {

	if l == nil {
		log.Fatal("Nil Method Pointer")
	}

	l.recvmx.Lock()
	defer l.recvmx.Unlock()

	arr = make([]LinkStateType, len(l.recvmap))

	for _, add := range l.recvmap {
		arr = append(arr, add.State)
	}

	return arr
}
