package tcp

import (
	"io"

	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (t *TcpStruct) RecvChan() <-chan *packet.PacketStruct {

	if t == nil {
		log.FatalStack("nil TcpStruct")
		return nil
	}
	if t.recvch == nil {
		log.Error("Nil recvch pointer")
		return nil
	}

	return t.recvch
}

// Run while connection is running
// Exit when closed
func (t *TcpStruct) goRecv() {
	defer t.emptyrecv()

	for {
		buf := make([]byte, 2048)

		l, err := t.conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Errorf("TCP Read() Error:%s", err)
			}
			return
		}

		if t.closed() {
			return
		}

		buf = buf[:l]

		packet, err := packet.MakePacket(buf)
		if err != nil {
			log.Errorf("TCP MakePacket() Err:%s", err)
			continue
		}

		t.recvch <- packet

		if t.closed() {
			return
		}
	}
}

func (t *TcpStruct) emptyrecv() {
	for {
		select {
		case <-t.recvch:
		default:
			close(t.recvch)
			return
		}
	}
}
