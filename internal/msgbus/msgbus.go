package msgbus

import (
	"sync"

	"github.com/seanmcadam/ctx"
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/internal/interfaces"
	"github.com/seanmcadam/octovpn/internal/msg"
	"github.com/seanmcadam/octovpn/octolib/instance"
)

// MsgBus is create with a Node Interface that will handle incoming messages, as well and provide messages for the creator to handle.
type MsgBus struct {
	cx               *ctx.Ctx
	closeonce        *sync.Once
	parentCh         chan *msg.MsgInterface
	parentMsgHandler func(*msg.MsgInterface)
	parentInst       *instance.InstanceName
	childCh          chan *msg.MsgInterface
	childMsgHandler  func(*msg.MsgInterface)
	childInst        *instance.InstanceName
}

//
// New()
// Parent (caller) routine provices:
// Parent Recv Handler func
// Parent Close Handler func
// NewNode func

func New[T interfaces.Node](cx *ctx.Ctx, parent T, child T) (mb *MsgBus) {

	mb = &MsgBus{
		cx:               cx,
		closeonce:        new(sync.Once),
		parentCh:         parent.GetParentRecvCh(),
		parentMsgHandler: parent.GetParentMsgHandler(),
		parentInst:       parent.GetInstanceName(),
		childCh:          child.GetChildRecvCh(),
		childMsgHandler:  child.GetChildMsgHandler(),
		childInst:        child.GetInstanceName(),
	}

	go mb.goRun()
	return mb

}

// -
//
// -
func (msgbus *MsgBus) goRun() {
	defer msgbus.close()

	// Could split this up for bi-directioanal flow
	for {
		select {
		case <-msgbus.cx.DoneChan():
			return
		case msg := <-msgbus.parentCh:
			if msg == nil {
				return
			}
			msgbus.childMsgHandler(msg)
		case msg := <-msgbus.childCh:
			if msg == nil {
				return
			}
			msgbus.parentMsgHandler(msg)
		}
	}
}

// -
// Just call it one time, but multiple places can call it
// -
func (msgbus *MsgBus) close() {

	c := func() {
		log.Debugf("MsgBus Closeing: %s <-> %s", msgbus.parentInst, msgbus.childInst)

		msgbus.cx.Cancel()
	}

	msgbus.closeonce.Do(c)
}
