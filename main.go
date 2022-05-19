package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/seanmcadam/octovpn/iface"
	"github.com/seanmcadam/octovpn/octoconfig"
)

func main() {

	configs := octoconfig.ReadConfigs()
	_ = configs

	fmt.Printf("Running...\n")

	//
	// Verify permissions (am I root?)
	//

	user, e := user.Current()

	fmt.Printf("UID:%s\n", user.Uid)

	if e != nil {
		panic("")
	}
	if user.Uid != "0" {
		os.Exit(-1)
	}

	//
	// Open Tun/Tap
	// Type /

	iface.OpenTap(configs.IFace)

	//
	// Setup Interface
	// IP Address / Netmask / MTU
	//

	//
	// Open Connections
	//

	//
	// Launch traffic manager
	//

	//
	// Run Packet Loop
	//

}
