package pinger

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type Ping32 counter.Counter32
type Pong32 counter.Counter32

func NewPong32(pong []byte) *Pong32 {
	if len(pong) < 4 {
		log.FatalfStack("Not enough pong data:%0x", pong)
	}
	p := Pong32(binary.LittleEndian.Uint32(pong))
	return &p
}

func (p *Ping32) ToByte() (ping []byte) {
	ping = make([]byte, 4)
	binary.LittleEndian.PutUint32(ping, uint32(*p))
	return ping
}

func (p *Pong32) ToByte() (pong []byte) {
	pong = make([]byte, 4)
	binary.LittleEndian.PutUint32(pong, uint32(*p))
	return pong
}

type Pinger32Struct struct {
	cx      *ctx.Ctx
	active  bool
	freq    time.Duration
	timeout time.Duration
	counter *counter.Counter32Struct
	Pingch  chan *Ping32
	Pongch  chan *Pong32
	Errorch chan error
}

func NewPinger32(ctx *ctx.Ctx, freq time.Duration, timeout time.Duration) (p *Pinger32Struct) {
	p = &Pinger32Struct{
		cx:      ctx,
		active:  false,
		freq:    freq,
		timeout: timeout,
		counter: counter.NewCounter32(ctx),
		Pingch:  make(chan *Ping32),
		Pongch:  make(chan *Pong32),
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
	pingmap := make(map[Ping32]time.Time)
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
				var c *counter.Counter32
				var d Ping32
				c = <-countch
				d = Ping32(*c)
				p.Pingch <- &d
				pingmap[d] = time.Now()
			}

		case pong := <-p.Pongch:
			if t, ok := pingmap[Ping32(*pong)]; ok {
				delete(pingmap, Ping32(*pong))
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
