package link

import (
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (ls *LinkStateStruct) goRecv() {
	if ls == nil {
		return
	}
	defer ls.Cancel()

	for {
		for i, ch := range ls.recvchan {
			var index uint64

			switch i.Uint().(type) {
			case uint32:
				index = uint64(i.Uint().(uint32))
			case uint64:
				index = i.Uint().(uint64)
			}

			//log.GDebugf("Link[%d] recvchan index:%d chanlen:%d", ls.linkname, index, len(ch))

			if index != 0 {
				log.GDebugf("Recv Link[%d] Msg", index)
				var ns LinkNoticeStateType
				select {
				case ns = <-ch:
					log.GDebugf("Recv Link[%d] Msg:%s", index, ns)
					ls.processMessage(ns)
					// Reload the channel
					ls.recvchan[i] = ls.recvfn[i]()
				default:
					// Channel closed, it is dead to me now.
					log.GDebugf("Got DEAD Link Delete:%d", index)
					ls.dellinkch <- i
				}
			} else {
				log.GDebug("Got Link Refresh")

				select {
				case <-ch:
				default:
					log.Error("Empty Channel")
				}

				for {
					select {
					case add := <-ls.addlinkch:
						if add == nil {
							log.FatalStack("nil pointer")
						}
						c := ls.recvcounter.Next()
						ls.recvfn[c] = add.LinkFunc
						ls.recvchan[c] = add.LinkFunc()
						ls.recvstate[c] = add.State
						if add.State != LinkStateNONE {
							ls.processStateChange(noticeState(LinkNoticeNONE, add.State))
						}
					case c := <-ls.dellinkch:
						delete(ls.recvfn, c)
						delete(ls.recvchan, c)
						delete(ls.recvstate, c)
					default:
						break
					}
				}
			}
		}
	}
}
