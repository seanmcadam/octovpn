package interfaces

import "github.com/seanmcadam/octovpn/octolib/instance"

type MsgInterface interface {
	Data() interface{}
	FromName() *instance.InstanceName
}
