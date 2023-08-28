package loopchan

import (
	"fmt"
	"log"

	"github.com/seanmcadam/octovpn/internal/chanconn"
	"github.com/seanmcadam/octovpn/internal/chanconn/tcpcli"
	"github.com/seanmcadam/octovpn/internal/chanconn/tcpsrv"
	"github.com/seanmcadam/octovpn/internal/chanconn/udpcli"
	"github.com/seanmcadam/octovpn/internal/chanconn/udpsrv"
	"github.com/seanmcadam/octovpn/internal/channel"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/netlib"
)

// NewUdpLoop()
// Returns a pair of connected UDP sockets as a
func NewUdpChanLoop(ctx *ctx.Ctx) (srv *channel.ChannelStruct, cli *channel.ChannelStruct, err error) {

	udpconfig := &settings.NetworkStruct{
		Name:  "LoopbackUDP",
		Proto: "udp",
		Host:  "127.0.0.1",
		Port:  fmt.Sprintf("%d", uint16(netlib.GetRandomNetworkPort())),
		Auth:  "",
	}

	connsrv, err := chanconn.NewConn(ctx, udpconfig, udpsrv.New)
	if err != nil {
		log.Fatalf("NewConn() Err:%s", err)
	}

	conncli, err := chanconn.NewConn(ctx, udpconfig, udpcli.New)
	if err != nil {
		log.Fatalf("NewConn() Err:%s", err)
	}

	srv, err = channel.NewChannel(ctx, connsrv)
	if err != nil {
		log.Fatalf("NewUdpChanLoop Err:%s", err)
	}

	cli, err = channel.NewChannel(ctx, conncli)
	if err != nil {
		log.Fatalf("NewUdpChanLoop Err:%s", err)
	}

	return srv, cli, err
}

func NewTcpChanLoop(ctx *ctx.Ctx) (srv *channel.ChannelStruct, cli *channel.ChannelStruct, err error) {

	tcpconfig := &settings.NetworkStruct{
		Name:  "LoopbackTCP",
		Proto: "tcp",
		Host:  "127.0.0.1",
		Port:  fmt.Sprintf("%d", uint16(netlib.GetRandomNetworkPort())),
		Auth:  "",
	}

	connsrv, err := chanconn.NewConn(ctx, tcpconfig, tcpsrv.New)
	if err != nil {
		log.Fatalf("NewConn() Err:%s", err)
	}

	conncli, err := chanconn.NewConn(ctx, tcpconfig, tcpcli.New)
	if err != nil {
		log.Fatalf("NewConn() Err:%s", err)
	}

	srv, err = channel.NewChannel(ctx, connsrv)
	if err != nil {
		log.Fatalf("NewUdpChanLoop Err:%s", err)
	}

	cli, err = channel.NewChannel(ctx, conncli)
	if err != nil {
		log.Fatalf("NewUdpChanLoop Err:%s", err)
	}

	return srv, cli, err
}
