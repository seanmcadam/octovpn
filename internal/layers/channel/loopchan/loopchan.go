package loopchan

import (
	"log"

	"github.com/seanmcadam/octovpn/internal/layers/chanconn"
	"github.com/seanmcadam/octovpn/internal/layers/channel"
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
func NewUdpChanLoop(ctx *ctx.Ctx) (srv *channel.ChannelStruct, cli *channel.ChannelStruct, err error) {

	udpconfig := &settings.ConnectionStruct{
		Name:  "LoopbackUDP",
		Proto: "udp",
		Host:  "127.0.0.1",
		Port:  settings.ConfigPortType(uint16(netlib.GetRandomNetworkPort())),
		Auth:  "",
	}

	connsrv, err := chanconn.NewConn32(ctx, udpconfig, udpsrv.New)
	if err != nil {
		log.Fatalf("NewConn() Err:%s", err)
	}

	conncli, err := chanconn.NewConn32(ctx, udpconfig, udpcli.New)
	if err != nil {
		log.Fatalf("NewConn() Err:%s", err)
	}

	srv, err = channel.NewChannel32(ctx, connsrv)
	if err != nil {
		log.Fatalf("NewUdpChanLoop Err:%s", err)
	}

	cli, err = channel.NewChannel32(ctx, conncli)
	if err != nil {
		log.Fatalf("NewUdpChanLoop Err:%s", err)
	}

	return srv, cli, err
}

func NewTcpChanLoop(ctx *ctx.Ctx) (srv *channel.ChannelStruct, cli *channel.ChannelStruct, err error) {

	tcpconfig := &settings.ConnectionStruct{
		Name:  "LoopbackTCP",
		Proto: "tcp",
		Host:  "127.0.0.1",
		Port:  settings.ConfigPortType(uint16(netlib.GetRandomNetworkPort())),
		Auth:  "",
	}

	connsrv, err := chanconn.NewConn32(ctx, tcpconfig, tcpsrv.New)
	if err != nil {
		log.Fatalf("NewConn() Err:%s", err)
	}

	conncli, err := chanconn.NewConn32(ctx, tcpconfig, tcpcli.New)
	if err != nil {
		log.Fatalf("NewConn() Err:%s", err)
	}

	srv, err = channel.NewChannel32(ctx, connsrv)
	if err != nil {
		log.Fatalf("NewUdpChanLoop Err:%s", err)
	}

	cli, err = channel.NewChannel32(ctx, conncli)
	if err != nil {
		log.Fatalf("NewUdpChanLoop Err:%s", err)
	}

	return srv, cli, err
}
