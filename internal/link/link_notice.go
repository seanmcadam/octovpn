package link

import (
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (ls *LinkStateStruct) sendNotice(notice *LinkNoticeStateType) {

	ls.linkNoticeState.send(notice)
	ns := notice.Notice()
	switch ns {
	case LinkNoticeCLOSED:
		ls.linkClose.send(notice)
	case LinkNoticeLATENCY:
		ls.linkLatency.send(notice)
	case LinkNoticeLOSS:
		ls.linkLoss.send(notice)
	case LinkNoticeSATURATED:
		ls.linkSaturation.send(notice)
	default:
		log.Fatalf("Unhandled notice:%s", ns)
	}
}
