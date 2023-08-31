package counter

import (
	"encoding/binary"

	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type Counter32 uint32

type Counter32Struct struct {
	cx      *ctx.Ctx
	CountCh chan *Counter32
}

func (c *Counter32) ToByte() (b []byte) {
	b = make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(*c))
	return b
}

func NewCounter32(ctx *ctx.Ctx) (c *Counter32Struct) {
	c = &Counter32Struct{
		cx:      ctx,
		CountCh: make(chan *Counter32, 5),
	}
	go c.goCount()
	return c
}

func (c *Counter32Struct) GetCountCh() (ch chan *Counter32) {
	return c.CountCh
}

// goCounter()
// Generates a UniqueID (int) and returns via supplied channel
// Runs forever
func (c *Counter32Struct) goCount() {

	log.GDebug("Counter32 Start")
	defer log.GDebug("Counter32 Stop")
	defer c.emptych()

	var counter Counter32 = 1
	for {
		select {
		case c.CountCh <- &counter:
			counter += 1
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
