package packet

import "testing"

func TestAll(t *testing.T) {
	var zerosig PacketSigType = PacketSigType(0x0000)
	var sig PacketSigType

	sig = VERSION_1
	if !sig.V1(){
			t.Fatal("Version 1")
	}
	if zerosig.V1(){
			t.Fatal("Zero  Version 1")
	}
	
}