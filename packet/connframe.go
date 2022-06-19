package packet

import (
	"encoding/gob"
	"sync"
	"time"

	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/octolib"
)

//
// Provides wrapper for packet.Eth frames
// Allows for tracking Eth frames and only allowing the first one recieved.
//

const frameMapSize uint = 16384
const packetTimeLimit time.Duration = time.Second * -5
const packetHighWater uint = 8192

type ConnFrameID uint64

type ConnFrame struct {
	ID    ConnFrameID
	Frame *EthFrame
}

type ConnFrameTrackerStruct struct {
	ctx      *ctx.Ctx
	m        sync.Mutex
	id       map[ConnFrameID]time.Time
	stamp    map[time.Time]ConnFrameID
	count    int
	runClean chan interface{}
}

var ConnFrameIDChan chan uint64

//
//
//
func init() {
	ConnFrameIDChan = octolib.RunGoCounter64()
	gob.Register(ConnFrame{})
}

func NewFrameTracker(cx *ctx.Ctx) (f *ConnFrameTrackerStruct) {
	f = &ConnFrameTrackerStruct{
		ctx:      cx,
		m:        sync.Mutex{},
		id:       make(map[ConnFrameID]time.Time, frameMapSize),
		stamp:    make(map[time.Time]ConnFrameID, frameMapSize),
		count:    0,
		runClean: make(chan interface{}, 1),
	}
	go f.goRunCleanUp()
	return f
}

func NewConnFrame(eth *EthFrame) (cf *ConnFrame) {

	cf = &ConnFrame{
		ID:    ConnFrameID(<-ConnFrameIDChan),
		Frame: eth,
	}

	return cf
}

//
// PassFrame()
// If this is the first encounter with the frame (ID) then pass it on
//
func (f *ConnFrameTrackerStruct) PassFrame(id ConnFrameID) (b bool) {

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
func (f *ConnFrameTrackerStruct) addID(id ConnFrameID) {

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
func (f *ConnFrameTrackerStruct) delID(id ConnFrameID) {
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
