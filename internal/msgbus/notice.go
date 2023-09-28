package msgbus

import "fmt"

type NoticeMsg uint8

const (
	NoticeNone       NoticeMsg = 0x00
	NoticeLoss       NoticeMsg = 0x01
	NoticeLatency    NoticeMsg = 0x02
	NoticeSaturation NoticeMsg = 0x04
	NoticeClosed     NoticeMsg = 0xFF
)

type NoticeMsgStruct struct {
	Notice NoticeMsg
}

// -
// Notice()
// -
func (mb *MsgBus) Notice(target MsgTarget, notice *NoticeMsg) {
	topic := fmt.Sprintf("%s:%s", string(target), string(MsgNotice))
	mb.eventbus.Publish(topic, notice)
}

// -
// Register Notice Handler
// -
func (mb *MsgBus) NoticeHandler(source MsgTarget, handler func(...interface{})) (err error) {
	topic := fmt.Sprintf("%s:%s", string(source), string(MsgNotice))
	return mb.eventbus.Subscribe(topic, handler)
}