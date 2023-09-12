package link

import (
	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func NewLinkState(ctx *ctx.Ctx, m ...LinkModeType) (ls *LinkStateStruct) {
	count := <-instanceCounter.GetCountCh()
	mode := LinkModePassALL // Default Pass All

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
		linkUpDown:      newLinkChan("UpDown"),
		linkUp:          newLinkChan("Up"),
		linkDown:        newLinkChan("Down"),
		linkLink:        newLinkChan("Link"),
		linkAuth:        newLinkChan("Auth"),
		linkChal:        newLinkChan("Chal"),
		linkNoLink:      newLinkChan("NoLink"),
		linkConnected:   newLinkChan("Connected"),
		linkLoss:        newLinkChan("Loss"),
		linkLatency:     newLinkChan("Latency"),
		linkSaturation:  newLinkChan("Saturation"),
		linkClose:       newLinkChan("Close"),
		recvcounter:     counter.NewCounter32(ctx),
		recvnew:         make(chan LinkNoticeStateType, 10),
		recvfn:          make(map[counter.Counter]LinkNoticeStateFunc),
		recvchan:        make(map[counter.Counter]LinkNoticeStateCh),
		recvstate:       make(map[counter.Counter]LinkStateType),
		addlinkch:       make(chan *AddLinkStruct, 5),
		dellinkch:       make(chan counter.Counter, 5),
	}

	ls.recvchan[counter.MakeCounter32(0)] = recvnew

	log.Debugf("Link[%d] Starting", ls.instance)
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
	if ls == nil{
		return
	}
	ls.setState(LinkStateNOLINK)
}

func (ls *LinkStateStruct) Link() {
	if ls == nil{
		return
	}
	ls.setState(LinkStateLINK)
}

func (ls *LinkStateStruct) Chal() {
	if ls == nil{
		return
	}
	ls.setState(LinkStateCHAL)
}

func (ls *LinkStateStruct) Auth() {
	if ls == nil{
		return
	}
	ls.setState(LinkStateAUTH)
}

func (ls *LinkStateStruct) Connected() {
	if ls == nil{
		return
	}
	ls.setState(LinkStateCONNECTED)
}

func (ls *LinkStateStruct) Close() {
	if ls == nil{
		return
	}
	ls.SendNotice(LinkNoticeCLOSED)
}

func (ls *LinkStateStruct) setState(s LinkStateType) {
	if ls == nil{
		return
	}
	if s == ls.state {
		return
	}
	log.GDebugf("Link State Change:%s -> %s", ls.state, s)
	ls.processMessage(noticeState(LinkNoticeNONE, s))
}

func (ls *LinkStateStruct) GetState() LinkStateType {
	if ls == nil{
		return 0
	}
	return ls.state
}

func (ls *LinkStateStruct) goRun() {
	if ls == nil{
		return 
	}

	defer log.Debugf("Link[%d] Shutdown", ls.instance)
	for {
		select {
		case <-ls.cx.DoneChan():
			ls.setState(LinkStateNOLINK)
			ls.processMessage(noticeState(LinkNoticeCLOSED, LinkStateNONE))
			return
		}
	}
}

func (ls *LinkStateStruct) goRecv() {
	if ls == nil{
		return 
	}
	defer ls.Cancel()

	for {
		for i, ch := range ls.recvchan {
			var index uint64

			switch i.Uint().(type) {
			case uint32:
				index = uint64(i.Uint().(uint32))
			case uint64:
				index = i.Uint().(uint64)
			}

			if index != 0 {
				log.GDebugf("Recv Link[%d] Msg", index)
				var ns LinkNoticeStateType
				select {
				case ns = <-ch:
					log.GDebugf("Recv Link[%d] Msg:%s", index, ns)
					ls.processMessage(ns)
					// Reload the channel
					ls.recvchan[i] = ls.recvfn[i]()
				default:
					// Channel closed, it is dead to me now.
					log.GDebugf("Got DEAD Link Delete:%d", index)
					ls.dellinkch <- i
				}
			} else {
				log.GDebug("Got Link Refresh")

				for {
					select {
					case add := <-ls.addlinkch:
						if add == nil {
							log.FatalStack("nil pointer")
						}
						c := ls.recvcounter.Next()
						ls.recvfn[c] = add.LinkFunc
						ls.recvchan[c] = add.LinkFunc()
						ls.recvstate[c] = add.State
						if add.State != LinkStateNONE {
							ls.processStateChange(noticeState(LinkNoticeNONE, add.State))
						}
					case c := <-ls.dellinkch:
						delete(ls.recvfn, c)
						delete(ls.recvchan, c)
						delete(ls.recvstate, c)
					default:
						break
					}
				}
			}
		}
	}
}

func noticeState(n LinkNoticeType, s LinkStateType) (ns LinkNoticeStateType) {
	return LinkNoticeStateType((uint16(n) << 8) | uint16(s))
}
