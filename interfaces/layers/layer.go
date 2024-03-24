package layers

import (
	"github.com/seanmcadam/bufferpool"
	"github.com/seanmcadam/octovpn/common/status"
)

type LayerInterface interface {
	Send(b *bufferpool.Buffer)
	Reset() // Closes or restart the layer
	RecvCh() chan *bufferpool.Buffer
	Status() status.LayerStatus
	StatusCh() chan status.LayerStatus
	//MtuCh() chan common.MtuType
}
