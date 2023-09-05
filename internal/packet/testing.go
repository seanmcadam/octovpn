package packet

import (
	"fmt"

	"github.com/seanmcadam/octovpn/internal/counter"
)

// -
//
// -
func Testpacket() (*PacketStruct, error) {
	return NewPacket(SIG_CONN_32_RAW, []byte("testpacket"), counter.MakeCounter32(1))
}

// -
//
// -
func Validatepackets(sent *PacketStruct, recv *PacketStruct) (err error) {

	if sent.Raw() == nil {
		panic("nil sent packet")
	}

	if recv.Raw() == nil {
		return fmt.Errorf("Recv Nil packet")
	}

	if string(sent.Raw()) != string(recv.Raw()) {
		return fmt.Errorf("Packets do not match '%s' != '%s'", sent.Raw(), recv.Raw())
	}

	return err
}
