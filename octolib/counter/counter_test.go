package counter

import (
	"testing"

	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestNewCounter64(t *testing.T) {
	cx := ctx.NewContext()
	_ = NewCounter64(cx)

	cx.Cancel()

}
