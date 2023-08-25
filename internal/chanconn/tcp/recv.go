package tcp

import (
	"io"

	"github.com/seanmcadam/octovpn/internal/chanconn"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Recv()
func (t *TcpStruct) Recv() (packet *chanconn.ConnPacket, err error) {

	packet = <-t.recvch

	return packet, err
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

		buf = buf[:l]

		packet, err := chanconn.MakePacket(buf)
		if err != nil {
			log.Errorf("Err:%s", err)
			continue
		}

		switch packet.GetType() {
		case chanconn.PACKET_TYPE_TCP:
		case chanconn.PACKET_TYPE_TCPAUTH:
		default:
			log.Errorf("Err:%s", err)
			continue
		}

		if t.closed() {
			return
		}

		t.recvch <- packet
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
