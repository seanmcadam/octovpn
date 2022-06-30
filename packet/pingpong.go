package packet

import (
	"encoding/binary"
	"sync"
	"time"

	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/octolib"
)

const defaultPingsPerSec uint16 = 5
const maximumPingsPerSec uint16 = 15
const minimumPingsPerSec uint16 = 1
const defaultPingPeriodMS uint16 = 1000 / defaultPingsPerSec
const minimumPingPeriodMS uint16 = 1000 / maximumPingsPerSec
const maximumPingPeriodMS uint16 = 1000 / minimumPingsPerSec

const pingMapSize uint = 32
const pingTimeLimit time.Duration = time.Second * -1
const pingHighWater uint = 1

type PingID uint64
type lossPercent uint8
type Deviation uint16
type Latency uint16
type Loss uint16
type Status uint16

const (
	Status1S   Status = 1
	Status5S   Status = 5
	Status15S  Status = 15
	Status60S  Status = 60
	Status300S Status = 300
)

type Mutex struct {
	mutex sync.Mutex
}

type PingerStruct struct {
	ctx        *ctx.Ctx
	mutex      *Mutex
	sendChan   chan interface{}
	runCalc    chan interface{}
	periodMS   uint16
	sent       map[PingID]int64
	rtt        map[PingID]int64
	currentID  PingID
	count      uint
	RTTAveS1   time.Duration
	RTTAveS5   time.Duration
	RTTAveS15  time.Duration
	RTTAveS60  time.Duration
	RTTAveS300 time.Duration
	LossS1     lossPercent
	LossS5     lossPercent
	LossS15    lossPercent
	LossS60    lossPercent
	LossS300   lossPercent
}

type Ping struct {
	id   PingID
	Send int64
}

func (p *Ping) PacketID() PacketID {
	return PacketID(p.id)
}

func (p *Ping) Type() PacketType {
	return PingPacket
}

func (p *Ping) Length() PacketPayloadLen {
	return PacketPayloadLen(8 + 8)
}

func (ch *Ping) ToByte() (buf []byte) {
	totallen := 8 + 8
	buf = make([]byte, 0, totallen)
	buf4a := make([]byte, 4)
	binary.BigEndian.PutUint32(buf4a, uint32(ch.id))

	buf4b := make([]byte, 4)
	binary.BigEndian.PutUint32(buf4b, uint32(ch.Send))

	buf = append(buf, buf4a...)
	buf = append(buf, buf4b...)

	return buf
}

type Pong struct {
	id   PingID
	Send int64
	Recv int64
}

func (p *Pong) PacketID() PacketID {
	return PacketID(p.id)
}

func (p *Pong) Type() PacketType {
	return PongPacket
}

func (p *Pong) Length() PacketPayloadLen {
	return PacketPayloadLen(8 + 8 + 8)
}

func (ch *Pong) ToByte() (buf []byte) {
	totallen := 8 + 8 + 8
	buf = make([]byte, 0, totallen)
	buf4a := make([]byte, 4)
	binary.BigEndian.PutUint32(buf4a, uint32(ch.id))

	buf4b := make([]byte, 4)
	binary.BigEndian.PutUint32(buf4b, uint32(ch.Send))

	buf4c := make([]byte, 4)
	binary.BigEndian.PutUint32(buf4c, uint32(ch.Recv))

	buf = append(buf, buf4a...)
	buf = append(buf, buf4b...)
	buf = append(buf, buf4c...)

	return buf
}

var counterPingChan chan uint64

//
//
//
func init() {
	counterPingChan = octolib.RunGoCounter64()
}

//
//
//
func NewPinger(cx *ctx.Ctx, pingsPerSec uint16) (pinger *PingerStruct) {

	cx.LogLocation()

	var period uint16

	//
	// Figure out what the period is
	// The period is the number of MS between each
	//
	if pingsPerSec == 0 {
		period = defaultPingPeriodMS
	} else if pingsPerSec > maximumPingsPerSec {
		period = minimumPingPeriodMS
	} else {
		period = maximumPingPeriodMS
	}

	//
	// (300 seconds/1000 ms/s)/period
	// Calculate the number of periods are in 300 seconds
	//
	count := uint((300*1000)/uint(period)) + 1

	pinger = &PingerStruct{
		ctx:        cx,
		mutex:      &Mutex{},
		sendChan:   make(chan interface{}, 3),
		periodMS:   period,
		sent:       make(map[PingID]int64, count),
		rtt:        make(map[PingID]int64, count),
		currentID:  0,
		count:      count,
		RTTAveS1:   0,
		RTTAveS5:   0,
		RTTAveS15:  0,
		RTTAveS60:  0,
		RTTAveS300: 0,
		LossS1:     0,
		LossS5:     0,
		LossS15:    0,
		LossS60:    0,
		LossS300:   0,
		runCalc:    make(chan interface{}),
	}

	return pinger
}

//
//
//
func (m *Mutex) Lock() {
	m.mutex.Lock()
}

//
//
//
func (m *Mutex) Unlock() {
	m.mutex.Unlock()
}

//
//
//
func (p *PingerStruct) Start() {
	go p.goPing()
	go p.goCalc()
}

//
//
//
func (p *PingerStruct) Stop() {
	p.ctx.Cancel()
}

//
//
//
func (p *PingerStruct) SendChan() chan interface{} {
	return p.sendChan
}

//
//
//
func (p *PingerStruct) ping() {

	var lastID uint64 = 0
	pingID := <-counterPingChan

	now := time.Now()
	send := now.UnixNano()
	ping := &Ping{
		id:   PingID(pingID),
		Send: send,
	}

	if pingID > uint64(p.count) {
		lastID = pingID - uint64(p.count)
	}

	select {
	case p.sendChan <- ping:
	default:
		p.ctx.Logf(ctx.LogLevelError, "Send Ping ID:%d failed", pingID)
		return
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.sent[ping.id] = ping.Send
	p.currentID = ping.id
	if lastID > 0 {
		_, ok := p.sent[PingID(lastID)]
		if ok {
			delete(p.sent, PingID(lastID))
		}
		_, ok = p.rtt[PingID(lastID)]
		if ok {
			delete(p.rtt, PingID(lastID))
		}
	}
}

//
//
//
func (p *PingerStruct) goCalc() {
	const count5time = 2
	const count15time = 5
	const count60time = 10
	const count300time = 20
	var count5 = count5time
	var count15 = count15time
	var count60 = count60time
	var count300 = count300time

	for {
		select {
		case <-p.ctx.DoneChan():
			return
		case <-time.After(time.Second):
			count5--
			count15--
			count60--
			count300--
			p.calc(Status1S)

			if count5 == 0 {
				count5 = count5time
				p.calc(Status5S)
			}
			if count15 == 0 {
				count15 = count15time
				p.calc(Status15S)
			}
			if count60 == 0 {
				count60 = count60time
				p.calc(Status60S)
			}
			if count300 == 0 {
				count300 = count300time
				p.calc(Status300S)
			}
		}
	}
}

//
//
//
func (p *PingerStruct) goPing() {

	p.ctx.LogLocation()

	for {
		select {
		case <-p.ctx.DoneChan():
			return
		case <-time.After(time.Millisecond * time.Duration(p.periodMS)):
			p.ping()
		}
	}
}

//
//
//
func (p *PingerStruct) GotPing(ping *Ping) {
	p.ctx.LogLocation()
	t := time.Now()
	recv := t.UnixNano()
	pong := &Pong{
		id:   ping.id,
		Send: ping.Send,
		Recv: recv,
	}
	p.sendChan <- pong
}

//
//
//
func (p *PingerStruct) GotPong(pong *Pong) {

	p.ctx.LogLocation()
	p.mutex.Lock()
	defer p.mutex.Unlock()

	nanosent, ok := p.sent[pong.id]
	if !ok {
		p.ctx.Logf(ctx.LogLevelError, "Pong with no ping ID:%d", pong.id)
		return
	}

	_, ok = p.rtt[pong.id]
	if ok {
		p.ctx.Logf(ctx.LogLevelError, "Pong RTT already exists, ID:%d", pong.id)
		return
	}

	if nanosent != pong.Send {
		p.ctx.Logf(ctx.LogLevelError, "Ping and Pong times do not match:%d, %d", nanosent, pong.Send)
		return
	}

	nanonow := time.Now().UnixNano()
	delta := nanosent - nanonow

	p.rtt[pong.id] = delta

}

//
// calc()
// update the statistic used to determine the health of the link
//
func (p *PingerStruct) calc(status Status) {

	var count uint64 = uint64(status) * 1000 / uint64(p.periodMS)

	// This is for start up, dont calc until the ID (count) is high enough
	if count <= (uint64(p.currentID) - 1) {
		return
	}

	// Collect the last count data points and avarage them
	sum := uint64(0)
	items := uint64(0)
	p.mutex.Lock()
	for i := uint64(0); i < count; i++ {
		data, ok := p.rtt[PingID(i)]
		if ok {
			sum += uint64(data)
			items++
		}
	}
	p.mutex.Unlock()

	var rtt time.Duration
	var loss lossPercent

	if items > 0 {
		rtt = time.Duration(sum / uint64(items))
		loss = lossPercent(100 - (100 * items / count))
	} else {
		rtt = 0
		loss = lossPercent(100)
	}

	switch status {
	case Status1S:
		p.RTTAveS1 = rtt
		p.LossS1 = loss
	case Status5S:
		p.RTTAveS5 = rtt
		p.LossS5 = loss
	case Status15S:
		p.RTTAveS15 = rtt
		p.LossS15 = loss
	case Status60S:
		p.RTTAveS60 = rtt
		p.LossS60 = loss
	case Status300S:
		p.RTTAveS300 = rtt
		p.LossS300 = loss
	default:
		p.ctx.Logf(ctx.LogLevelPanic, "Bast Status:%s", status)
	}
}
