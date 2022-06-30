package transit

//
// Transit is used to bridge the local VPN interface to the various VPN paths leading to the other end of the VPN
//
//

// func New(cx *ctx.Ctx, iface *iface.IFace, route *routes.RouteStruct) (transit *TransitStruct, e error) {
// func (t *TransitStruct) Start() {
// func (t *TransitStruct) Stop() {
// func (t *TransitStruct) goRun() {

import (
	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/iface"
	"github.com/seanmcadam/octovpn/routes"
)

type TransitStruct struct {
	ctx   *ctx.Ctx
	iface *iface.IFace
	route *routes.RouteStruct
}

func New(cx *ctx.Ctx, iface *iface.IFace, route *routes.RouteStruct) (transit *TransitStruct, e error) {

	transit = &TransitStruct{
		ctx:   cx,
		iface: iface,
		route: route,
	}

	return transit, e
}

func (t *TransitStruct) Start() {
	t.iface.Start()
	t.route.Start()
	go t.goRun()
}

func (t *TransitStruct) Stop() {
	t.iface.Stop()
	t.route.Stop()
	t.ctx.Cancel()
}

//
// goRun()
// Handle Eth packets going back and forth between local internface and vpn connections
// Add IP filters
// Add IP Accounting data collection
func (t *TransitStruct) goRun() {
	for {
		select {
		case <-t.ctx.DoneChan():
			return

		//
		// Get Eth packet from the local interface
		// Filter it
		// Send it
		case packet := <-t.iface.ReadEthChan():
			// Add IP filtering HERE
			t.route.Write(packet)

		//
		// Get Eth packet from the other side of the VPN
		// Send it up the stack
		case packet := <-t.route.ReadEthChan():
			// Filtering HERE
			t.iface.Write(packet)
		}
	}
}
