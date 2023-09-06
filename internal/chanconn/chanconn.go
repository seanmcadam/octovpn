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
	width   packet.PacketWidth
	recvch  chan *packet.PacketStruct
	pinger  pinger.PingerStruct
	counter counter.CounterStruct
}

func (c *ConnStruct) MaxLocalMtu() (size packet.PacketSizeType) {
	size = packet.PacketSigSize + packet.PacketSize16Size
	if c.width == packet.PacketWidth32 {
		size += packet.PacketCounter32Size
		size += packet.PacketPing32Size
	if c.width == packet.PacketWidth64 {
		size += packet.PacketCounter64Size
		size += packet.PacketPing64Size
	} else {
		log.FatalfStack("CannedStruct:%v", c)
	}
	return size
}

func NewConn32(ctx *ctx.Ctx, config *settings.NetworkStruct, confFunc NewConnFunc) (ci interfaces.ChannelInterface, err error) {

	conn, err := confFunc(ctx, config)
	if err != nil {
		return nil, err
	}

	cs := &ChanconnStruct{
		cx:      ctx,
		width:   packet.PacketWidth32,
		conn:    conn,
		recvch:  make(chan *packet.PacketStruct, 16),
		pinger:  pinger.NewPinger32(ctx, PingFreq, PingTimeout),
		counter: counter.NewCounter32(ctx),
	}

	go cs.goRecv()

	return cs, err
}


func NewConn64(ctx *ctx.Ctx, config *settings.NetworkStruct, confFunc NewConnFunc) (ci interfaces.ChannelInterface, err error) {

	conn, err := confFunc(ctx, config)
	if err != nil {
		return nil, err
	}

	cs := &ChanconnStruct{
		cx:      ctx,
		width:   packet.PacketWidth64,
		conn:    conn,
		recvch:  make(chan *packet.PacketStruct, 16),
		pinger:  pinger.NewPinger64(ctx, PingFreq, PingTimeout),
		counter: counter.NewCounter64(ctx),
	}

	go cs.goRecv()

	return cs, err
}
