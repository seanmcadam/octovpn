package tcpsrv

import (
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/internal/msg"
)

func (tcp *TcpServerStruct) ReceiveHandler(data msg.MsgInterface) {

	if data.FromName() == tcp.msgnode.ParentName() {
		// From Parent

		switch m := data.(type) {
		case *msg.PacketStruct:
			if uint16(m.Packet.Size()) > uint16(tcp.config.Mtu) {
				// return errors.ErrNetPacketTooBig(log.Errf(""))
				log.Errorf("Size() > MTU")
				return
			}

			m.From = tcp.me
			tcp.msgnode.SendChild(m)
		default:
			log.Fatalf("default parent recv reached %T", data)
		}

	} else {
		// From Child

		switch m := data.(type) {
		case *msg.NoticeStruct:
			switch m.Notice {
			case msg.NoticeCLOSED:
				tcp.setState(msg.StateNOLINK)
				tcp.notice(msg.NoticeCLOSED)
				tcp.msgnode.CloseAllChidren()

			default:
				log.Fatalf("[%s] default notice reached:%s", *tcp.me, m.Notice)
			}
		case *msg.StateStruct:
			log.Fatalf("[%s] received State:%s", *tcp.me, m.State)
		case *msg.PacketStruct:
			m.From = tcp.me
			tcp.msgnode.SendParent(m)
		default:
			log.Fatalf("default child recv reached %T", data)
		}

	}
}
