package link

import (
	"testing"
	"time"
)

//func (ls *LinkStateStruct) ToggleState(s LinkStateType) {
//func (ls *LinkStateStruct) StateToggleCh() (newch chan LinkStateType) {
//func (ls *LinkStateStruct) GetState() LinkStateType {

func TestLinkState_compile(t *testing.T) {
	_ = NewLinkState()
}

func TestLinkState_StateToggles(t *testing.T) {
	ls := NewLinkState()

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

}

func TestLinkState_StateTogglesChannel(t *testing.T) {
	ls := NewLinkState()

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

}
