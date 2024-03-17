package interfaces

import "github.com/seanmcadam/bufferpool"

type LayerStatus uint8

const (
	Closed LayerStatus = 0
	Down   LayerStatus = 1
	Up     LayerStatus = 2
)

type LayerInterface interface {
	Send(b *bufferpool.Buffer)
	RecvCh() chan *bufferpool.Buffer
	Reset()
	StatusCh() chan LayerStatus
}
