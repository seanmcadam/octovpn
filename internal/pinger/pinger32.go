package pinger

import (
	"fmt"
	"time"

	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/errors"
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

func NewBytePing32(b []byte) (p Ping, err error) {
	if len(b) != 4 {
		return nil, errors.ErrCounterBadParameter(log.Errf("Count data len:%d, :%0x", len(b), b))
	}

	c, err := counter.NewByteCounter32(b)
	if err == nil {
		return nil, err
	}
	p = counter.Counter(c)
	return p, nil

}

func NewBytePong32(b []byte) (p Pong, err error) {
	if len(b) != 8 {
		return nil, errors.ErrCounterBadParameter(log.Errf("Count data len:%d, :%0x", len(b), b))
	}

	c, err := counter.NewByteCounter64(b)
	if err != nil {
		return nil, err
	}

	p = counter.Counter(c)
	return p, nil

}

func NewPinger32(ctx *ctx.Ctx, freq time.Duration, timeout time.Duration) (p PingerStruct, err error) {
	if ctx == nil{
		return p, errors.ErrPingerBadParameter(log.Errf(""))
	}

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
	return p, nil
}

func (p *Pinger32Struct) Width() (s PingWidth) {
	return PingWidth32
}

func (p *Pinger32Struct) TurnOn() {
	if p == nil {
		return
	}

	p.active = true
}

func (p *Pinger32Struct) TurnOff() {
	if p == nil {
		return
	}

	p.active = false
}

func (p *Pinger32Struct) RecvPong(pong Pong) {
	if p == nil {
		return 
	}

	p.pongch <- pong
}
func (p *Pinger32Struct) GetPingChan() <-chan Ping {
	if p == nil {
		return nil
	}

	return p.pingch
}

func (p *Pinger32Struct) goRun() {
	if p == nil {
		return
	}

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

func (ps *Pinger32Struct) NewPong(pong []byte) (p Pong, err error) {
	if p == nil {
		return
	}

	if len(pong) != 4 {
		log.FatalfStack("Not enough pong data:%0x", pong)
	}

	p32, err := ps.counter.NewByteCounter(pong)
	if p == nil {
		return nil, err
	}

	p = p32
	return p, nil
}
