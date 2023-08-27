package tcp

import (
	"io"

	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

func (t *TcpStruct) RecvChan() <-chan *packetconn.ConnPacket {

	if t == nil {
		log.FFatal("Nil struct pointer")
		return nil
	}
	if t.recvch == nil {
		log.FFatal("Nil recvch pointer")
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

		packet, err := packetconn.MakePacket(buf)
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
