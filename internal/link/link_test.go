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

	ls1 := NewLinkState(cx)
	ls1.NoLink()

	ls := NewLinkState(cx)
	stateCh := ls.LinkStateCh()
	upCh := ls.LinkUpCh()
	connectCh := ls.LinkConnectCh()
	ls.NoLink()
	ls.AddLinkStateCh(ls1)
	time.Sleep(1*time.Millisecond)

	ls.AddLinkStateCh(ls1)
	time.Sleep(1*time.Millisecond)
	//upCh := ls.LinkConnectCh()
	//upCh := ls.LinkUpCh()
	//upCh := ls.LinkUpCh()

	ls1.Connected()
	time.Sleep(1*time.Millisecond)

	select {
	case <-stateCh:
	case <-upCh:
	case <-connectCh:
	default:
		t.Error("State did not change to Up")
	}

	cx.Cancel()
}
