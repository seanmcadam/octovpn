package iface

import (
	"log"

	"github.com/songgao/packets/ethernet"
	"github.com/songgao/water"
)

func OpenTap(octoconfig.ConfigInterface conf) (ifce *water.Interface, e error) {

	name := conf.Name
	ip :=  conf.IP
	netmask := conf.NetMask
	mtu := conf.MTU

	config := water.Config{
		DeviceType: water.TAP,
		Name: name,
	}

	ifce, e = water.New(config)

	// Set IP / NetMask / MTU
	

	return ifce, e
}



func SomeFunction(ifce *water.Interface) {

	var frame ethernet.Frame
	for {
		frame.Resize(1500)
		n, err := ifce.Read([]byte(frame))
		if err != nil {
			log.Fatal(err)
		}
		frame = frame[:n]
		log.Printf("Dst: %s\n", frame.Destination())
		log.Printf("Src: %s\n", frame.Source())
		log.Printf("Ethertype: % x\n", frame.Ethertype())
		log.Printf("Payload: % x\n", frame.Payload())
	}
}
