package link

import (
	"sync"

	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type LinkStateType uint8

const DefaultListeners int = 2

const LinkStateNone LinkStateType = 0x00
const LinkStateUp LinkStateType = 0x01
const LinkStateDown LinkStateType = 0x02
const LinkStateClose LinkStateType = 0x0F
const LinkStateError LinkStateType = 0xF0

type LinkStateStruct struct {
	mx        sync.Mutex
	cx        *ctx.Ctx
	instance  uint32
	state     LinkStateType
	listeners []chan LinkStateType
}

var instanceCounter counter.CounterStruct

func init() {
	instanceCounter = counter.NewCounter32(ctx.NewContext())
}

func NewLinkState(ctx *ctx.Ctx) (ls *LinkStateStruct) {
	count := <-instanceCounter.GetCountCh()
	ls = &LinkStateStruct{
		cx:        ctx,
		instance:  count.Uint().(uint32),
		state:     LinkStateNone,
		listeners: []chan LinkStateType{},
	}
	log.Debugf("Link[%d] Starting", ls.instance)
	go ls.goRun()
	return ls
}

func (ls *LinkStateStruct) ToggleState(s LinkStateType) {
	if s == LinkStateClose {
		log.FatalStack("ToggleState() got a Close... not allowed")
	}
	if s == LinkStateNone {
		log.Debug("ToggleState() got a None... skipping")
		return
	}

	ls.toggleState(s)
}

func (ls *LinkStateStruct) toggleState(s LinkStateType) {
	if s != ls.state && s != LinkStateNone {
		log.Debugf("Link[%d] Toggled State:%s", ls.instance, s)
		ls.mx.Lock()
		defer ls.mx.Unlock()
		ls.state = s
		for _, l := range ls.listeners {
			l <- ls.state
			close(l)
		}
		ls.listeners = ls.listeners[:0]
	} else {
		if s == LinkStateNone {
			log.FatalStack("ToggleState to None")
		}
	}
}

func (ls *LinkStateStruct) StateToggleCh() (newch chan LinkStateType) {
	ls.mx.Lock()
	defer ls.mx.Unlock()
	newch = make(chan LinkStateType, 1)
	ls.listeners = append(ls.listeners, newch)
	return newch
}

func (ls *LinkStateStruct) GetState() LinkStateType {
	return ls.state
}

func (ls *LinkStateStruct) goRun() {
	for {
		select {
		case <-ls.cx.DoneChan():
			log.Debugf("Link[%d] Shutdown", ls.instance)
			ls.toggleState(LinkStateClose)
		}
	}
}

func (state LinkStateType) String() string {
	switch state {
	case LinkStateNone:
		return "NONE"
	case LinkStateUp:
		return "UP"
	case LinkStateDown:
		return "DOWN"
	case LinkStateClose:
		return "CLOSE"
	case LinkStateError:
		return "ERROR"
	default:
		log.FatalfStack("unsupported LinkStateType:%02X", uint8(state))
	}
	return ""
}
