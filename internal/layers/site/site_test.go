package site

import (
	"testing"

	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestNewSite32_compile(t *testing.T) {
	cx := ctx.NewContext()
	defer cx.Cancel()
	NewSite32(cx, nil, nil)

}
func TestAssembleSite_no_config(t *testing.T) {
	cx := ctx.NewContext()
	defer cx.Cancel()

	if _, err := AssembleSite(nil, &settings.ConfigSiteStruct{}); err == nil {
		t.Error("Expected Error")
	}
	if _, err := AssembleSite(cx, nil); err == nil {
		t.Error("Expected Error")
	}

}
