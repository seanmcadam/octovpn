package site

import (
	"testing"

	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestNewSite32_compile(t *testing.T) {
	cx := ctx.NewContext()
	defer cx.Cancel()
	NewSite32(cx, nil)

}
