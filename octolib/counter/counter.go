package counter

import (
	"encoding/binary"

	"github.com/seanmcadam/octovpn/octolib/log"
)

type Counter64 uint64

type Counter64Struct struct {
	CountCh chan Counter64
	closech chan interface{}
}

func (c Counter64) ToByte() (b []byte) {
	b = make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(c))
	return b
}

func NewCounter64() (c *Counter64Struct) {
	c = &Counter64Struct{
		CountCh: make(chan Counter64, 5),
		closech: make(chan interface{}),
	}
	go c.goCount()
	return c
}

func (c *Counter64Struct) GetCountCh() (ch chan Counter64) {
	return c.CountCh
}

func (c *Counter64Struct) Close() {
	close(c.closech)
}

// goCounter()
// Generates a UniqueID (int) and returns via supplied channel
// Runs forever
func (c *Counter64Struct) goCount() {

	log.GDebug("Start")
	defer log.GDebug("Stop")
	defer c.emptych()

	var counter Counter64 = 1
	for {
		select {
		case c.CountCh <- counter:
			counter += 1
		case <-c.closech:
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
