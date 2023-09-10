package packet

import (
	"testing"
)

func TestNewEth_compile(t *testing.T) {
	_, _ = NewEth()
}

func TestNewEth_nil_methods(t *testing.T) {

	var e *EthPacket

	e.Size()
	e.ToByte()
	e, _ = MakeEth([]byte{})
	e.Size()
	e.ToByte()

}
