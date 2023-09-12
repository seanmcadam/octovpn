package link

import (
	"sync"

	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type LinkModeType uint8
type LinkStateType uint8
type LinkNoticeType uint8
type LinkNoticeStateType uint16
type LinkNoticeStateListenCh chan<- LinkNoticeStateType
type LinkNoticeStateCh <-chan LinkNoticeStateType
type LinkNoticeStateFunc func() LinkNoticeStateCh

const DefaultListeners int = 2

const LinkModePassALL LinkModeType = 0x00
const LinkModePassState LinkModeType = 0x00
const LinkModePassNotice LinkModeType = 0x00
const LinkModeConnectedAND LinkModeType = 0x01 // All Links are connected
const LinkModeConnectedOR LinkModeType = 0x02  // One link is connected
const LinkModeUpAND LinkModeType = 0x01        // All Links are up (any up state)
const LinkModeUpOR LinkModeType = 0x02         // One link is up (any up state)
const LinkModeDownAND LinkModeType = 0x04      // All Links are down
const LinkModeDownOR LinkModeType = 0x08       // One link is down

const LinkModeFilterNotices LinkModeType = 0x08

const LinkStateNONE LinkStateType = 0x00
const LinkStateNOLINK LinkStateType = 0x01
const LinkStateLINK LinkStateType = 0x10
const LinkStateCHAL LinkStateType = 0x20
const LinkStateAUTH LinkStateType = 0x30
const LinkStateCONNECTED LinkStateType = 0x80
const LinkStateERROR LinkStateType = 0xFF
const LinkStateUpMASK LinkStateType = 0xF0
const LinkStateDownMASK LinkStateType = 0x0F

const LinkNoticeNONE LinkNoticeType = 0x00
const LinkNoticeLOSS LinkNoticeType = 0x01
const LinkNoticeLATENCY LinkNoticeType = 0x02
const LinkNoticeSATURATED LinkNoticeType = 0x04
const LinkNoticeCLOSED LinkNoticeType = 0x80
const LinkNoticeERROR LinkNoticeType = 0xFF

var recvnew chan LinkNoticeStateType

type AddLinkStruct struct {
	State LinkStateType
	LinkFunc LinkNoticeStateFunc
}

type LinkChan struct {
	name       string
	listenChan chan LinkNoticeStateListenCh
}

type LinkStateStruct struct {
	cx              *ctx.Ctx
	mx              sync.Mutex
	mode            LinkModeType
	instance        uint32
	state           LinkStateType
	linkNoticeState *LinkChan
	linkState       *LinkChan
	linkNotice      *LinkChan
	linkUpDown      *LinkChan
	linkUp          *LinkChan
	linkDown        *LinkChan
	linkConnected   *LinkChan
	linkNoLink      *LinkChan
	linkLink        *LinkChan
	linkAuth        *LinkChan
	linkChal        *LinkChan
	linkLoss        *LinkChan
	linkLatency     *LinkChan
	linkSaturation  *LinkChan
	linkClose       *LinkChan
	recvcounter     counter.CounterStruct
	recvfn          map[counter.Counter]LinkNoticeStateFunc
	recvchan        map[counter.Counter]LinkNoticeStateCh
	recvstate       map[counter.Counter]LinkStateType
	recvnew         chan LinkNoticeStateType
	addlinkch       chan *AddLinkStruct
	dellinkch       chan counter.Counter
}

var instanceCounter counter.CounterStruct

func init() {
	instanceCounter = counter.NewCounter32(ctx.NewContext())
}

func (ns LinkNoticeStateType) Notice() LinkNoticeType {
	return LinkNoticeType(ns >> 8)
}
func (ns LinkNoticeStateType) State() LinkStateType {
	return LinkStateType(ns & 0xFF)
}

func (state LinkStateType) String() string {
	switch state {
	case LinkStateNONE:
		return "NONE"
	case LinkStateNOLINK:
		return "NOLINK"
	case LinkStateLINK:
		return "LINK"
	case LinkStateCHAL:
		return "CHALLENGE"
	case LinkStateAUTH:
		return "AUTHENTICATED"
	case LinkStateCONNECTED:
		return "CONNECTED"
	case LinkStateERROR:
		return "ERROR"
	default:
		log.FatalfStack("unsupported LinkStateType:%02X", uint8(state))
	}
	return ""
}

func (state LinkNoticeType) String() string {
	switch state {
	case LinkNoticeNONE:
		return "NONE"
	case LinkNoticeLOSS:
		return "LOSS"
	case LinkNoticeLATENCY:
		return "LATENCY"
	case LinkNoticeSATURATED:
		return "SATURATED"
	case LinkNoticeCLOSED:
		return "CLOSED"
	case LinkNoticeERROR:
		return "ERROR"
	default:
		log.FatalfStack("unsupported LinkNoticeType:%02X", uint8(state))
	}
	return ""
}

func (ns LinkNoticeStateType) String() (s string) {
	s = ns.Notice().String() + "|" + ns.State().String()
	return s
}
