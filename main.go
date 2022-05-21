package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/iface"
	"github.com/seanmcadam/octovpn/octoconfig"
)

func main() {

	cx := ctx.NewContext()

	configs, e := octoconfig.ReadConfigs()
	if e != nil {
		cx.LogPanic(e)
	}

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

	Iface, e := iface.NewIface(cx, configs.Iface)
	if e != nil {
		cx.Logf(ctx.LogLevelPanic, "NewIface() error:%s", e)
	}

	// Set up IFace reader routine

	//
	// Setup Interface
	// IP Address / Netmask / MTU
	// Open Connections
	//

	if configs.Conn != nil {

	}

	//
	// Setup Listeners
	//

	if configs.List != nil {

	}

	//
	// Launch traffic manager
	//

	//
	// Run Packet Loop
	//

	Iface.Start()
	defer Iface.Stop()

	<-cx.Done()
}
