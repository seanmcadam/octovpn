package link

import (
	"fmt"

	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func newLinkState(ctx *ctx.Ctx, name string, m ...LinkModeType) (ls *LinkStateStruct) {
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

	ls = &LinkStateStruct{
		cx:              ctx,
		linkname:        fmt.Sprintf("link#%d ", count.Uint().(uint32)),
		mode:            mode,
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
		recvcounter:     counter.NewCounter32(ctx),
		recvnew:         make(chan LinkNoticeStateType, 5),
		recvfn:          make(map[counter.Counter]LinkNoticeStateFunc),
		recvchan:        make(map[counter.Counter]LinkNoticeStateCh),
		recvstate:       make(map[counter.Counter]LinkStateType),
		addlinkch:       make(chan *AddLinkStruct, 5),
		dellinkch:       make(chan counter.Counter, 5),
	}

	ls.recvchan[counter.MakeCounter32(0)] = recvnew

	// log.Debugf("[%s] Starting", ls.linkname)
	go ls.goRun()
	go ls.goRecv()
	return ls
}

//func (ls *LinkStateStruct) AddLink(fn LinkNoticeStateFunc) {
//	ls.addlinkch <- fn
//	ls.recvnew <- noticeState(LinkNoticeNONE, LinkStateNONE)
//}

func (ls *LinkStateStruct) SendNotice(n LinkNoticeType) {
	if n == LinkNoticeNONE {
		return
	}
	ls.processMessage(noticeState(n, LinkStateNONE))
}

func (ls *LinkStateStruct) NoLink() {
	if ls == nil {
		return
	}
	ls.setState(LinkStateNOLINK)
}

func (ls *LinkStateStruct) Listen() {
	if ls == nil {
		return
	}
	ls.setState(LinkStateLISTEN)
}

func (ls *LinkStateStruct) Link() {
	if ls == nil {
		return
	}
	ls.setState(LinkStateLINK)
}

func (ls *LinkStateStruct) Chal() {
	if ls == nil {
		return
	}
	ls.setState(LinkStateCHAL)
}

func (ls *LinkStateStruct) Auth() {
	if ls == nil {
		return
	}
	ls.setState(LinkStateAUTH)
}

func (ls *LinkStateStruct) Connected() {
	if ls == nil {
		return
	}
	ls.setState(LinkStateCONNECTED)
}

func (ls *LinkStateStruct) Close() {
	if ls == nil {
		return
	}
	ls.SendNotice(LinkNoticeCLOSED)
}

func (ls *LinkStateStruct) setState(s LinkStateType) {
	if ls == nil {
		return
	}
	if s == ls.state {
		return
	}
	log.GDebugf("Link State Change:%s -> %s", ls.state, s)
	ls.processMessage(noticeState(LinkNoticeNONE, s))
}

func (ls *LinkStateStruct) IsUp() bool {
	s := ls.state & LinkStateUpMASK 
	return s > 0
	//return ls.state & LinkStateUpMASK > 0
}

func (ls *LinkStateStruct) IsDown() bool {
	return ls.state & LinkStateDownMASK > 0
}

func (ls *LinkStateStruct) GetState() LinkStateType {
	if ls == nil {
		return 0
	}
	return ls.state
}

func (ls *LinkStateStruct) goRun() {
	if ls == nil {
		return
	}

	defer log.Debugf("[%s] Shutdown", ls.linkname)
	for {
		select {
		case <-ls.cx.DoneChan():
			ls.setState(LinkStateNOLINK)
			ls.processMessage(noticeState(LinkNoticeCLOSED, LinkStateNONE))
			return
		}
	}
}

func noticeState(n LinkNoticeType, s LinkStateType) (ns LinkNoticeStateType) {
	return LinkNoticeStateType((uint16(n) << 8) | uint16(s))
}
