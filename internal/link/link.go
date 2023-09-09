package link

import (
	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func NewLinkState(ctx *ctx.Ctx, m ...LinkModeType) (ls *LinkStateStruct) {
	count := <-instanceCounter.GetCountCh()
	mode := LinkModePassALL // Default PAss All

	if len(m) > 0 {
		mode = m[0]
	}

	ls = &LinkStateStruct{
		cx:              ctx,
		mode:            mode,
		instance:        count.Uint().(uint32),
		state:           LinkStateNONE,
		linkNoticeState: newLinkChan("NoticeState"),
		linkState:       newLinkChan("State"),
		linkNotice:      newLinkChan("Notice"),
		linkUp:          newLinkChan("Up"),
		linkLink:        newLinkChan("Link"),
		linkDown:        newLinkChan("Down"),
		linkLoss:        newLinkChan("Loss"),
		linkLatency:     newLinkChan("Latency"),
		linkSaturation:  newLinkChan("Saturation"),
		linkClose:       newLinkChan("Close"),
		recvcounter:     counter.NewCounter32(ctx),
		recvnew:         make(chan LinkNoticeStateType, 10),
		recvfn:          make(map[counter.Counter]LinkNoticeStateFunc),
		recvchan:        make(map[counter.Counter]LinkNoticeStateCh),
		recvstate:       make(map[counter.Counter]LinkStateType),
		addlinkch:       make(chan LinkNoticeStateFunc, 5),
		dellinkch:       make(chan counter.Counter, 5),
	}

	ls.recvchan[counter.MakeCounter32(0)] = recvnew

	log.Debugf("Link[%d] Starting", ls.instance)
	go ls.goRun()
	go ls.goRecv()
	return ls
}

func (ls *LinkStateStruct) AddLink(fn LinkNoticeStateFunc) {
	ls.addlinkch <- fn
	ls.recvnew <- noticeState(LinkNoticeNONE, LinkStateNONE)
}

func (ls *LinkStateStruct) SendNotice(n LinkNoticeType) {
	if n == LinkNoticeNONE {
		return
	}
	ls.processMessage(noticeState(n, LinkStateNONE))
}

func (ls *LinkStateStruct) Up() {
	ls.setState(LinkStateUP)
}

func (ls *LinkStateStruct) Link() {
	ls.setState(LinkStateLINK)
}

func (ls *LinkStateStruct) Down() {
	ls.setState(LinkStateDOWN)
}

func (ls *LinkStateStruct) Chal() {
	ls.setState(LinkStateCHAL)
}

func (ls *LinkStateStruct) Auth() {
	ls.setState(LinkStateAUTH)
}

func (ls *LinkStateStruct) setState(s LinkStateType) {
	if s == ls.state {
		return
	}
	log.Debugf("Link State Change:%s -> %s", ls.state, s)
	ls.state = s
	ls.processMessage(noticeState(LinkNoticeNONE, s))
}

func (ls *LinkStateStruct) GetState() LinkStateType {
	return ls.state
}

func (ls *LinkStateStruct) goRun() {
	defer log.Debugf("Link[%d] Shutdown", ls.instance)
	for {
		select {
		case <-ls.cx.DoneChan():
			ls.setState(LinkStateDOWN)
			ls.processMessage(noticeState(LinkNoticeCLOSED, LinkStateNONE))
			return
		}
	}
}

func (ls *LinkStateStruct) goRecv() {
FORLOOP:
	for {
		select {
		case fn := <-ls.addlinkch:
			if fn == nil {
				log.Debug("nil pointer")
			}
			c := ls.recvcounter.Next()
			ls.recvfn[c] = fn
			ls.recvchan[c] = fn()
			ls.recvstate[c] = LinkStateNONE
		case c := <-ls.dellinkch:
			delete(ls.recvfn, c)
			delete(ls.recvchan, c)
			delete(ls.recvstate, c)
		default:
		}

		for i, ch := range ls.recvchan {
			var ns LinkNoticeStateType
			select {
			case ns = <-ch:
				if i.Uint().(uint32) == 0 {
					continue FORLOOP
				}
				ls.processMessage(ns)
			default:
				// Channel closed, it is dead to me now.
				ls.dellinkch <- i
			}

			continue FORLOOP
		}
	}
}

func noticeState(n LinkNoticeType, s LinkStateType) (ns LinkNoticeStateType) {
	return LinkNoticeStateType((uint16(n) << 8) | uint16(s))
}
