package layers

import (
	"github.com/seanmcadam/bufferpool"
	"github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/common"
)

type layerStruct struct {
	Pair   *layerStruct
	recvch chan *bufferpool.Buffer
}

func CreateLayerPair() (a LayerInterface, b LayerInterface) {

	la := &layerStruct{
		recvch: make(chan *bufferpool.Buffer, 1),
	}
	lb := &layerStruct{
		recvch: make(chan *bufferpool.Buffer, 1),
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
func (l *layerStruct) Reset()                            { loggy.FatalStack("Not implemented") }
func (l *layerStruct) StatusCh() chan common.LayerStatus { loggy.FatalStack("Not implemented"); return nil }
