package interfaces

import (
	"github.com/seanmcadam/bufferpool"
	"github.com/seanmcadam/octovpn/common"
)

type LayerInterface interface {
	Send(b *bufferpool.Buffer)
	Reset() // Closes or restart the layer
	RecvCh() chan *bufferpool.Buffer
	StatusCh() chan common.LayerStatus
	//MtuCh() chan common.MtuType
}
