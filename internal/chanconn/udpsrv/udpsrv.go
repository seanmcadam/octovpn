package udpsrv

import (
	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func New( config *settings.NetworkStruct) (udpserver interfaces.ChannelInterface, err error){

	log.Debug("Not implemented")
	return udpserver, err
}