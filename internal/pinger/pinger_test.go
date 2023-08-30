package pinger

import (
	"testing"

	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestNewPinger64(t *testing.T) {
	cx := ctx.NewContext()
	_ = NewPinger64(cx, 1, 5)

	cx.Cancel()
}


func TestNewPinger32(t *testing.T) {
	cx := ctx.NewContext()
	_ = NewPinger32(cx, 1, 5)

	cx.Cancel()
}
