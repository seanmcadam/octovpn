package counter

import (
	"encoding/binary"

	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/errors"
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
	if c == nil {
		return b
	}

	b = make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(*c))
	return b
}

func (c *Counter64) Width() CounterWidth {
	if c == nil {
		return 0
	}

	return CounterWidth64
}

func (c *Counter64) Uint() interface{} {
	if c == nil {
		return nil
	}
	var c64 uint64
	c64 = uint64(*c)
	return c64
}

func (c *Counter64) Copy() Counter {
	if c == nil {
		return nil
	}
	return MakeCounter64(uint64(*c))
}

func NewByteCounter64(b []byte) (c Counter, err error) {
	if len(b) != 8 {
		return nil, log.Errf("Count data len:%d, :%0x", len(b), b)
	}
	c64 := Counter64(binary.BigEndian.Uint64(b))
	c = &c64
	return c, nil
}

func NewCounter64(ctx *ctx.Ctx) (c CounterStruct) {
	if ctx == nil {
		log.Fatal()
	}
	c64 := &Counter64Struct{
		cx:      ctx,
		CountCh: make(chan Counter, 5),
	}
	go c64.goCount()
	c = c64
	return c
}

func (*Counter64Struct) NewByteCounter(b []byte) (c Counter, err error) {
	if c == nil {
		return nil, errors.ErrCounterNilMethod(log.Errf(""))
	}

	if len(b) != 4 {
		return nil, log.Errf("Count data len:%d, :%0x", len(b), b)
	}
	c64 := Counter64(binary.BigEndian.Uint64(b))
	c = &c64
	return c, nil
}

func (c *Counter64Struct) Width() CounterWidth {
	if c == nil {
		log.Fatal()
	}
	return CounterWidth64
}

func (c *Counter64Struct) Next() Counter {
	if c == nil {
		log.Fatal()
	}
	return <-c.CountCh
}

func (c *Counter64Struct) GetCountCh() (ch <-chan Counter) {
	if c == nil {
		log.Fatal()
	}
	return c.CountCh
}

// goCounter()
// Generates a UniqueID (int) and returns via supplied channel
// Runs forever
func (c *Counter64Struct) goCount() {
	if c == nil {
		log.Fatal()
	}

	defer c.emptych()

	var counter Counter64 = 1
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

func (c *Counter64Struct) emptych() {
	if c == nil {
		return
	}

	for {
		select {
		case <-c.CountCh:
		default:
			close(c.CountCh)
			return
		}
	}
}
