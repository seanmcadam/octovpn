package link

import (
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// @TODO - validate the Notice actions

func TestLinkState_compile(t *testing.T) {
	cx := ctx.NewContext()
	_ = NewLinkState(cx)
	cx.Cancel()
}

func TestLinkState_nil_method(t *testing.T) {

	var ls *LinkStateStruct
	var lc *LinkChan

	//lc.LinkCh()
	lc.send(nil)

	ls.goRecv()
	ls.setState(0)
	ls.Auth()
	ls.NoLink()
	ls.Start()
	ls.Chal()
	ls.Close()
	ls.Connected()
	ls.Link()
	ls.LinkAuthCh()
	ls.LinkCloseCh()
	ls.LinkChalCh()
	ls.LinkConnectCh()
	ls.LinkDownCh()
	ls.LinkLatencyCh()
	ls.LinkLinkCh()
	ls.LinkListenCh()
	ls.LinkLossCh()
	ls.LinkNoticeCh()
	ls.LinkNoticeStateCh()
	ls.LinkNoLinkCh()
	ls.LinkSaturationCh()
	ls.LinkStartCh()
	ls.LinkStateCh()
	ls.LinkUpCh()
	ls.LinkUpDownCh()
	ls.Listen()
	ls.AddLinkCloseCh(nil)
	ls.AddLinkChalCh(nil)
	ls.AddLinkConnectCh(nil)
	ls.AddLinkDownCh(nil)
	ls.AddLinkLatencyCh(nil)
	ls.AddLinkLinkCh(nil)
	ls.AddLinkListenCh(nil)
	ls.AddLinkLossCh(nil)
	ls.AddLinkNoLinkCh(nil)
	ls.AddLinkNoticeCh(nil)
	ls.AddLinkNoticeStateCh(nil)
	ls.AddLinkSaturationCh(nil)
	ls.AddLinkStartCh(nil)
	ls.AddLinkStateCh(nil)
	ls.AddLinkUpCh(nil)
	ls.AddLinkUpDownCh(nil)

	ls.LinkStateCh()
	ls.LinkNoticeCh()
	ls.LinkNoticeStateCh()
	ls.LinkChalCh()
	ls.LinkAuthCh()
	ls.LinkLinkCh()
	ls.LinkConnectCh()
	ls.LinkUpDownCh()
	ls.LinkListenCh()
	ls.LinkNoLinkCh()
	ls.LinkUpCh()
	ls.LinkStartCh()
	ls.LinkDownCh()
	ls.LinkCloseCh()
	ls.LinkLossCh()
	ls.LinkLatencyCh()
	ls.LinkSaturationCh()

}

func TestLinkState_test_links(t *testing.T) {
	cx := ctx.NewContext()
	ls := NewLinkState(cx)

	ls.LinkStateCh()
	ls.LinkNoticeCh()
	ls.LinkNoticeStateCh()
	ls.LinkChalCh()
	ls.LinkAuthCh()
	ls.LinkLinkCh()
	ls.LinkConnectCh()
	ls.LinkUpDownCh()
	ls.LinkListenCh()
	ls.LinkNoLinkCh()
	ls.LinkUpCh()
	ls.LinkStartCh()
	ls.LinkDownCh()
	ls.LinkCloseCh()
	ls.LinkLossCh()
	ls.LinkLatencyCh()
	ls.LinkSaturationCh()
}

func TestLinkState_StateToggles(t *testing.T) {
	cx := ctx.NewContext()
	ls := NewLinkState(cx)

	if ls.GetState() != LinkStateNONE {
		t.Error("State not NONE")
	}

	ls.Auth()
	if ls.GetState() != LinkStateAUTH {
		t.Error("State not AUTH")
	}

	ls.Chal()
	if ls.GetState() != LinkStateCHAL {
		t.Error("State not CHAL")
	}

	ls.Listen()
	if ls.GetState() != LinkStateLISTEN {
		t.Error("State not LISTEN")
	}

	ls.Start()
	if ls.GetState() != LinkStateSTART {
		t.Error("State not START")
	}

	ls.Link()
	if ls.GetState() != LinkStateLINK {
		t.Error("State not LINK")
	}

	ls.Connected()
	if ls.GetState() != LinkStateCONNECTED {
		t.Error("State not Connected")
	}

	ls.NoLink()

	if ls.GetState() != LinkStateNOLINK {
		t.Error("State not NoLink")
	}

	cx.Cancel()

}

func TestLinkState_StateTogglesChannel(t *testing.T) {
	cx := ctx.NewContext()

	ls := NewLinkState(cx)

	ch := ls.LinkStateCh()
	ls.NoLink()
	select {
	case ns := <-ch:
		if ns.State() != LinkStateNOLINK {
			t.Error("State did not toggle to NOLINK")
		}
	case <-time.After(time.Millisecond):
		t.Fatal("Timeout...")
	}

	ch = ls.LinkStateCh()
	ls.Start()
	select {
	case ns := <-ch:
		if ns.State() != LinkStateSTART {
			t.Error("State did not toggle to START")
		}
	case <-time.After(time.Millisecond):
		t.Fatal("Timeout...")
	}

	ch = ls.LinkStateCh()
	ls.Auth()
	select {
	case ns := <-ch:
		if ns.State() != LinkStateAUTH {
			t.Error("State did not toggle to AUTH")
		}
	case <-time.After(time.Millisecond):
		t.Fatal("Timeout...")
	}

	ch = ls.LinkStateCh()
	ls.Chal()
	select {
	case ns := <-ch:
		if ns.State() != LinkStateCHAL {
			t.Error("State did not toggle to CHAL")
		}
	case <-time.After(time.Millisecond):
		t.Fatal("Timeout...")
	}

	ch = ls.LinkStateCh()
	ls.Listen()
	select {
	case ns := <-ch:
		if ns.State() != LinkStateLISTEN {
			t.Error("State did not toggle to LISTEN")
		}
	case <-time.After(time.Millisecond):
		t.Fatal("Timeout...")
	}

	ch = ls.LinkStateCh()
	ls.Link()
	select {
	case ns := <-ch:
		if ns.State() != LinkStateLINK {
			t.Error("State did not toggle to LINK")
		}
	case <-time.After(time.Millisecond):
		t.Fatal("Timeout...")
	}

	ch = ls.LinkStateCh()

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

func TestLinkState_send_recv_link(t *testing.T) {
	cx := ctx.NewContext()
	defer cx.Cancel()

	rec := newLinkState(cx, "reciever")
	rec.NoLink()
	recupdownch := rec.LinkUpDownCh()

	send := newLinkState(cx, "sender")
	send.NoLink()

	log.Infof("Rec:%s, Send:%s", rec.GetState(), send.GetState())

	rec.AddLinkStateCh(send)
	time.Sleep(time.Millisecond)

	send.Connected()
	state := <-recupdownch
	log.Debugf("Up State:%s", state)

	send.NoLink()
	state = <-recupdownch
	log.Debugf("Down State:%s", state)

	cx.Cancel()
}

func TestLinkState_to_other_links(t *testing.T) {
	cx := ctx.NewContext()
	defer cx.Cancel()

	rec := newLinkState(cx, "reciever")
	rec.NoLink()

	send := newLinkState(cx, "sender")
	send.NoLink()

	log.Infof("Rec:%s, Send:%s", rec.GetState(), send.GetState())

	rec.AddLinkAuthCh(send)
	rec.AddLinkChalCh(send)
	rec.AddLinkConnectCh(send)
	rec.AddLinkLinkCh(send)
	rec.AddLinkListenCh(send)
	rec.AddLinkNoLinkCh(send)
	rec.AddLinkStartCh(send)

	var test map[LinkStateType]func() = make(map[LinkStateType]func())
	test[LinkStateNOLINK] = send.NoLink
	test[LinkStateAUTH] = send.Auth
	test[LinkStateCHAL] = send.Chal
	test[LinkStateCONNECTED] = send.Connected
	test[LinkStateLINK] = send.Link
	test[LinkStateLISTEN] = send.Listen
	test[LinkStateSTART] = send.Start

	for l, f := range test {
		statech := rec.LinkStateCh()

		log.Debugf("Test: Set to state:%s", l)

		time.Sleep(1 * time.Millisecond)
		f()

		select {
		case ns := <-statech:
			if ns.State() != l {
				t.Errorf("State %s != %s ", ns.State(), l)
			}
		case <-time.After(time.Second):
			t.Errorf("Transition time out on %s", l)
			break
		}
	}

	//	rec.AddLinkDownCh(send)
	//	rec.AddLinkNoticeCh(send)
	//	rec.AddLinkNoticeStateCh(send)
	//	rec.AddLinkStateCh(send)
	//	rec.AddLinkUpCh(send)
	//	rec.AddLinkUpDownCh(send)

	//	rec.AddLinkCloseCh(send)
	//	rec.AddLinkLatencyCh(send)
	//	rec.AddLinkLossCh(send)
	//	rec.AddLinkSaturationCh(send)

	//	send.Close()
	//	if rec.GetState() != State

	cx.Cancel()
}
