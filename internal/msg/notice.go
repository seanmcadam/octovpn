package msg

import (
	"github.com/seanmcadam/octovpn/octolib/instance"
	log "github.com/seanmcadam/loggy"
)

type MsgNotice string

const (
	NoticeCLOSED    MsgNotice = "Close"
	NoticeLOSS      MsgNotice = "Loss"
	NoticeLATENCY   MsgNotice = "Latency"
	NoticeSATURATED MsgNotice = "Saturated"
)

type NoticeHandleTable struct {
	handler map[MsgNotice]func(*NoticeStruct)
}

func NewNoticeHandler() (nht *NoticeHandleTable) {
	nht = &NoticeHandleTable{
		handler: map[MsgNotice]func(*NoticeStruct){},
	}
	nht.handler[NoticeCLOSED] = emptyNoticeHandle
	nht.handler[NoticeLOSS] = emptyNoticeHandle
	nht.handler[NoticeLATENCY] = emptyNoticeHandle
	nht.handler[NoticeSATURATED] = emptyNoticeHandle
	return nht
}

func (nht *NoticeHandleTable) Run(ns *NoticeStruct) {
	if handlerFn, ok := nht.handler[ns.Notice]; ok {
		handlerFn(ns)
	} else {
		log.FatalfStack("Unknown Run Notice Type %v", ns)
	}
}

func (nht *NoticeHandleTable) AddHandle(ms MsgNotice, fn func(*NoticeStruct)) {
	nht.handler[ms] = fn
}

func (nht *NoticeHandleTable) CallHandle(ss *NoticeStruct) {
	nht.handler[ss.Notice](ss)
}

type NoticeStruct struct {
	Notice MsgNotice
	From   *instance.InstanceName
}

func NewNotice(from *instance.InstanceName, notice MsgNotice) (ns *NoticeStruct) {
	ns = &NoticeStruct{
		Notice: notice,
		From:   from,
	}
	return ns
}

func (n *NoticeStruct) FromName() *instance.InstanceName {
	return n.From
}

func (n *NoticeStruct) Data() interface{} {
	return &n.Notice
}

func emptyNoticeHandle(ns *NoticeStruct) {
	log.ErrorfStack("StateEmptyHandler From:%s Notice:%s", ns.From, ns.Notice)
}
