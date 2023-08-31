package counter

import (
	"encoding/binary"

	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type Counter64 uint64

type Counter64Struct struct {
	cx      *ctx.Ctx
	CountCh chan *Counter64
}

func (c *Counter64) ToByte() (b []byte) {
	b = make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(*c))
	return b
}

func NewCounter64(ctx *ctx.Ctx) (c *Counter64Struct) {
	c = &Counter64Struct{
		cx:      ctx,
		CountCh: make(chan *Counter64, 5),
	}
	go c.goCount()
	return c
}

func (c *Counter64Struct) GetCountCh() (ch chan *Counter64) {
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
