package iface

import (
	"github.com/songgao/packets/ethernet"
	"github.com/songgao/water"
)

type DeviceType string

const (
	TAP DeviceType = "TAP"
	TUN DeviceType = "TUN"
)

type IFace struct {
	ctx     *ctx
	device  DeviceType
	iface   *water.Interface
	name    string
	ip      string
	netmask string
	mtu     int
}

type Ethernet struct {
	frame ethernet.Frame
}

func NewIFace(ctx ctx, device DeviceType, name string, ip string, netmask string, mtu int) (iface *IFace, e error) {

	ctx.Logf(ctx.LogTrace, " called ")

	var dev water.DeviceType
	switch device {
	case TAP:
		dev = water.TAP
	case TUN:
		dev = water.TUN
	default:
		ctx.Logf(ctx.LogPanic, " default reached ")
	}

	_ = water.DevicePermissions{
		Owner: 0, //-1 default, but value is uint
		Group: 0, //-1,
	}

	psp := water.PlatformSpecificParams{
		Name:        name,
		Persist:     false, //(bool)
		Permissions: nil,   // (*DevicePermissions)
		MultiQueue:  false, // (bool)
	}

	config := water.Config{
		DeviceType:             dev,
		PlatformSpecificParams: psp,
	}

	// Validate name,ip,netmask, mtu

	i, e := water.New(config)
	iface = &IFace{
		ctx:     ctx,
		device:  device,
		iface:   i,
		name:    name,
		ip:      ip,
		netmask: netmask,
		mtu:     mtu,
	}

	return iface, e
}

//
// Name()
// return the device name (in case it was OS assigned)
//
func (i *IFace) Name() (name string) {
	return i.iface.Name()
}
