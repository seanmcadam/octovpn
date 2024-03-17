package connmgr

import (
	"testing"

	"github.com/seanmcadam/ctx"
)

func TestCompile(t *testing.T) {
}
func TestNew(t *testing.T) {

	cx := ctx.Ctx.New()

	cm := New(cx, "")

	_ = cm

}
