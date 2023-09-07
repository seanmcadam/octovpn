package link

import (
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/octolib/ctx"
)

func TestLinkState_compile(t *testing.T) {
	cx := ctx.NewContext()
	_ = NewLinkState(cx)
	cx.Cancel()
}

func TestLinkState_StateToggles(t *testing.T) {
	cx := ctx.NewContext()
	ls := NewLinkState(cx)

	if ls.GetState() != LinkStateNone {
		t.Error("State not NONE")
	}

	ls.ToggleState(LinkStateUp)

	if ls.GetState() != LinkStateUp {
		t.Error("State not Up")
	}

	ls.ToggleState(LinkStateDown)

	if ls.GetState() != LinkStateDown {
		t.Error("State not Down")
	}

	cx.Cancel()

}

func TestLinkState_StateTogglesChannel(t *testing.T) {
	cx := ctx.NewContext()

	ls := NewLinkState(cx)

	ch := ls.StateToggleCh()
	ls.ToggleState(LinkStateUp)
	select {
	case state := <-ch:
		if state != LinkStateUp {
			t.Error("State did not toggle to Up")
		}
	case <-time.After(time.Microsecond):
		t.Error("Timeout...")
	}

	ch = ls.StateToggleCh()
	ls.ToggleState(LinkStateUp)
	select {
	case state := <-ch:
		if state != LinkStateUp {
			t.Error("State toggled")
		}
	case <-time.After(time.Microsecond):
	}

	ch = ls.StateToggleCh()
	ls.ToggleState(LinkStateDown)
	select {
	case state := <-ch:
		if state != LinkStateDown {
			t.Error("State did not toggle to Down")
		}
	case <-time.After(time.Microsecond):
		t.Error("Timeout...")
	}

	cx.Cancel()
}
