package loopconn

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

// NewUdpLoop()
// Returns a pair of connected UDP sockets as a
func NewUdpConnLoop(ctx *ctx.Ctx) (srv interfaces.ConnInterface, cli interfaces.ConnInterface, err error) {

	udpconfig := &settings.ConnectionStruct{
		Name:  "LoopbackUDP",
		Proto: "udp",
		Host:  "127.0.0.1",
		Port:  settings.ConfigPortType(uint16(netlib.GetRandomNetworkPort())),
		Auth:  "",
	}

	srv, err = udpsrv.New(ctx, udpconfig)
	if err != nil {
		return nil, nil, err
	}

	cli, err = udpcli.New(ctx, udpconfig)
	if err != nil {
		return nil, nil, err
	}

	return srv, cli, err
}

func NewTcpConnLoop(ctx *ctx.Ctx) (loop1 interfaces.ConnInterface, loop2 interfaces.ConnInterface, err error) {

	config := &settings.ConnectionStruct{
		Name:  "LoopbackTCP",
		Proto: "tcp",
		Host:  "127.0.0.1",
		Port:  settings.ConfigPortType(uint16(netlib.GetRandomNetworkPort())),
		Auth:  "",
	}

	loop1, err = tcpsrv.New(ctx, config)
	if err != nil {
		return nil, nil, err
	}

	loop2, err = tcpcli.New(ctx, config)
	if err != nil {
		return nil, nil, err
	}

	return loop1, loop2, err
}
