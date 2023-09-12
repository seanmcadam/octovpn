package main

import (
	"github.com/seanmcadam/octovpn/internal/layers/site"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Functional order
// Read Config
// Probe the network for the interfaces
// Launch the network layers (Connections, channels, sites, routers)
//
//
//

func main() {
	var config *settings.ConfigStruct
	var err error

	cx := ctx.NewContext()
	defer cx.Cancel()

	var sites []*site.SiteStruct


	if config, err = settings.ReadConfig("config.json"); err != nil {
		log.Fatalf("Config file Error:%s", err)
	}

	for _, siteconfig := range config.Sites {
		if site, err := site.AssembleSite(cx, &siteconfig); err != nil {
			log.Fatal("Assemble Site Err:%s", err)
		}else{
		sites = append(sites, site)
		}
	}


	<-sites[0].Link().LinkCloseCh()
	<-sites[1].Link().LinkCloseCh()
	
}
