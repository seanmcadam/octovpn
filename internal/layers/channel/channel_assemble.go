package channel

import (
	"github.com/seanmcadam/octovpn/internal/layers/chanconn"
	"github.com/seanmcadam/octovpn/internal/layers/conn/tcpcli"
	"github.com/seanmcadam/octovpn/internal/layers/conn/tcpsrv"
	"github.com/seanmcadam/octovpn/internal/layers/conn/udpcli"
	"github.com/seanmcadam/octovpn/internal/layers/conn/udpsrv"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

//func NewConn32(ctx *ctx.Ctx, config *settings.ConnectionStruct, confFunc NewConnFunc) (ci interfaces.ChannelInterface, err error) {

func ChannelAssembleServer(ctx *ctx.Ctx, config *settings.ConnectionStruct) (cs *ChannelStruct, err error) {
	return channelAssemble(ctx, config, true)
}

func ChannelAssembleClient(ctx *ctx.Ctx, config *settings.ConnectionStruct) (cs *ChannelStruct, err error) {
	return channelAssemble(ctx, config, false)
}

func channelAssemble(ctx *ctx.Ctx, config *settings.ConnectionStruct, server bool) (cs *ChannelStruct, err error) {

	log.Debugf("confg:%v", config)

	var connFunc chanconn.NewConnFunc

	switch config.Proto {
	case settings.TCP:
		fallthrough
	case settings.TCP4:
		fallthrough
	case settings.TCP6:
		if server {
			connFunc = tcpsrv.New
		} else {
			connFunc = tcpcli.New
		}

	case settings.UDP:
		fallthrough
	case settings.UDP4:
		fallthrough
	case settings.UDP6:
		if server {
			connFunc = udpsrv.New
		} else {
			connFunc = udpcli.New
		}

	default:
		return nil, errors.ErrChanUnhandledProtocol(log.Errf("default reached"))
	}

	if config.Width == 32 || (config.Width == 0 && settings.Width32 == settings.WidthDefault) {
		if ci, err := chanconn.NewConn32(ctx, config, connFunc); err != nil {
			return nil, errors.ErrChanNew(log.Errf("NewConn32() Err:%s", err))
		} else {
			if cs, err = NewChannel32(ctx, ci); err != nil {
				return nil, errors.ErrChanNew(log.Errf("NewChannel32() Err:%s", err))
			}
		}
	} else if config.Width == 64 || (config.Width == 0 && settings.Width64 == settings.WidthDefault) {
		if ci, err := chanconn.NewConn64(ctx, config, connFunc); err != nil {
			return nil, errors.ErrChanNew(log.Errf("NewConn64() Err:%s", err))
		} else {
			if cs, err = NewChannel64(ctx, ci); err != nil {
				return nil, errors.ErrChanNew(log.Errf("NewChannel64() Err:%s", err))
			}
		}

	}

	return cs, err
}
