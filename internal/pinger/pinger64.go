package pinger

import (
	"fmt"
	"time"

	"github.com/seanmcadam/counter"
	"github.com/seanmcadam/ctx"
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/octolib/errors"
)

type Ping64 counter.Count
type Pong64 counter.Count

type Pinger64Struct struct {
	cx      *ctx.Ctx
	active  bool
	freq    time.Duration
	timeout time.Duration
	counter counter.Counter
	pingch  chan Ping
	pongch  chan Pong
	Errorch chan error
}

func NewBytePing64(b []byte) (p Ping, err error) {
	if len(b) != 4 {
		return nil, errors.ErrCounterBadParameter(log.Errf("Count data len:%d, :%0x", len(b), b))
	}

	c, err := counter.ByteToCount(b)
	if err == nil {
		return nil, err
	}

	return c, nil

}

func NewBytePong64(b []byte) (p Ping, err error) {
	if len(b) != 8 {
		return nil, errors.ErrCounterBadParameter(log.Errf("Count data len:%d, :%0x", len(b), b))
	}

	c, err := counter.ByteToCount(b)
	if err == nil {
		return nil, err
	}

	return c, nil

}

func NewPinger64(ctx *ctx.Ctx, freq time.Duration, timeout time.Duration) (p PingerStruct, err error) {
	if ctx == nil {
		return p, errors.ErrPingerBadParameter(log.Errf(""))
	}

	p64 := &Pinger64Struct{
		cx:      ctx,
		active:  false,
		freq:    freq,
		timeout: timeout,
		counter: counter.New(ctx, counter.BIT64),
		pingch:  make(chan Ping),
		pongch:  make(chan Pong),
		Errorch: make(chan error),
	}

	//log.Debug("Ping64 Start")

	go p64.goRun()
	p = p64
	return p, nil
}

func (p *Pinger64Struct) Width() (s PingWidth) {
	return PingWidth64
}

func (p *Pinger64Struct) TurnOn() {
	if p == nil {
		return
	}
	p.active = true
}

func (p *Pinger64Struct) TurnOff() {
	if p == nil {
		return
	}
	p.active = false
}

func (p *Pinger64Struct) RecvPong(pong Pong) {
	if p == nil {
		return
	}

	p.pongch <- pong
}

func (p *Pinger64Struct) GetPingChan() <-chan Ping {
	return p.pingch
}

func (p *Pinger64Struct) goRun() {
	pingmap := make(map[Ping64]time.Time)
	tick := time.NewTicker(p.freq)
	tickch := tick.C
	//countch := p.counter.GetCountCh()

	log.GDebug("Pinger64 Run")
	defer log.GDebug("Pinger Stop")

	defer tick.Stop()
	defer close(p.pingch)
	defer close(p.Errorch)

	for {
		select {
		case <-tickch:
			if p.active {
				var c counter.Count
				var d Ping64
				c = p.counter.Next()
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

func (ps *Pinger64Struct) NewPong(pong []byte) (p Pong, err error) {
	if len(pong) != 8 {
		return nil, errors.ErrCounterBadParameter(log.Errf("Not enough pong data:%0x", pong))
	}

	c64, err := counter.ByteToCount(pong)
	if p == nil {
		return nil, err
	}

	p = c64
	return p, nil
}
