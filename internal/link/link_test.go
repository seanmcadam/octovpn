package link

import (
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/octolib/ctx"
)

// @TODO - validate the Notice actions

func TestLinkState_compile(t *testing.T) {
	cx := ctx.NewContext()
	_ = NewLinkState(cx)
	cx.Cancel()
}

func TestLinkState_StateToggles(t *testing.T) {
	cx := ctx.NewContext()
	ls := NewLinkState(cx)

	if ls.GetState() != LinkStateNONE {
		t.Error("State not NONE")
	}

	ls.Up()

	if ls.GetState() != LinkStateUP {
		t.Error("State not Up")
	}

	ls.Down()

	if ls.GetState() != LinkStateDOWN {
		t.Error("State not Down")
	}

	cx.Cancel()

}

func TestLinkState_StateTogglesChannel(t *testing.T) {
	cx := ctx.NewContext()

	ls := NewLinkState(cx)

	ch := ls.LinkStateCh()
	ls.Up()
	select {
	case ns := <-ch:
		if ns.State() != LinkStateUP {
			t.Error("State did not toggle to Up")
		}
	case <-time.After(time.Millisecond):
		t.Error("Timeout...")
	}

	ch = ls.LinkStateCh()
	ls.Up()
	select {
	case ns := <-ch:
		if ns.State() != LinkStateUP {
			t.Error("State toggled")
		}
	case <-time.After(time.Millisecond):
	}

	ch = ls.LinkStateCh()
	ls.Down()
	select {
	case ns := <-ch:
		if ns.State() != LinkStateDOWN {
			t.Error("State did not toggle to Down")
		}
	case <-time.After(time.Second):
		t.Error("Timeout...")
	}

	cx.Cancel()
}
