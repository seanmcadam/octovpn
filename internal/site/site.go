package site

import (
	"time"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/internal/pinger"
	"github.com/seanmcadam/octovpn/internal/tracker"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type SiteStruct struct {
	cx       *ctx.Ctx
	width    packet.PacketWidth
	channels []interfaces.ChannelInterface
	pinger   pinger.PingerStruct
	counter  counter.CounterStruct
	tracker  *tracker.TrackerStruct
	recvch   chan *packet.PacketStruct
}

func (c *SiteStruct) MaxLocalMtu() (size packet.PacketSizeType) {
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

		var max packet.PacketSizeType = 0
		for _, channel := range c.channels {
			s := channel.MaxLocalMtu()
			if s < max {
				max = s
			}
		}

		size += max

		return size
	}
}

func NewSite32(ctx *ctx.Ctx, si []interfaces.SiteInterface) (s *SiteStruct, err error) {

	ss := &SiteStruct{
		cx:       ctx,
		width:    packet.PacketWidth32,
		channels: si,
		pinger:   pinger.NewPinger32(ctx, 1, 2),
		counter:  counter.NewCounter32(ctx),
		tracker:  tracker.NewTracker(ctx, time.Second),
		recvch:   make(chan *packet.PacketStruct, 16),
	}

	go ss.goRecv()
	return ss, err
}


func NewSite64(ctx *ctx.Ctx, si []interfaces.SiteInterface) (s *SiteStruct, err error) {

	ss := &SiteStruct{
		cx:       ctx,
		width:    packet.PacketWidth32,
		channels: si,
		pinger:   pinger.NewPinger32(ctx, 1, 2),
		counter:  counter.NewCounter32(ctx),
		tracker:  tracker.NewTracker(ctx, time.Second),
		recvch:   make(chan *packet.PacketStruct, 16),
	}

	go ss.goRecv()
	return ss, err
}
