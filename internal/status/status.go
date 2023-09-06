package status

import "fmt"

type StatusType uint8

type StatusLatency struct {
}

type StatusLoss struct {
}

type StatusStruct struct {
	status  StatusType
	linkup  bool
	latency *StatusLatency
	loss    *StatusLoss
}

const (
	statusLinkChange    StatusType = 0x01
	statusLatencyChange StatusType = 0x02
	statusLossChange    StatusType = 0x04
	statusError         StatusType = 0x80
)

func (s *StatusStruct) LinkChange() bool {
	return (s.status & statusLinkChange) > 0
}
func (s *StatusStruct) LatencyChange() bool {
	return (s.status & statusLatencyChange) > 0
}
func (s *StatusStruct) LossChange() bool {
	return (s.status & statusLossChange) > 0
}
func (s *StatusStruct) Error() bool {
	return (s.status & statusError) > 0
}

func (s *StatusStruct) LinkUp() bool {
	return s.linkup
}

func (s *StatusStruct) Latency() (l *StatusLatency, err error) {
	if !s.LatencyChange() {
		return nil, fmt.Errorf("No Latency change")
	}
	return s.latency, err
}
func (s *StatusStruct) Loss() (l *StatusLoss, err error) {
	if !s.LossChange() {
		return nil, fmt.Errorf("No Loss change")
	}
	return s.loss, err
}
