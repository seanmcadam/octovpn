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

func TestLinkState_nil_method(t *testing.T) {

	var ls *LinkStateStruct

	ls.goRecv()
	ls.goRun()
	ls.processMessage(0)
	ls.setState(0)
	ls.NoLink()
	ls.Listen()
	ls.Link()
	ls.Chal()
	ls.Auth()
	ls.Connected()
	ls.Close()
	ls.LinkStateCh()
	ls.LinkNoticeCh()
	ls.LinkNoticeStateCh()
	ls.LinkChalCh()
	ls.LinkAuthCh()
	ls.LinkLinkCh()
	ls.LinkConnectCh()
	ls.LinkNoLinkCh()
	ls.LinkListenCh()
	ls.LinkUpDownCh()
	ls.LinkUpCh()
	ls.LinkDownCh()
	ls.LinkCloseCh()
	ls.LinkLossCh()
	ls.LinkLatencyCh()
	ls.LinkSaturationCh()
	ls.refreshRecvLinks()
	ls.AddLinkStateCh(nil)
	ls.AddLinkNoticeCh(nil)
	ls.AddLinkNoticeStateCh(nil)
	ls.AddLinkLinkCh(nil)
	ls.AddLinkUpDownCh(nil)
	ls.AddLinkUpCh(nil)
	ls.AddLinkConnectCh(nil)
	ls.AddLinkDownCh(nil)
	ls.AddLinkCloseCh(nil)
	ls.AddLinkLossCh(nil)
	ls.AddLinkLatencyCh(nil)
	ls.AddLinkSaturationCh(nil)

}
func TestLinkState_StateToggles(t *testing.T) {
	cx := ctx.NewContext()
	ls := NewLinkState(cx)

	if ls.GetState() != LinkStateNONE {
		t.Error("State not NONE")
	}

	ls.Connected()

	if ls.GetState() != LinkStateCONNECTED {
		t.Error("State not Up")
	}

	ls.NoLink()

	if ls.GetState() != LinkStateNOLINK {
		t.Error("State not Down")
	}

	cx.Cancel()

}

func TestLinkState_StateTogglesChannel(t *testing.T) {
	cx := ctx.NewContext()

	ls := NewLinkState(cx)

	ch := ls.LinkStateCh()
	ls.Connected()
	select {
	case ns := <-ch:
		if ns.State() != LinkStateCONNECTED {
			t.Error("State did not toggle to Up")
		}
	case <-time.After(time.Millisecond):
		t.Error("Timeout...")
	}

	ch = ls.LinkStateCh()
	ls.Connected()
	select {
	case ns := <-ch:
		if ns.State() != LinkStateCONNECTED {
			t.Error("State toggled")
		}
	case <-time.After(time.Millisecond):
	}

	ch = ls.LinkStateCh()
	ls.NoLink()
	select {
	case ns := <-ch:
		if ns.State() != LinkStateNOLINK {
			t.Error("State did not toggle to Down")
		}
	case <-time.After(time.Second):
		t.Error("Timeout...")
	}

	cx.Cancel()
}

func TestLinkState_UP_send(t *testing.T) {
	cx := ctx.NewContext()
	defer cx.Cancel()

	ls1 := newLinkState(cx, "Slave")
	ls1.NoLink()

	ls := newLinkState(cx, "Master")
	connectedCh := ls.LinkConnectCh()
	connectedCh1 := ls1.LinkConnectCh()
	ls.AddLinkConnectCh(ls1)

	ls.NoLink()

	ls1.Connected()
	time.Sleep(1 * time.Millisecond)

	select {
	case <-connectedCh:
	case <-connectedCh1:
	default:
		t.Error("StateCh did not change to Up")
	}

	cx.Cancel()
}
