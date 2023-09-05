package counter

import (
	"encoding/binary"

	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type Counter64 uint64

type Counter64Struct struct {
	cx      *ctx.Ctx
	CountCh chan Counter
}

func MakeCounter64(c64 uint64) (c Counter) {
	cc := Counter64(c64)
	c = &cc
	return c
}

func (c *Counter64) ToByte() (b []byte) {
	b = make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(*c))
	return b
}

func (c *Counter64) Width() CounterWidth {
	return CounterWidth64
}

func (c *Counter64) Uint() interface{} {
	var c64 uint64
	c64 = uint64(*c)
	return c64
}

func (c *Counter64) Copy() Counter {
	return MakeCounter64(uint64(*c))
	var copy = *c
	return Counter(&copy)
}

func NewByteCounter64(b []byte) (c Counter) {
	if len(b) != 8 {
		log.FatalfStack("Count data len:%d, :%0x", len(b), b)
	}
	c64 := Counter64(binary.LittleEndian.Uint64(b))
	c = &c64
	return c
}

func NewCounter64(ctx *ctx.Ctx) (c CounterStruct) {
	c64 := &Counter64Struct{
		cx:      ctx,
		CountCh: make(chan Counter, 5),
	}
	go c64.goCount()
	c = c64
	return c
}

func (*Counter64Struct) NewByteCounter(b []byte) (c Counter) {
	if len(b) != 4 {
		log.FatalfStack("Count data len:%d, :%0x", len(b), b)
	}
	c64 := Counter64(binary.LittleEndian.Uint64(b))
	c = &c64
	return c
}

func (c *Counter64Struct) Width() CounterWidth {
	return CounterWidth64
}

func (c *Counter64Struct) GetCountCh() (ch <-chan Counter) {
	return c.CountCh
}

// goCounter()
// Generates a UniqueID (int) and returns via supplied channel
// Runs forever
func (c *Counter64Struct) goCount() {

	log.GDebug("Counter64 Start")
	defer log.GDebug("Counter64 Stop")
	defer c.emptych()

	var counter Counter64 = 1
	for {
		select {
		case c.CountCh <- &counter:
			counter += 1
		case <-c.cx.DoneChan():
			return
		}
	}
}

func (c *Counter64Struct) emptych() {
	for {
		select {
		case <-c.CountCh:
		default:
			close(c.CountCh)
			return
		}
	}
}
