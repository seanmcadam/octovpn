package counter

import (
	"encoding/binary"

	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type Counter32 uint32

type Counter32Struct struct {
	cx      *ctx.Ctx
	CountCh chan Counter
}

func MakeCounter32(c32 uint32) (c Counter) {
	cc := Counter32(c32)
	c = &cc
	return c
}

func (c *Counter32) ToByte() (b []byte) {
	b = make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(*c))
	return b
}

func (c *Counter32) Width() CounterWidth {
	return CounterWidth32
}

func (c *Counter32) Uint() interface{} {
	var c32 uint32
	c32 = uint32(*c)
	return c32
}

func (c *Counter32) Copy() Counter {
	return MakeCounter32(uint32(*c))
	//var copy = *c
	//return Counter(&copy)
}

func NewByteCounter32(b []byte) (c Counter) {
	if len(b) != 4 {
		log.FatalfStack("Count data len:%d, :%0x", len(b), b)
	}
	c32 := Counter32(binary.LittleEndian.Uint32(b))
	c = &c32
	return c
}

func NewCounter32(ctx *ctx.Ctx) (c CounterStruct) {
	c32 := &Counter32Struct{
		cx:      ctx,
		CountCh: make(chan Counter, 5),
	}
	go c32.goCount()

	c = c32
	return c
}

func (*Counter32Struct) NewByteCounter(b []byte) (c Counter) {
	if len(b) != 4 {
		log.FatalfStack("Count data len:%d, :%0x", len(b), b)
	}
	c32 := Counter32(binary.LittleEndian.Uint32(b))
	c = &c32
	return c
}

func (c *Counter32Struct) Width() CounterWidth {
	return CounterWidth32
}

func (c *Counter32Struct) Next() (Counter) {
	return <-c.CountCh
}

func (c *Counter32Struct) GetCountCh() (ch <-chan Counter) {
	return c.CountCh
}

// goCounter()
// Generates a UniqueID (int) and returns via supplied channel
// Runs forever
func (c *Counter32Struct) goCount() {

	defer c.emptych()

	var counter Counter32 = 0
	for {
		counter += 1
		localc := counter
		select {
		case c.CountCh <- &localc:
		case <-c.cx.DoneChan():
			return
		}
	}
}

func (c *Counter32Struct) emptych() {
	for {
		select {
		case <-c.CountCh:
		default:
			close(c.CountCh)
			return
		}
	}
}
