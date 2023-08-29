package tracker

import (
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestNewCompile(t *testing.T) {

	cx := ctx.NewContext()

	_ = NewTracker(cx, 1*time.Second)

	cx.Cancel()

}
