package udpcli

import (
	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func New( config *settings.NetworkStruct) (udpclient interfaces.ChannelInterface, err error){

	log.Debug("Not implemented")
	return udpclient, err
}