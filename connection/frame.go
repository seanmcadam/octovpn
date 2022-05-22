package connection

import (
	"sync"
	"time"

	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/iface"
)

const packetTimeLimit time.Duration = time.Second * 5

type connFrameID uint64

type ConnFrame struct {
	id    connFrameID
	frame *iface.Frame
}

type ConnFrameTrackerStruct struct {
	ctx        *ctx.Ctx
	m          sync.Mutex
	id         map[connFrameID]time.Time
	stamp      map[time.Time]connFrameID
	count      int
	deleteChan chan connFrameID
}

func newFrameTracker(cx *ctx.Ctx) (f *ConnFrameTrackerStruct) {
	f = &ConnFrameTrackerStruct{
		ctx:        cx,
		m:          sync.Mutex{},
		id:         make(map[connFrameID]time.Time),
		stamp:      make(map[time.Time]connFrameID),
		count:      0,
		deleteChan: make(chan connFrameID),
	}
	go f.goDeleteID()
}

func (f *ConnFrameTrackerStruct) PassFrame(id connFrameID) (b bool) {
	f.m.Lock()
	defer f.m.Unlock()

	// Does the id already exist, meaning that the frame has been seen before
	_, ok := f.id[id]
	if ok {
		return false
	}

	t := time.Now()
	_, ok := f.stamp[t]
	if ok {
		f.ctx.Logf(ctx.LogLevelPanic, "Duplicate times stamp:%s, First ID:%d Second ID:%d", t, f.stamp[t], id)
	}

	f.id[id] = t
	f.stamp[t] = id
	f.count++

}

func (f *ConnFrameTrackerStruct) CleanUp() {
	f.m.Lock()
	defer f.m.Unlock()

	// Feed IDs here to the cleanup channel

}

func (f *ConnFrameTrackerStruct) goDeleteID() {
	for {
		select {
		case <-f.ctx.Done():
			return

		case id := <-f.deleteChan:
			f.m.Lock()
			t, ok := f.id[id]
			if !ok {
				f.ctx.Logf(ctx.LogLevelError, "Deleteing non-existant ID:%d", id)
				return
			}
			_, ok = f.stamp[t]
			if !ok {
				f.ctx.Logf(ctx.LogLevelError, "Deleteing non-existant time:%s for ID:%d", t, id)
			} else {
				delete(f.stamp, t)
			}

			delete(f.id, id)
			f.count--
			f.m.Unlock()

		}
	}
}
