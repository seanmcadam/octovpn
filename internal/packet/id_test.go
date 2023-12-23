package packet

import (
	"testing"
)

func TestNewID_compile(t *testing.T) {
	_, _ = NewID()
}

func TestNewID_nil_methods(t *testing.T) {

	var i *IDPacket

	i.Size()
	i.ToByte()
	i, _ = MakeID([]byte{})
	i.Size()
	i.ToByte()

}
