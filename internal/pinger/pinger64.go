package pinger

import (
	"fmt"
	"time"

	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type Ping64 counter.Counter
type Pong64 counter.Counter

type Pinger64Struct struct {
	cx      *ctx.Ctx
	active  bool
	freq    time.Duration
	timeout time.Duration
	counter counter.CounterStruct
	pingch  chan Ping
	pongch  chan Pong
	Errorch chan error
}

func NewBytePing64(b []byte) (p Ping) {
	if len(b) != 4 {
		log.FatalfStack("Count data len:%d, :%0x", len(b), b)
	}

	c := counter.NewByteCounter32(b)
	p = counter.Counter(c)
	return p

}
func NewBytePong64(b []byte) (p Ping) {
	if len(b) != 8 {
		log.FatalfStack("Count data len:%d, :%0x", len(b), b)
	}

	c := counter.NewByteCounter64(b)
	p = counter.Counter(c)
	return p

}

func NewPinger64(ctx *ctx.Ctx, freq time.Duration, timeout time.Duration) (p PingerStruct) {
	p64 := &Pinger64Struct{
		cx:      ctx,
		active:  false,
		freq:    freq,
		timeout: timeout,
		counter: counter.NewCounter64(ctx),
		pingch:  make(chan Ping),
		pongch:  make(chan Pong),
		Errorch: make(chan error),
	}

	//log.Debug("Ping64 Start")

	go p64.goRun()
	p = p64
	return p
}

func (p *Pinger64Struct) Width() (s PingWidth) {
	return PingWidth64
}

func (p *Pinger64Struct) TurnOn() {
	p.active = true
}

func (p *Pinger64Struct) TurnOff() {
	p.active = false
}

func (p *Pinger64Struct) RecvPong(pong Pong) {
	p.pongch <- pong
}

func (p *Pinger64Struct) GetPingChan() <-chan Ping {
	return p.pingch
}

func (p *Pinger64Struct) goRun() {
	pingmap := make(map[Ping64]time.Time)
	tick := time.NewTicker(p.freq)
	tickch := tick.C
	countch := p.counter.GetCountCh()

	log.GDebug("Pinger Start")
	defer log.GDebug("Pinger Stop")

	defer tick.Stop()
	defer close(p.pingch)
	defer close(p.Errorch)

	for {
		select {
		case <-tickch:
			if p.active {
				var c counter.Counter
				var d Ping64
				c = <-countch
				p64 := Ping64(c)
				pingmap[p64] = time.Now()
				d = p64
				p.pingch <- d
			}

		case pong := <-p.pongch:
			if t, ok := pingmap[Ping64(pong)]; ok {
				delete(pingmap, Ping64(pong))
				dur := time.Since(t)
				if dur > p.timeout {
					log.Warnf("Ping Timeout: %d", pong)
					p.Errorch <- fmt.Errorf("Ping Timeout: %d", pong)
				}
			} else {
				log.Warnf("Ping Missing: %d", pong)
				p.Errorch <- fmt.Errorf("Ping Missing: %d", pong)
			}

		case <-p.cx.DoneChan():
			return
		}

		// Scan the ping table
		expired := []Ping64{}
		for count, t := range pingmap {
			dur := time.Since(t)
			if dur > p.timeout {
				log.Warnf("Ping Expired: %d", count)
				p.Errorch <- fmt.Errorf("Ping Expired: %d", count)
				expired = append(expired, count)
			}
		}
		for _, count := range expired {
			delete(pingmap, count)
		}
	}
}

func (ps *Pinger64Struct) NewPong(pong []byte) (p Pong) {
	if len(pong) != 8 {
		log.FatalfStack("Not enough pong data:%0x", pong)
	}
	c64 := ps.counter.NewByteCounter(pong)
	p = c64
	return p
}
