package status

import (
	"github.com/seanmcadam/ctx"
	"github.com/seanmcadam/loggy"
)

type LayerStatus string

//
//
//

type LayerStatusStruct struct {
	cx     *ctx.Ctx
	status LayerStatus
	ch     chan LayerStatus
}

const (
	LayerStatusInit   LayerStatus = "Init"
	LayerStatusClosed LayerStatus = "Closed"
	LayerStatusDown   LayerStatus = "Down"
	LayerStatusUp     LayerStatus = "Up"
)

func New(cx *ctx.Ctx) (ls *LayerStatusStruct) {
	ls = &LayerStatusStruct{
		cx:     cx,
		status: LayerStatusInit,
		ch:     make(chan LayerStatus, 1),
	}

	return ls
}

func (lss *LayerStatusStruct) GetCh() chan LayerStatus {
	if lss == nil {
		loggy.FatalfStack("nil pointer")
	}

	return lss.ch
}

func (lss *LayerStatusStruct) Get() LayerStatus {
	if lss == nil {
		loggy.FatalfStack("nil pointer")
	}

	if !lss.cx.Done() {
		return lss.status
	}
	return LayerStatusClosed
}

func (lss *LayerStatusStruct) Set(s LayerStatus) {
	if lss == nil {
		loggy.FatalfStack("nil pointer")
	}

	if !lss.cx.Done() {
		if s != lss.status {
			lss.status = s
			select {
			case lss.ch <- lss.status:
			default:
				loggy.Warnf("Status Chan is backed up, waiting...")
				go func() {
					lss.ch <- lss.status
				}()
			}
		}
	}
}

func (lss *LayerStatusStruct) goRun() {
	select {
	case <-lss.cx.DoneChan():
		close(lss.ch)
	}
}

func (ls LayerStatus) IsDown() (s bool) {
	if ls != LayerStatusUp {
		s = true
	}
	return s
}

func (ls LayerStatus) IsUp() (s bool) {
	if ls == LayerStatusUp {
		s = true
	}
	return s
}

func (ls LayerStatus) String() (s string) {
	return string(ls)
}
