package chanconn

import (
	"time"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/auth"
	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/internal/link"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/internal/pinger"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

const PingFreq = 1 * time.Second
const PingTimeout = 2 * time.Second

type NewConnFunc func(*ctx.Ctx, *settings.ConnectionStruct) (interfaces.ConnInterface, error)

type ChanconnStruct struct {
	cx      *ctx.Ctx
	name    string
	link    *link.LinkStateStruct
	auth    *auth.AuthStruct
	conn    interfaces.ConnInterface
	width   packet.PacketWidth
	recvch  chan *packet.PacketStruct
	pinger  pinger.PingerStruct
	counter counter.CounterStruct
}

func NewConn32(ctx *ctx.Ctx, config *settings.ConnectionStruct, confFunc NewConnFunc) (ci interfaces.ChannelInterface, err error) {
	conn, err := confFunc(ctx, config)
	if err != nil {
		return nil, err
	}

	auth, err := auth.NewAuthStruct(ctx, config.Auth)
	if err != nil {
		return nil, err
	}

	cs := &ChanconnStruct{
		cx:      ctx,
		name:    config.Name,
		link:    link.NewLinkState(ctx, link.LinkModeUpAND),
		auth:    auth,
		width:   packet.PacketWidth32,
		conn:    conn,
		recvch:  make(chan *packet.PacketStruct, packet.DefaultChannelDepth),
		pinger:  pinger.NewPinger32(ctx, PingFreq, PingTimeout),
		counter: counter.NewCounter32(ctx),
	}

	cs.link.AddLink(cs.conn.Link().LinkStateCh)
	cs.link.AddLink(cs.auth.Link().LinkStateCh)

	go cs.goRecv()

	return cs, err
}

func NewConn64(ctx *ctx.Ctx, config *settings.ConnectionStruct, confFunc NewConnFunc) (ci interfaces.ChannelInterface, err error) {

	conn, err := confFunc(ctx, config)
	if err != nil {
		return nil, err
	}

	auth, err := auth.NewAuthStruct(ctx, config.Auth)
	if err != nil {
		return nil, err
	}

	cs := &ChanconnStruct{
		cx:      ctx,
		link:    link.NewLinkState(ctx),
		auth:    auth,
		width:   packet.PacketWidth64,
		conn:    conn,
		recvch:  make(chan *packet.PacketStruct, packet.DefaultChannelDepth),
		pinger:  pinger.NewPinger64(ctx, PingFreq, PingTimeout),
		counter: counter.NewCounter64(ctx),
	}

	cs.link.AddLink(cs.conn.Link().LinkStateCh)
	cs.link.AddLink(cs.auth.Link().LinkStateCh)

	go cs.goRecv()

	return cs, err
}

func (c *ChanconnStruct) MaxLocalMtu() (size packet.PacketSizeType) {
	size = packet.PacketSigSize + packet.PacketSize16Size
	if c.width == packet.PacketWidth32 {
		size += packet.PacketCounter32Size
		size += packet.PacketPing32Size
		if c.width == packet.PacketWidth64 {
			size += packet.PacketCounter64Size
			size += packet.PacketPing64Size
		} else {
			log.FatalfStack("ChanconnStruct:%v", c)
		}
	}
	return size
}

func (c *ChanconnStruct) Link() (link *link.LinkStateStruct) {
	if c == nil {
		return nil
	}
	return c.link
}

func (c *ChanconnStruct) Name() string {
	if c == nil {
		return ""
	}
	return c.name
}
