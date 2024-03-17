package common

type LayerStatus uint8

const (
	LayerClosed LayerStatus = 0
	LayerDown   LayerStatus = 1
	LayerUp     LayerStatus = 2
)