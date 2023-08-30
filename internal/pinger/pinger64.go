package pinger

import (
	"fmt"
	"time"

	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type Pinger64Struct struct {
	cx      *ctx.Ctx
	active  bool
	freq    time.Duration
	timeout time.Duration
	counter *counter.Counter64Struct
	Pingch  chan counter.Counter64
	Pongch  chan counter.Counter64
	Errorch chan error
}

func NewPinger64(ctx *ctx.Ctx, freq time.Duration, timeout time.Duration) (p *Pinger64Struct) {
	p = &Pinger64Struct{
		cx:      ctx,
		active:  false,
		freq:    freq,
		timeout: timeout,
		counter: counter.NewCounter64(ctx),
		Pingch:  make(chan counter.Counter64),
		Pongch:  make(chan counter.Counter64),
		Errorch: make(chan error),
	}

	//log.Debug("Ping64 Start")

	go p.goRun()

	return p
}

func (p *Pinger64Struct) TurnOn() {
	p.active = true
}

func (p *Pinger64Struct) TurnOff() {
	p.active = false
}

func (p *Pinger64Struct) goRun() {
	pingmap := make(map[counter.Counter64]time.Time)
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
		expired := []counter.Counter64{}
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
