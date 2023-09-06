package link

import (
	"sync"

	"github.com/seanmcadam/octovpn/octolib/log"
)

type LinkStateType uint8

const DefaultListeners int = 2

const LinkStateNone LinkStateType = 0x00
const LinkStateUp LinkStateType = 0x01
const LinkStateDown LinkStateType = 0x02
const LinkStateError LinkStateType = 0xF0

type LinkStateStruct struct {
	mx        sync.Mutex
	state     LinkStateType
	listeners []chan LinkStateType
}

func NewLinkState() (ls *LinkStateStruct) {
	ls = &LinkStateStruct{
		state:     LinkStateNone,
		listeners: []chan LinkStateType{},
	}
	return ls
}

func (ls *LinkStateStruct) ToggleState(s LinkStateType) {
	if s != ls.state && s != LinkStateNone {
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
