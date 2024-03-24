package status

import (
	"testing"

	"github.com/seanmcadam/ctx"
)

func TestCompile(t *testing.T) {}

func TestNew(t *testing.T) {

	cx := ctx.New()
	status := New(cx)

	if LayerStatusInit != status.Get() {
		t.Errorf("wrong status")
	}

	select {
	case <-status.GetCh():
		t.Errorf("returned status on chan")
	default:
	}

	checkstatus := func(l LayerStatus) {
		var s LayerStatus

		current := status.Get()
		status.Set(l)
		s = status.Get()
		if s != l {
			t.Errorf("wrong status: Set:%s Get:%s", l, s)
		}

		if current != l {
			select {
			case s = <-status.GetCh():
				if s != l {
					t.Errorf("wrong status: Set:%s GetCh:%s", l, s)
				}
			default:
				t.Errorf("did not return status Set:%s", s)
			}
		}
	}

	for _, i := range []LayerStatus{
		LayerStatusInit,
		LayerStatusClosed,
		LayerStatusDown,
		LayerStatusUp,
	} {
		checkstatus(i)
	}

	status.Set(LayerStatusUp)
	if LayerStatusUp != status.Get() {
		t.Errorf("wrong status")
	}

}
