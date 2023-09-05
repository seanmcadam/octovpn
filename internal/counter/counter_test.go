package counter

import (
	"testing"

	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestNewCounter_interface32(t *testing.T) {
	var count Counter
	cx := ctx.NewContext()
	c32 := NewCounter32(cx)

	width32 := c32.Width()
	if width32 != CounterWidth32 {
		t.Fatal("Size mismatch")
	}

	var c CounterStruct = c32
	count = <-c.GetCountCh()
	count = <-c.GetCountCh()
	count = <-c.GetCountCh()
	count = <-c.GetCountCh()
	_ = count.ToByte()
	cx.Cancel()

}

func TestNewCounter_interface64(t *testing.T) {
	var count Counter
	cx := ctx.NewContext()
	c64 := NewCounter64(cx)

	width64 := c64.Width()
	if width64 != CounterWidth64 {
		t.Fatal("Size mismatch")
	}

	var c CounterStruct = c64
	count = <-c.GetCountCh()
	count = <-c.GetCountCh()
	count = <-c.GetCountCh()
	count = <-c.GetCountCh()
	_ = count.ToByte()
	cx.Cancel()

}

