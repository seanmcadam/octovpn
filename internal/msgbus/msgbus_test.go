package msgbus

import (
	"testing"

	"github.com/seanmcadam/ctx"
	"github.com/seanmcadam/octovpn/internal/interfaces"
	"github.com/seanmcadam/octovpn/internal/msg"
)

func TestMsgBus_compile(t *testing.T) {
}

func TestMsgBus_Node(t *testing.T) {
	cx := ctx.New()
	childnode := NewTestNode()
	parentnode := NewTestNode()

	mb := New[interfaces.Node](cx, parentnode, childnode)
	_ = mb

	noticemsg := msg.MsgInterface(msg.NewNotice(childnode.GetInstanceName(), msg.NoticeCLOSED))
	childnode.SendAsChildMsgTest(&noticemsg)
	parentnode.SendAsParentMsgTest(&noticemsg)

	//cx.Cancel()
	<-cx.DoneChan()
}
