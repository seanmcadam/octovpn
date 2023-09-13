package site

import (
	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/layers/channel"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func AssembleSite(cx *ctx.Ctx, config *settings.ConfigSiteStruct) (s *SiteStruct, err error) {

	if cx == nil {
		return nil, errors.ErrSiteBadParameters(log.Errf("no context"))
	}

	if config == nil {
		return nil, errors.ErrSiteBadParameters(log.Errf("no configs"))
	}

	log.Debugf("confg:%v", config)

	var channels []interfaces.ChannelSiteInterface

	if len(config.Servers) > 0 {
		for _, server := range config.Servers {
			if channel, err := channel.ChannelAssembleServer(cx, &server); err != nil {
				return nil, errors.ErrSiteBadParameters(err)
			} else {
				channels = append(channels, channel)
			}
		}
	}

	if len(config.Clients) > 0 {
		for _, client := range config.Clients {
			if channel, err := channel.ChannelAssembleClient(cx, &client); err != nil {
				return nil, errors.ErrSiteBadParameters(err)
			} else {
				channels = append(channels, channel)
			}
		}
	}

	if len(channels) == 0 {
		return nil, log.Errf("No connections for site:%s", config.Sitename)
	}

	if config.Width == 32 || (config.Width == 0 && settings.Width32 == settings.WidthDefault) {

		if s, err = NewSite32(cx, config, channels); err != nil {
			return nil, log.Errf("NewSite32 err:%s", err)
		}

	} else if config.Width == 64 || (config.Width == 0 && settings.Width64 == settings.WidthDefault) {

		if s, err = NewSite64(cx, config, channels); err != nil {
			return nil, log.Errf("NewSite64 err:%s", err)
		}

	} else {
		log.FatalfStack("Bad width:%d", config.Width)
	}

	return s, err
}
