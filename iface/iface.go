package iface

import (
	"errors"
	"fmt"

	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/octoconfig"
	"github.com/seanmcadam/octovpn/octolib"
	"github.com/vishvananda/netlink"

	"github.com/songgao/water"
)

type DeviceType string

const readChanDepth = 10
const writeChanDepth = 10

const (
	TAP DeviceType = "TAP"
	TUN DeviceType = "TUN"
)

type IFace struct {
	ctx        *ctx.Ctx
	device     DeviceType
	iface      *water.Interface // Used to create an interface
	link       netlink.Link
	name       string
	ip         string
	netmask    string
	mtu        int
	readFrame  chan *Frame
	writeFrame chan *Frame
}

func NewIface(cx *ctx.Ctx, conf *octoconfig.ConfigInterface) (iface *IFace, e error) {

	cx.Logf(ctx.LogLevelTrace, " called ")

	var dev DeviceType
	var waterdev water.DeviceType
	switch conf.TunTap {
	case string(TAP):
		dev = TAP
		waterdev = water.TAP
	case string(TUN):
		dev = TUN
		waterdev = water.TUN
	default:
		cx.Logf(ctx.LogLevelPanic, " default reached ")
	}

	_ = water.DevicePermissions{
		Owner: 0, //-1 default, but value is uint
		Group: 0, //-1,
	}

	psp := water.PlatformSpecificParams{
		Name:        conf.Name,
		Persist:     false, //(bool)
		Permissions: nil,   // (*DevicePermissions)
		MultiQueue:  false, // (bool)
	}

	config := water.Config{
		DeviceType:             waterdev,
		PlatformSpecificParams: psp,
	}

	// Validate name,ip,netmask, mtu

	i, e := water.New(config)

	if e == nil {
		l, e := netlink.LinkByName(conf.Name)
		if e == nil {

			iface = &IFace{
				ctx:        cx,
				device:     dev,
				iface:      i,
				link:       l,
				name:       conf.Name,
				ip:         conf.IP,
				netmask:    conf.Netmask,
				mtu:        conf.MTU,
				readFrame:  make(chan *Frame, readChanDepth),
				writeFrame: make(chan *Frame, writeChanDepth),
			}

			iface.addIP()
		}
	}

	return iface, e
}

//
// Ctx()
//
func (i *IFace) Ctx() *ctx.Ctx {
	return i.ctx
}

//
//
//
func (i *IFace) Start() {

	e := i.LinkUp()
	if e != nil {
		i.ctx.Logf(ctx.LogLevelPanic, "LinkUp() error:%s", e)
	}

	go i.goReader()
	go i.goWriter()
}

//
//
//
func (i *IFace) Stop() {
	e := i.LinkDown()
	if e != nil {
		i.ctx.Logf(ctx.LogLevelPanic, "LinkDown() error:%s", e)
	}
	i.Ctx().Cancel()
}

//
// addIP() adds the ip address to the link (tun or tap)
//
func (i *IFace) addIP() {

	bits, e := octolib.IP4netmask2netbits(i.netmask)
	if e != nil {
		i.ctx.Logf(ctx.LogLevelPanic, "IP4netmask2netbits() error:%s", e)
	}
	mask := fmt.Sprintf("/%s", bits)

	addrstr := i.ip + mask
	addr, e := netlink.ParseAddr(addrstr)
	if e != nil {
		i.ctx.Logf(ctx.LogLevelPanic, " ParseAddr() failed:%s", e)
	}

	//validate addr, or the componants

	netlink.AddrAdd(i.link, addr)
}

//
// Read()
//
func (i *IFace) Read() (eth *Frame, e error) {
	select {
	case <-i.ctx.Done():
		e := errors.New("Interface Closed")
		return eth, e
	case eth := <-i.readFrame:
		return eth, e
	}
}

//
// Write()
//
func (i *IFace) Write(eth *Frame) {
	i.writeFrame <- eth
}

//
// goReader()
//
func (i *IFace) goReader() {
	for {
		eth, e := i.read()
		if e != nil {
			i.ctx.Logf(ctx.LogLevelError, "read() error:%s", e)
			break
		}

		et := eth.Ethertype()
		switch et {
		case IPv4:
			fallthrough
		case IPv6:
			fallthrough
		case ARP:
			i.ctx.Logf(ctx.LogLevelTrace, "frame Source:%s Dest:%s, Type:%s", eth.Source(), eth.Destination(), eth.Ethertype().String())
		// i.readFrame <- eth
		default:
			// Drop it
		}
		// i.readFrame <- eth
	}
}

//
// goWriter()
//
func (i *IFace) goWriter() {
	for {
		select {
		case <-i.ctx.Done():
			break
		case frame := <-i.writeFrame:
			e := i.write(frame)
			if e != nil {
				i.ctx.Logf(ctx.LogLevelError, "write() error:%s", e)
				break
			}
		}
	}
}

//
// Name()
// return the device name (in case it was OS assigned)
//
func (i *IFace) Name() (name string) {
	return i.iface.Name()
}

//
// SetLinkUp() adds the ip address to the link (tun or tap)
//
func (i *IFace) LinkUp() error {
	return netlink.LinkSetUp(i.link)
}

//
// SetLinkDown() adds the ip address to the link (tun or tap)
//
func (i *IFace) LinkDown() error {
	return netlink.LinkSetDown(i.link)
}

//
// read()
//
func (i *IFace) read() (eth *Frame, e error) {
	var f Frame
	eth = &f
	eth.ResizePayload(1500)
	count, e := i.iface.Read([]byte(*eth))
	if e != nil {
		i.ctx.Logf(ctx.LogLevelPanic, "iface.Read() error:%s", e)
	}
	eth.ResizePayload(count)
	return eth, e
}

//
// write()
//
func (i *IFace) write(eth *Frame) (e error) {
	len := len([]byte(*eth))
	sentlen, e := i.iface.Write([]byte(*eth))
	if e != nil {
		i.ctx.Logf(ctx.LogLevelPanic, "iface.Write() error:%s", e)
	}
	if len != sentlen {
		i.ctx.Logf(ctx.LogLevelPanic, "iface.Write() lenths do not match:%d != %d", len, sentlen)
	}
	return e
}
