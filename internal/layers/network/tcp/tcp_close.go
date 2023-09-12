package tcp

import (
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

var closePacketByte []byte

func init() {
	if closePacket, err := packet.NewPacket(packet.SIG_CONN_CLOSE); err != nil {
		log.Fatalf("TCP Initalization Err for closePacket:%s", err)
	} else {
		closePacketByte = closePacket.ToByte()
	}
}

func (t *TcpStruct) doneChan() <-chan struct{} {
	if t == nil {
		return nil
	}
	return t.cx.DoneChan()
}

func (t *TcpStruct) Cancel() {
	if t == nil {
		return
	}
	if t.conn != nil {
		t.sendclose()
		t.conn.Close()
	}
	t.link.NoLink() // Let upstream know we are down.
	t.link.Close()  // Send notify
	t.cx.Cancel()
}

func (t *TcpStruct) closed() bool {
	if t == nil {
		return true
	}
	return t.cx.Done()
}

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
