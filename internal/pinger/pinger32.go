package pinger

import (
	"fmt"
	"time"

	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type Pinger32Struct struct {
	cx      *ctx.Ctx
	active  bool
	freq    time.Duration
	timeout time.Duration
	counter *counter.Counter32Struct
	Pingch  chan counter.Counter32
	Pongch  chan counter.Counter32
	Errorch chan error
}

func NewPinger32(ctx *ctx.Ctx, freq time.Duration, timeout time.Duration) (p *Pinger32Struct) {
	p = &Pinger32Struct{
		cx:      ctx,
		active:  false,
		freq:    freq,
		timeout: timeout,
		counter: counter.NewCounter32(ctx),
		Pingch:  make(chan counter.Counter32),
		Pongch:  make(chan counter.Counter32),
		Errorch: make(chan error),
	}

	go p.goRun()

	return p
}

func (p *Pinger32Struct) TurnOn() {
	p.active = true
}

func (p *Pinger32Struct) TurnOff() {
	p.active = false
}

func (p *Pinger32Struct) goRun() {
	pingmap := make(map[counter.Counter32]time.Time)
	tick := time.NewTicker(p.freq)
	tickch := tick.C
	countch := p.counter.GetCountCh()

	log.GDebug("Pinger Start")
	defer log.GDebug("Pinger Stop")

	defer tick.Stop()
	defer close(p.Pingch)
	defer close(p.Errorch)

	for {
		select {
		case <-tickch:
			if p.active {
				c := <-countch
				p.Pingch <- c
				pingmap[c] = time.Now()
			}

		case pong := <-p.Pongch:
			if t, ok := pingmap[pong]; ok {
				delete(pingmap, pong)
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
		expired := []counter.Counter32{}
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
