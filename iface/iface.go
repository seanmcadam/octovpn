package iface

import (
	"fmt"

	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/octoconfig"
	"github.com/seanmcadam/octovpn/octolib"
	"github.com/seanmcadam/octovpn/packet"
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
	readFrame  chan *packet.EthFrame
	writeFrame chan *packet.EthFrame
}

func NewIface(cx *ctx.Ctx, conf *octoconfig.ConfigInterface) (iface *IFace, e error) {

	cx = cx.NewWithCancel()

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
		var l netlink.Link
		l, e = netlink.LinkByName(conf.Name)
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
				readFrame:  make(chan *packet.EthFrame, readChanDepth),
				writeFrame: make(chan *packet.EthFrame, writeChanDepth),
			}

			iface.addIP()
		}
	}

	return iface, e
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
	i.ctx.Cancel()
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
// Blocking read from the readFrame channel
//
func (i *IFace) ReadEthChan() <-chan *packet.EthFrame {
	// i.ctx.LogLocation()
	return i.readFrame
}

//
// Write()
// Blocking write to the writeFrame channel
//
func (i *IFace) Write(eth *packet.EthFrame) {
	i.ctx.LogLocation()
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
		case packet.ET_IPv4:
			//payload := eth.Payload()
			//i.ctx.Logf(ctx.LogLevelTrace, "Eth\n\tSource:%s Dest:%s\n\t%s Source:%s Dest:%s",
			//	eth.Source(),
			//	eth.Destination(),
			//	packet.IPv4Frame(payload).Protocol().String(),
			//	packet.IPv4Frame(payload).Source(),
			//	packet.IPv4Frame(payload).Dest(),
			//)
			i.readFrame <- eth
		case packet.ET_ARP:
			i.ctx.Logf(ctx.LogLevelTrace, "Frame Source:%s Dest:%s, Type:%s", eth.Source(), eth.Destination(), eth.Ethertype().String())
			i.readFrame <- eth
		case packet.ET_IPv6:
			// i.ctx.Logf(ctx.LogLevelTrace, "DROP IPv6 Source:%s Dest:%s, Type:%s", eth.Source(), eth.Destination(), eth.Ethertype().String())
		default:
			i.ctx.Logf(ctx.LogLevelTrace, "DROP frame Source:%s Dest:%s, Type:%s", eth.Source(), eth.Destination(), eth.Ethertype().String())
			// Drop it
		}
	}
}

//
// goWriter()
//
func (i *IFace) goWriter() {
	for {
		select {
		case <-i.ctx.DoneChan():
			return
		case frame := <-i.writeFrame:
			e := i.write(frame)
			if e != nil {
				i.ctx.Logf(ctx.LogLevelError, "write() error:%s", e)
				return
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
func (i *IFace) read() (eth *packet.EthFrame, e error) {
	var f packet.EthFrame
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
func (i *IFace) write(eth *packet.EthFrame) (e error) {
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
