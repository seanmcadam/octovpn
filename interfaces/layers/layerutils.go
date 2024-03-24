package layers

import (
	"github.com/seanmcadam/bufferpool"
	"github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/common/status"
)

type layerStruct struct {
	Pair   *layerStruct
	recvch chan *bufferpool.Buffer
	status status.LayerStatus
}

func CreateLayerPair() (a LayerInterface, b LayerInterface) {

	la := &layerStruct{
		recvch: make(chan *bufferpool.Buffer, 1),
		status: status.LayerStatusUp,
	}
	lb := &layerStruct{
		recvch: make(chan *bufferpool.Buffer, 1),
		status: status.LayerStatusUp,
	}

	la.Pair = lb
	lb.Pair = la

	return LayerInterface(la), LayerInterface(lb)

}

func (l *layerStruct) Send(b *bufferpool.Buffer) {
	l.Pair.recvch <- b
}

func (l *layerStruct) RecvCh() chan *bufferpool.Buffer {
	return l.recvch

}
func (l *layerStruct) Reset() { loggy.FatalStack("Not implemented") }
func (l *layerStruct) Status() (status status.LayerStatus) {
	return l.status
}
func (l *layerStruct) StatusCh() (statusCh chan status.LayerStatus) {
	statusCh = make(chan status.LayerStatus)
	statusCh <- l.status
	return
}
