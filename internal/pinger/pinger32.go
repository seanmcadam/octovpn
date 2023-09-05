package pinger

import (
	"fmt"
	"time"

	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type Ping32 counter.Counter
type Pong32 counter.Counter

type Pinger32Struct struct {
	cx      *ctx.Ctx
	active  bool
	freq    time.Duration
	timeout time.Duration
	counter counter.CounterStruct
	pingch  chan Ping
	pongch  chan Pong
	Errorch chan error
}

func NewBytePing32(b []byte) (p Ping) {
	if len(b) != 4 {
		log.FatalfStack("Count data len:%d, :%0x", len(b), b)
	}

	c := counter.NewByteCounter32(b)
	p = counter.Counter(c)
	return p

}

func NewBytePong32(b []byte) (p Pong) {
	if len(b) != 8 {
		log.FatalfStack("Count data len:%d, :%0x", len(b), b)
	}

	c := counter.NewByteCounter64(b)
	p = counter.Counter(c)
	return p

}

func NewPinger32(ctx *ctx.Ctx, freq time.Duration, timeout time.Duration) (p PingerStruct) {
	p32 := &Pinger32Struct{
		cx:      ctx,
		active:  false,
		freq:    freq,
		timeout: timeout,
		counter: counter.NewCounter32(ctx),
		pingch:  make(chan Ping),
		pongch:  make(chan Pong),
		Errorch: make(chan error),
	}

	go p32.goRun()
	p = p32
	return p
}

func (p *Pinger32Struct) Width() (s PingWidth) {
	return PingWidth32
}

func (p *Pinger32Struct) TurnOn() {
	p.active = true
}

func (p *Pinger32Struct) TurnOff() {
	p.active = false
}

func (p *Pinger32Struct) RecvPong(pong Pong) {
	p.pongch <- pong
}
func (p *Pinger32Struct) GetPingChan() <- chan Ping {
	return p.pingch
}

func (p *Pinger32Struct) goRun() {
	pingmap := make(map[Ping32]time.Time)
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
				var d Ping32
				c = <-countch
				p32 := Ping32(c)
				pingmap[p32] = time.Now()
				d = p32
				p.pingch <- d
			}

		case pong := <-p.pongch:
			if t, ok := pingmap[Ping32(pong)]; ok {
				delete(pingmap, Ping32(pong))
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
		expired := []Ping32{}
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

func (ps *Pinger32Struct) NewPong(pong []byte) (p Pong) {
	if len(pong) != 4 {
		log.FatalfStack("Not enough pong data:%0x", pong)
	}
	p32 := ps.counter.NewByteCounter(pong)
	p = p32
	return p
}
