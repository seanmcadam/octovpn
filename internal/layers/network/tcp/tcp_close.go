package tcp

import (
	"github.com/seanmcadam/octovpn/internal/msgbus"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

var closePacketByte []byte

// -
//
// -
func init() {
	if closePacket, err := packet.NewPacket(packet.SIG_CONN_CLOSE); err != nil {
		log.Fatalf("TCP Initalization Err for closePacket:%s", err)
	} else {
		if closePacketByte, err = closePacket.ToByte(); err != nil {
			log.Fatalf("TCP Initalization Err for closePacket:%s", err)
		}
	}
}

// -
//
// -
func (t *TcpStruct) doneChan() <-chan struct{} {
	if t == nil {
		return nil
	}
	return t.cx.DoneChan()
}

// -
//
// -
func (t *TcpStruct) Cancel() {
	if t == nil {
		return
	}

	t.setState(msgbus.StateNOLINK)

	if t.conn != nil {
		t.sendclose()
		t.conn.Close()
	}
	t.cx.Cancel()
}

// -
// sendclose()
// tries to send a last ditch effort close packet
// letting the other side know that the link is closing out side of the IP protocol
// -
func (t *TcpStruct) sendclose() {
	if t == nil {
		return
	}

	if t.conn != nil {
		_, err := t.conn.Write(closePacketByte)

		if err != nil {
			log.Warnf("Err:%s", err)
		}
	}
}
