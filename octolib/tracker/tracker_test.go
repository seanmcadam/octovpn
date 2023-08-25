package tracker

import (
	"testing"
)

func TestNewTracker(t *testing.T) {

	closech := make(chan interface{})
	_ = NewTracker(closech)

	close(closech)

}
