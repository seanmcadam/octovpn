package connection

import (
	"encoding/gob"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/octolib"
	"github.com/seanmcadam/octovpn/packet"
)

type headerString string

const frameMapSize uint = 16384
const packetTimeLimit time.Duration = time.Second * -5
const packetHighWater uint = 8192
const UDPHeaderSignature headerString = "-OctoVPN-UDP-1"
const TCPHeaderSignature headerString = "-OctoVPN-TCP-1"

type connFrameID uint64

type ConnectionHeader struct {
	count uint64
	lenth uint16
	frame *ConnFrame
}

type ConnFrame struct {
	ID    connFrameID
	Frame *packet.EthFrame
}

type ConnFrameTrackerStruct struct {
	ctx      *ctx.Ctx
	m        sync.Mutex
	id       map[connFrameID]time.Time
	stamp    map[time.Time]connFrameID
	count    int
	runClean chan interface{}
}

type ProtoHeader struct {
	Signature headerString
	ID        uint64
	Payload   interface{}
}

var counterHeaderIDChan chan uint64

//
//
//
func init() {
	counterHeaderIDChan = octolib.RunGoCounter64()
	gob.Register(ConnFrame{})
	gob.Register(ProtoHeader{})
}

func NewProtoHeader(sig headerString, payload interface{}) (p *ProtoHeader, e error) {

	switch payload.(type) {
	case *ConnFrame:
	case ConnFrame:
	case *Ping:
	case Ping:
	case *Pong:
	case Pong:
	default:
		return nil, errors.New(fmt.Sprintf("Invalid payload type:%t", payload))
	}

	p = &ProtoHeader{
		Signature: sig,
		ID:        <-counterHeaderIDChan,
		Payload:   payload,
	}
	return p, e
}

func newFrameTracker(cx *ctx.Ctx) (f *ConnFrameTrackerStruct) {
	f = &ConnFrameTrackerStruct{
		ctx:      cx,
		m:        sync.Mutex{},
		id:       make(map[connFrameID]time.Time, frameMapSize),
		stamp:    make(map[time.Time]connFrameID, frameMapSize),
		count:    0,
		runClean: make(chan interface{}, 1),
	}
	go f.goRunCleanUp()
	return f
}

//
// PassFrame()
// If this is the first encounter with the frame (ID) then pass it on
//
func (f *ConnFrameTrackerStruct) PassFrame(id connFrameID) (b bool) {

	// Does the id already exist, meaning that the frame has been seen before
	_, ok := f.id[id]
	if ok {
		return false
	}

	f.addID(id)

	if f.count > int(packetHighWater) {
		select {
		case f.runClean <- 1:
		default:
		}
	}

	return true
}

//
//
//
func (f *ConnFrameTrackerStruct) goRunCleanUp() {

	for {
		select {
		case <-f.ctx.DoneChan():
			return
		case <-f.runClean:
			t := time.Now()
			t = t.Add(packetTimeLimit) // Add the negative time limit

			for i, j := range f.id {
				if j.Before(t) {
					f.delID(i)
				}
			}
		}
	}
}

//
// addID()
//
func (f *ConnFrameTrackerStruct) addID(id connFrameID) {

	t := time.Now()
	_, ok := f.stamp[t]
	if ok {
		f.ctx.Logf(ctx.LogLevelPanic, "Duplicate times stamp:%s, First ID:%d Second ID:%d", t, f.stamp[t], id)
		return
	}

	f.m.Lock()
	defer f.m.Unlock()

	f.id[id] = t
	f.stamp[t] = id
	f.count++
}

//
// delID()
//
func (f *ConnFrameTrackerStruct) delID(id connFrameID) {
	t, ok := f.id[id]
	if !ok {
		f.ctx.Logf(ctx.LogLevelError, " non-existant ID:%d", id)
	} else {
		i, ok := f.stamp[t]
		if ok {
			if i != id {
				f.ctx.Logf(ctx.LogLevelPanic, "IDs do not match timeID:%d for ID:%d", i, id)
			}
			f.m.Lock()
			defer f.m.Unlock()
			delete(f.stamp, t)
			delete(f.id, id)
			f.count--
		} else {
			f.ctx.Logf(ctx.LogLevelError, "Deleteing non-existant time:%s for ID:%d", t, id)
		}
	}
}
