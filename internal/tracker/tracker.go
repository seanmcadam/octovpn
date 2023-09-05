package tracker

import (
	"time"

	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

//
//
//
//
//

const DefaultTrackerDataDepth = 4096
const DefaultTrackerChanDepth = 64
const MaintTimeoutDuration = 5 * time.Second

type PacketTracker struct {
	packet *packet.PacketStruct
	tm     time.Time
}

type CounterTracker struct {
	counter counter.Counter
	tm      time.Time
}

type DataTracker struct {
	interval    time.Duration
	sendbytes   uint64
	recvbytes   uint64
	sendpackets uint64
	recvpackets uint64
	ackcount    uint64
	naccount    uint64
}

type TrackerStruct struct {
	cx     *ctx.Ctx
	freq   time.Duration
	ticker *time.Ticker
	recvch chan *PacketTracker
	sendch chan *PacketTracker
	ackch  chan *CounterTracker
	nakch  chan *CounterTracker
	recv   map[counter.Counter]*PacketTracker
	sent   map[counter.Counter]*PacketTracker
	acknak map[counter.Counter]*PacketTracker
	ack    map[counter.Counter]*CounterTracker
	nak    map[counter.Counter]*CounterTracker
}

func newPacketTracker(d *packet.PacketStruct) (p *PacketTracker) {
	p = &PacketTracker{
		packet: d,
		tm:     time.Now(),
	}
	return p
}

func newCounterTracker(c counter.Counter) (p *CounterTracker) {
	p = &CounterTracker{
		counter: c,
		tm:      time.Now(),
	}
	return p
}

func NewTracker(ctx *ctx.Ctx, freq time.Duration) (t *TrackerStruct) {

	t = &TrackerStruct{
		cx:     ctx,
		freq:   freq,
		ticker: time.NewTicker(freq),
		recvch: make(chan *PacketTracker, DefaultTrackerChanDepth),
		sendch: make(chan *PacketTracker, DefaultTrackerChanDepth),
		ackch:  make(chan *CounterTracker, DefaultTrackerChanDepth),
		nakch:  make(chan *CounterTracker, DefaultTrackerChanDepth),
		recv:   make(map[counter.Counter]*PacketTracker, DefaultTrackerDataDepth),
		sent:   make(map[counter.Counter]*PacketTracker, DefaultTrackerDataDepth),
		acknak: make(map[counter.Counter]*PacketTracker, DefaultTrackerDataDepth),
		ack:    make(map[counter.Counter]*CounterTracker, DefaultTrackerDataDepth),
		nak:    make(map[counter.Counter]*CounterTracker, DefaultTrackerDataDepth),
	}

	go t.goRun()

	return t
}

// For acting on the onject to be serialized with Send, Recv, Ack, and Nak
func (t *TrackerStruct) Send(p *packet.PacketStruct) {
	t.sendch <- newPacketTracker(p)
}

func (t *TrackerStruct) Recv(p *packet.PacketStruct) {
	t.sendch <- newPacketTracker(p)
}

func (t *TrackerStruct) Ack(counter counter.Counter) {
	t.ackch <- newCounterTracker(counter)
}

func (t *TrackerStruct) Nak(counter counter.Counter) {
	t.nakch <- newCounterTracker(counter)
}

// This is the serializer
func (t *TrackerStruct) goRun() {

	defer t.ticker.Stop()

	for {
		select {
		case <-t.ticker.C:
			t.maint()

		case p := <-t.sendch:
			t.sendHandler(p)

		case p := <-t.recvch:
			t.recvHandler(p)

		case ack := <-t.ackch:
			t.ackHandler(ack)

		case nak := <-t.nakch:
			t.nakHandler(nak)

		case <-t.cx.DoneChan():
			return
		}
	}
}

// send()
// Take an interface (some sort of packet) and push the data into
// Tracker for latency, loss and bandwith calc
func (t *TrackerStruct) sendHandler(pt *PacketTracker) {
	count := counter.Counter(pt.packet.Counter())
	t.sent[count] = pt
	t.acknak[count] = pt
}

// rev()
// track recieved packets for bandwidth
func (t *TrackerStruct) recvHandler(pt *PacketTracker) {
	count := counter.Counter(pt.packet.Counter())
	t.recv[count] = pt
}

// ack
// Recieve ACK packet
func (t *TrackerStruct) ackHandler(ct *CounterTracker) {
	t.ack[ct.counter] = ct
	delete(t.acknak, ct.counter)
}

func (t *TrackerStruct) nakHandler(ct *CounterTracker) {
	t.nak[ct.counter] = ct
	delete(t.acknak, ct.counter)
}

// Tasks:
// Find Stale entries (no ACKs nor NAKs)
func (tracker *TrackerStruct) maint() {

	//log.Debug("Tracker maint running")

	dt := &DataTracker{
		interval:    tracker.freq,
		sendbytes:   0,
		recvbytes:   0,
		sendpackets: 0,
		recvpackets: 0,
		ackcount:    0,
		naccount:    0,
	}

	for count, pt := range tracker.sent {
		dt.sendbytes += uint64(pt.packet.Size())
		dt.sendpackets++
		delete(tracker.sent, count)
	}

	for count, pt := range tracker.recv {
		dt.recvbytes += uint64(pt.packet.Size())
		dt.recvpackets++
		delete(tracker.recv, count)
	}

	for count, ct := range tracker.ack {
		dt.ackcount++
		_ = ct
		_ = count
	}

	for count, ct := range tracker.nak {
		dt.ackcount++
		_ = ct
		_ = count
	}

	// Calc Send BW
	// Calc Send Count
	// Calc Recv BW
	// Calc Recv Count
	// Calc Ack Count
	// Calc Nak Count

	log.Debugf("Data:%v", dt)

}
