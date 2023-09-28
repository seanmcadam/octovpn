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

const ConnPingFreqDefault = 1 * time.Second
const ConnPingTimeoutDefault = 2 * time.Second

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

	pinger, err := pinger.NewPinger32(ctx, ConnPingFreqDefault, ConnPingTimeoutDefault)
	counter := counter.NewCounter32(ctx)

	cs := &ChanconnStruct{
		cx:      ctx,
		name:    config.Name,
		link:    link.NewNameLinkState(ctx, "Conn32", link.LinkModeUpAND),
		auth:    auth,
		width:   packet.PacketWidth32,
		conn:    conn,
		recvch:  make(chan *packet.PacketStruct, packet.DefaultChannelDepth),
		pinger:  pinger,
		counter: counter,
	}

	cs.link.AddLinkStateCh(cs.conn.Link())
	cs.link.AddLinkStateCh(cs.auth.Link())

	go cs.goStart()

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

	pinger, err := pinger.NewPinger64(ctx, ConnPingFreqDefault, ConnPingTimeoutDefault)
	counter := counter.NewCounter64(ctx)

	cs := &ChanconnStruct{
		cx:      ctx,
		link:    link.NewNameLinkState(ctx, "Conn64", link.LinkModeUpAND),
		auth:    auth,
		width:   packet.PacketWidth64,
		conn:    conn,
		recvch:  make(chan *packet.PacketStruct, packet.DefaultChannelDepth),
		pinger:  pinger,
		counter: counter,
	}

	cs.link.NoLink()
	cs.link.AddLinkStateCh(cs.conn.Link())
	cs.link.AddLinkStateCh(cs.auth.Link())

	go cs.goStart()

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
