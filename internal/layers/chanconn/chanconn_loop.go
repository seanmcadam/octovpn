package chanconn

import (
	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/layers/conn/tcpcli"
	"github.com/seanmcadam/octovpn/internal/layers/conn/tcpsrv"
	"github.com/seanmcadam/octovpn/internal/layers/conn/udpcli"
	"github.com/seanmcadam/octovpn/internal/layers/conn/udpsrv"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/netlib"
)

func NewTcp32ConnLoop(ctx *ctx.Ctx) (loop1 interfaces.ConnInterface, loop2 interfaces.ConnInterface, err error) {

	config := &settings.ConnectionStruct{
		Name:  "Server Loopback TCP",
		Proto: "tcp",
		Host:  "127.0.0.1",
		Port:  settings.ConfigPortType(uint16(netlib.GetRandomNetworkPort())),
		Auth:  "",
	}

	loop1, err = NewConn32(ctx, config, tcpsrv.New)
	if err != nil {
		return nil, nil, err
	}

	config.Name = "Client Loopback TCP"
	loop2, err = NewConn32(ctx, config, tcpcli.New)
	if err != nil {
		return nil, nil, err
	}

	return loop1, loop2, err
}
func NewTcp64ConnLoop(ctx *ctx.Ctx) (loop1 interfaces.ConnInterface, loop2 interfaces.ConnInterface, err error) {

	config := &settings.ConnectionStruct{
		Name:  "LoopbackTCP",
		Proto: "tcp",
		Host:  "127.0.0.1",
		Port:  settings.ConfigPortType(uint16(netlib.GetRandomNetworkPort())),
		Auth:  "",
	}

	loop1, err = NewConn64(ctx, config, tcpsrv.New)
	if err != nil {
		return nil, nil, err
	}

	loop2, err = NewConn64(ctx, config, tcpcli.New)
	if err != nil {
		return nil, nil, err
	}

	return loop1, loop2, err
}

// NewUdpLoop()
// Returns a pair of connected UDP sockets as a
func NewUdp32ConnLoop(ctx *ctx.Ctx) (srv interfaces.ConnInterface, cli interfaces.ConnInterface, err error) {

	udpconfig := &settings.ConnectionStruct{
		Name:  "LoopbackUDP",
		Proto: "udp",
		Host:  "127.0.0.1",
		Port:  settings.ConfigPortType(uint16(netlib.GetRandomNetworkPort())),
		Auth:  "",
	}

	srv, err = NewConn32(ctx, udpconfig, udpsrv.New)
	if err != nil {
		return nil, nil, err
	}

	cli, err = NewConn32(ctx, udpconfig, udpcli.New)
	if err != nil {
		return nil, nil, err
	}

	return srv, cli, err
}

func NewUdp64ConnLoop(ctx *ctx.Ctx) (srv interfaces.ConnInterface, cli interfaces.ConnInterface, err error) {

	udpconfig := &settings.ConnectionStruct{
		Name:  "LoopbackUDP",
		Proto: "udp",
		Host:  "127.0.0.1",
		Port:  settings.ConfigPortType(uint16(netlib.GetRandomNetworkPort())),
		Auth:  "",
	}

	srv, err = NewConn64(ctx, udpconfig, udpsrv.New)
	if err != nil {
		return nil, nil, err
	}

	cli, err = NewConn64(ctx, udpconfig, udpcli.New)
	if err != nil {
		return nil, nil, err
	}

	return srv, cli, err
}
