package pinger

import (
	"testing"

	"github.com/seanmcadam/octovpn/octolib/ctx"
)

// @TODO
// Add read pingch
// Wait for Ping time out
// Send pong
// Turn On/Off

func TestNewPinger32_compile(t *testing.T) {
	cx := ctx.NewContext()
	_, _ = NewPinger32(cx, 1, 5)

	cx.Cancel()
}

func TestNewPinger64_compile(t *testing.T) {
	cx := ctx.NewContext()
	_, _ = NewPinger64(cx, 1, 5)

	cx.Cancel()
}
