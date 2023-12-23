package packet

import (
	"testing"
)

func TestNewIPv4_compile(t *testing.T) {
	_, _ = NewIPv4()
}

func TestNewIPv4_nil_methods(t *testing.T) {

	var i *IPv4Packet

	i.Size()
	i.ToByte()
	i, _ = MakeIPv4([]byte{})
	i.Size()
	i.ToByte()

}
