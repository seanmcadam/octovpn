package tracker

import (
	"testing"
)

func TestNewCompile(t *testing.T) {

	closech := make(chan interface{})
	_ = NewTracker(closech)
	close(closech)

}
