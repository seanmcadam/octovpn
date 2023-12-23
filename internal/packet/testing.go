package packet

import (
	"fmt"

	"github.com/seanmcadam/counter"
)

// -
//
// -
func TestConn32Packet() (*PacketStruct, error) {
	return NewPacket(SIG_CONN_32_RAW, []byte("conn 32 testpacket"), counter.NewCount(uint32(1)))
}

// -
//
// -
func TestChan32Packet() (*PacketStruct, error) {
	return NewPacket(SIG_CHAN_32_RAW, []byte("chan 32 testpacket"), counter.NewCount(uint32(1)))
}

// -
//
// -
func TestSite32Packet() (*PacketStruct, error) {
	return NewPacket(SIG_SITE_32_RAW, []byte("site 32 testpacket"), counter.NewCount(uint32(1)))
}

// -
//
// -
func TestConn64Packet() (*PacketStruct, error) {
	return NewPacket(SIG_CONN_64_RAW, []byte("conn 64 testpacket"), counter.NewCount(uint64(1)))
}

// -
//
// -
func TestChan64Packet() (*PacketStruct, error) {
	return NewPacket(SIG_CHAN_64_RAW, []byte("chan 64 testpacket"), counter.NewCount(uint64(1)))
}

// -
//
// -
func TestSite64Packet() (*PacketStruct, error) {
	return NewPacket(SIG_SITE_64_RAW, []byte("site 64 testpacket"), counter.NewCount(uint64(1)))
}

// -
//
// -
func TestRouterPacket() (*PacketStruct, error) {
	return NewPacket(SIG_ROUTE_RAW, []byte("testpacket"))
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
