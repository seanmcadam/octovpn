package packet

import (
	"testing"
)

func TestNewIPv6_compile(t *testing.T) {
	_, _ = NewIPv6()
}

func TestNewIPv6_nil_methods(t *testing.T) {

	var i *IPv6Packet

	i.Size()
	i.ToByte()
	i, _ = MakeIPv6([]byte{})
	i.Size()
	i.ToByte()

}
