package tracker

import (
	"time"

	"github.com/seanmcadam/octovpn/octolib/counter"
	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/packet/packetchan"
)

const DefaultTrackerDataDepth = 4096
const DefaultTrackerChanDepth = 64
const MaintTimeoutDuration = 5 * time.Second

type TrackerStruct struct {
	ticker  *time.Ticker
	entry   map[counter.Counter64]time.Time
	track   map[counter.Counter64]interface{}
	latency map[time.Time]time.Duration
	acks    map[time.Time]int
	naks    map[time.Time]int
	send    map[time.Time]int
	recv    map[time.Time]int
	recvch  chan interface{}
	sendch  chan interface{}
	ackch   chan counter.Counter64
	nakch   chan counter.Counter64
	closech chan interface{}
}

func NewTracker(closech chan interface{}) (t *TrackerStruct) {

	t = &TrackerStruct{
		ticker:  time.NewTicker(1 * time.Second),
		entry:   make(map[counter.Counter64]time.Time, DefaultTrackerDataDepth),
		track:   make(map[counter.Counter64]interface{}, DefaultTrackerDataDepth),
		latency: make(map[time.Time]time.Duration, DefaultTrackerDataDepth),
		acks:    make(map[time.Time]int, DefaultTrackerDataDepth),
		naks:    make(map[time.Time]int, DefaultTrackerDataDepth),
		recv:    make(map[time.Time]int, DefaultTrackerDataDepth),
		send:    make(map[time.Time]int, DefaultTrackerDataDepth),
		recvch:  make(chan interface{}, DefaultTrackerChanDepth),
		sendch:  make(chan interface{}, DefaultTrackerChanDepth),
		ackch:   make(chan counter.Counter64, DefaultTrackerChanDepth),
		nakch:   make(chan counter.Counter64, DefaultTrackerChanDepth),
		closech: closech,
	}

	go t.goRun()

	return t
}

// For acting on the onject to be serialized with Send, Recv, Ack, and Nak
func (t *TrackerStruct) Send(data interface{}) {
	t.sendch <- data
}

func (t *TrackerStruct) Recv(data interface{}) {
	t.sendch <- data
}

func (t *TrackerStruct) Ack(counter counter.Counter64) {
	t.ackch <- counter
}

func (t *TrackerStruct) Nck(counter counter.Counter64) {
	t.nakch <- counter
}

// This is the serializer
func (t *TrackerStruct) goRun() {

	defer t.ticker.Stop()

	for {
		select {
		case <-t.ticker.C:
			t.maint()

		case p := <-t.sendch:
			t.snd(p)

		case p := <-t.recvch:
			t.rcv(p)

		case ack := <-t.ackch:
			t.ack(ack)

		case nak := <-t.nakch:
			t.nak(nak)
		case <-t.closech:
			return
		}
	}
}

// send()
// Take an interface (some sort of packet) and push the data into
// Tracker for latency, loss and bandwith calc
func (t *TrackerStruct) snd(p interface{}) {
	// Decode data type
	switch data := p.(type) {
	case *packetchan.ChanPacket:
		count := counter.Counter64(data.GetCounter())
		t.entry[count] = time.Now()
		t.track[count] = p

	default:
		log.Fatalf("Unhandled type %t", data)
	}

}

// rev()
// track recieved packets for bandwidth
func (t *TrackerStruct) rcv(p interface{}) {
	// Decode data type
	switch data := p.(type) {
	case *packetchan.ChanPacket:
		tm := time.Now()
		count := counter.Counter64(data.GetCounter())
		if _, ok := t.recv[tm]; !ok {
			t.recv[tm] = int(count)
		} else {
			t.recv[tm] += int(count)
		}

	default:
		log.Fatalf("Unhandled type %t", data)
	}
}

//
// ack
// Recieve ACK packet
// Add to send bandwidth, add and latecny
// Remove from entry and track
//
func (t *TrackerStruct) ack(c counter.Counter64) {

	if p, ok := t.track[c]; !ok {
		log.Warnf("Tracker ACK lost %d", c)
		return
	} else {
		switch data := p.(type) {
		case *packetchan.ChanPacket:
			tm := time.Now()
			count := counter.Counter64(data.GetCounter())
			if _, ok := t.send[tm]; !ok {
				t.send[tm] = int(count)
			} else {
				t.send[tm] += int(count)
			}

		default:
			log.Fatalf("Unhandled type %t", data)
		}

	}

	if entry, ok := t.entry[c]; !ok {
		log.Fatalf("Tracker ACK missing entry %d", c)
		return
	} else {
		t.latency[time.Now()] = time.Since(entry)
	}

	delete(t.track, c)
	delete(t.entry, c)

}

func (t *TrackerStruct) nak(c counter.Counter64) {
	if _, ok := t.track[c]; !ok {
		log.Warnf("Tracker NAK lost %d", c)
		return
	}

	if _, ok := t.entry[c]; !ok {
		log.Fatalf("Tracker NAK missing entry %d", c)
		return
	}

	delete(t.track, c)
	delete(t.entry, c)

}


//
// Tasks:
// Find Stale entries (no ACKs nor NAKs)
//
func (tracker *TrackerStruct) maint() {
	for c, t := range tracker.entry {
		if time.Since(t) > MaintTimeoutDuration {
			log.Warnf("Tracker Stale Entry %d", c)
			delete(tracker.track, c)
			delete(tracker.entry, c)
		}
	}
}
