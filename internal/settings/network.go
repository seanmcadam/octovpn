package settings

import (
	"fmt"
	"log"
	"strconv"

	"github.com/seanmcadam/octovpn/octolib/netlib"
)

type NetworkStruct struct {
	Name  string `json:"name"`
	Proto string `json:"proto"`
	Host  string `json:"host"`
	Port  string `json:"port"`
	Mtu   string `json:"mtu,omitempty"`
	Auth  string `json:"auth,omitempty"`
	port  *netlib.NetworkPort
	mtu   *uint16
}

type NetworkConfig struct {
	Server []NetworkStruct `json:"server,omitempty"`
	Client []NetworkStruct `json:"client,omitempty"`
}

func (ns *NetworkStruct) Validate() (err error) {

	switch ns.Proto {
	case "tcp":
	case "udp":
		if !(netlib.ValidIPHost(ns.Host) || netlib.ValidIP(ns.Host)) {
			return fmt.Errorf("NetworkStruct Invalid IP Address:%s", ns.Host)
		}
	case "tcp4":
	case "udp4":
		if !(netlib.ValidIP4Host(ns.Host) || netlib.ValidIP4(ns.Host)) {
			return fmt.Errorf("NetworkStruct Invalid IP4 Address:%s", ns.Host)
		}
	case "tcp6":
	case "udp6":
		if !(netlib.ValidIP6Host(ns.Host) || netlib.ValidIP6(ns.Host)) {
			return fmt.Errorf("NetworkStruct Invalid IP6 Address:%s", ns.Host)
		}
	default:
		return fmt.Errorf("NetworkStruct Proto:%s", ns.Proto)
	}

	err = ns.validatePort()
	return err
}

func (ns *NetworkStruct) validatePort() (err error) {
	_, err = strconv.ParseUint(ns.Port, 10, 16)
	return err
}

func (ns *NetworkStruct) GetHost() (h string) {
	return ns.Host
}

func (ns *NetworkStruct) GetPort() (p netlib.NetworkPort) {
	if ns.port == nil {
		val, err := strconv.ParseUint(ns.Port, 10, 16)
		if err != nil {
			log.Fatalf("NetoworkStruct bad Port val %s", ns.Port)
		}
		port := netlib.NetworkPort(uint16(val))
		ns.port = &port
	}
	return *ns.port
}

func (ns *NetworkStruct) GetMtu() uint16 {
	if ns.mtu == nil {
		if ns.Mtu == "" {
			ns.Mtu = fmt.Sprintf("%d", netlib.DefaultMaxPacketSize)
		}
		val, err := strconv.ParseUint(ns.Mtu, 10, 16)
		if err != nil {
			log.Fatalf("NetoworkStruct bad Mtu val %s", ns.Mtu)
		}
		mtu := uint16(val)
		ns.mtu = &mtu
	}
	return *ns.mtu
}
