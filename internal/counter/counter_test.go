package counter

import (
	"testing"

	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestNewCounter64(t *testing.T) {
	cx := ctx.NewContext()
	c := NewCounter64(cx)
	count := <-c.GetCountCh()
	_ = count.ToByte()
	cx.Cancel()

}

func TestNewCounter32(t *testing.T) {
	cx := ctx.NewContext()
	c := NewCounter32(cx)
	count := <-c.GetCountCh()
	_ = count.ToByte()

	cx.Cancel()

}
