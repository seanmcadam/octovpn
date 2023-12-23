package msgbus

import (
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/internal/msg"
	"github.com/seanmcadam/octovpn/octolib/instance"
)

//
// This is a test node for handling messages and passing them up or down the chain of handlers
//

var inst *instance.Instance

func init() {
	inst = instance.New()
}

type TestNode struct {
	name          *instance.InstanceName
	parentCh      chan *msg.MsgInterface
	childCh       chan *msg.MsgInterface
	stateHandler  *msg.StateHandleTable
	noticeHandler *msg.NoticeHandleTable
	packetHandler *msg.PacketHandleTable
}

func NewTestNode() (tn *TestNode) {
	tn = &TestNode{
		name:          inst.Next(),
		parentCh:      make(chan *msg.MsgInterface),
		childCh:       make(chan *msg.MsgInterface),
		stateHandler:  msg.NewStateHandler(),
		noticeHandler: msg.NewNoticeHandler(),
		packetHandler: msg.NewPacketHandler(),
	}

	return tn
}

func (tn *TestNode) GetInstanceName() *instance.InstanceName {
	return tn.name
}

func (tn *TestNode) GetParentRecvCh() chan *msg.MsgInterface {
	return tn.parentCh
}

func (tn *TestNode) GetChildRecvCh() chan *msg.MsgInterface {
	return tn.childCh
}

func (tn *TestNode) GetChildMsgHandler() func(*msg.MsgInterface) {
	return tn.childMsgHandler
}

func (tn *TestNode) GetParentMsgHandler() func(*msg.MsgInterface) {
	return tn.parentMsgHandler
}

// -
//
// -
func (tn *TestNode) parentMsgHandler(msg *msg.MsgInterface) {

	log.Debugf("[%s]:%v", tn.name, msg)
	tn.childCh <- msg
}

// -
//
// -
func (tn *TestNode) childMsgHandler(msg *msg.MsgInterface) {
	log.Debugf("[%s]:%v", tn.name, msg)
	tn.parentCh <- msg
}

// ---------------------------------
// -
//
// -
func (tn *TestNode) SendToParentMsgTest(msg *msg.MsgInterface) {
	log.Debugf("[%s]:%v", tn.name, msg)
	select {
	case tn.parentCh <- msg:
	default:
		log.Debugf("Unable to send as Child: %v", msg)
	}
}

// -
//
// -
func (tn *TestNode) SendToChildMsgTest(msg *msg.MsgInterface) {
	log.Debugf("[%s]:%v", tn.name, msg)
	select {
	case tn.childCh <- msg:
	default:
		log.Debugf("Unable to send as Parent: %v", msg)
	}
}
