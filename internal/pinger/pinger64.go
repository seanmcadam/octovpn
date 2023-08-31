package pinger

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type Ping64 counter.Counter64
type Pong64 counter.Counter64

func NewPong64(pong []byte) *Pong64 {
	if len(pong) < 8 {
		log.FatalfStack("Not enough pong data:%0x", pong)
	}
	p := Pong64(binary.LittleEndian.Uint64(pong))
	return &p
}

func (p *Ping64) ToByte() (ping []byte) {
	ping = make([]byte, 4)
	binary.LittleEndian.PutUint64(ping, uint64(*p))
	return ping
}

func (p *Pong64) ToByte() (pong []byte) {
	pong = make([]byte, 4)
	binary.LittleEndian.PutUint64(pong, uint64(*p))
	return pong
}

type Pinger64Struct struct {
	cx      *ctx.Ctx
	active  bool
	freq    time.Duration
	timeout time.Duration
	counter *counter.Counter64Struct
	Pingch  chan *Ping64
	Pongch  chan *Pong64
	Errorch chan error
}

func NewPinger64(ctx *ctx.Ctx, freq time.Duration, timeout time.Duration) (p *Pinger64Struct) {
	p = &Pinger64Struct{
		cx:      ctx,
		active:  false,
		freq:    freq,
		timeout: timeout,
		counter: counter.NewCounter64(ctx),
		Pingch:  make(chan *Ping64),
		Pongch:  make(chan *Pong64),
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
	pingmap := make(map[Ping64]time.Time)
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
				var c *counter.Counter64
				var d Ping64
				c = <-countch
				d = Ping64(*c)
				p.Pingch <- &d
				pingmap[d] = time.Now()
			}

		case pong := <-p.Pongch:
			if t, ok := pingmap[Ping64(*pong)]; ok {
				delete(pingmap, Ping64(*pong))
				dur := time.Since(t)
				if dur > p.timeout {
					log.Warnf("Ping Timeout: %d", *pong)
					p.Errorch <- fmt.Errorf("Ping Timeout: %d", *pong)
				}
			} else {
				log.Warnf("Ping Missing: %d", *pong)
				p.Errorch <- fmt.Errorf("Ping Missing: %d", *pong)
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
