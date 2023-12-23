package interfaces

import (
	"github.com/seanmcadam/octovpn/octolib/instance"
)

type Node interface {
	GetInstanceName() *instance.InstanceName
	GetParentRecvCh() chan MsgInterface
	GetChildRecvCh() chan MsgInterface
	GetParentMsgHandlerFn() func(MsgInterface)
	GetChildMsgHandlerFn() func(MsgInterface)
}
