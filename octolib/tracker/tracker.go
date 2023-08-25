package tracker

import (
	"time"

	"github.com/seanmcadam/octovpn/octolib/counter"
	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/packet/packetchan"
)

const DefaultTrackerDepth = 64
const MaintTimeoutDuration = 5 * time.Second

type TrackerStruct struct {
	ticker *time.Ticker
	entry  map[counter.Counter64]time.Time
	track  map[counter.Counter64]interface{}
	pushch chan interface{}
	ackch  chan counter.Counter64
	nakch  chan counter.Counter64
	closech chan interface{}
}

func NewTracker(closech chan interface{}) (t *TrackerStruct) {

	t = &TrackerStruct{
		ticker: time.NewTicker(1 * time.Second),
		entry:  make(map[counter.Counter64]time.Time),
		track:  make(map[counter.Counter64]interface{}, 4096),
		pushch: make(chan interface{}, DefaultTrackerDepth),
		ackch:  make(chan counter.Counter64, DefaultTrackerDepth),
		nakch:  make(chan counter.Counter64, DefaultTrackerDepth),
		closech: closech,
	}

	go t.goRun()

	return t
}

func (t *TrackerStruct) Push(data interface{}) {
	t.pushch <- data
}

func (t *TrackerStruct) Ack(counter counter.Counter64) {
	t.ackch <- counter
}

func (t *TrackerStruct) Nck(counter counter.Counter64) {
	t.nakch <- counter
}

func (t *TrackerStruct) goRun() {

	defer t.ticker.Stop()

	for {
		select {
		case <-t.ticker.C:
			t.maint()

		case push := <-t.pushch:
			// Decode data type
			switch data := push.(type) {
			case packetchan.ChanPacket:
				count := counter.Counter64(data.GetCounter())
				t.entry[count] = time.Now()
				t.track[count] = push

			default:
				log.Fatalf("Unhandled type %t", data)
			}

		case ack := <-t.ackch:
			if _, ok := t.track[ack]; !ok {
				log.Warnf("Tracker ACK lost %d", ack)
				continue
			}

			delete(t.track, ack)
			delete(t.entry, ack)

		case nak := <-t.nakch:
			if _, ok := t.track[nak]; !ok {
				log.Warnf("Tracker NAK lost %d", nak)
				continue
			}
			delete(t.track, nak)
			delete(t.entry, nak)

		case <-t.closech:
			return
		}
	}
}

func (tracker *TrackerStruct) maint() {
	for c, t := range tracker.entry {
		if time.Since(t) > MaintTimeoutDuration {
			log.Warnf("Tracker Stale Entry %d", c)
			delete(tracker.track, c)
			delete(tracker.entry, c)
		}
	}
}
