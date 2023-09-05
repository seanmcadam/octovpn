package chanconn

import (
	"time"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/internal/pinger"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
)

const PingFreq = 1 * time.Second
const PingTimeout = 2 * time.Second

type NewConnFunc func(*ctx.Ctx, *settings.NetworkStruct) (interfaces.ConnInterface, error)

type ChanconnStruct struct {
	cx      *ctx.Ctx
	conn    interfaces.ConnInterface
	recvch  chan *packet.PacketStruct
	pinger  pinger.PingerStruct
	counter counter.CounterStruct
}

func NewConn(ctx *ctx.Ctx, config *settings.NetworkStruct, confFunc NewConnFunc) (ci interfaces.ChannelInterface, err error) {

	conn, err := confFunc(ctx, config)
	if err != nil {
		return nil, err
	}

	cs := &ChanconnStruct{
		cx:      ctx,
		conn:    conn,
		recvch:  make(chan *packet.PacketStruct, 16),
		pinger:  pinger.NewPinger32(ctx, PingFreq, PingTimeout),
		counter: counter.NewCounter32(ctx),
	}

	go cs.goRecv()

	return cs, err
}
