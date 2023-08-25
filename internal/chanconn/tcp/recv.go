package tcp

import (
	"io"

	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

// Recv()
func (t *TcpStruct) Recv() (packet *packetconn.ConnPacket, err error) {

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

		packet, err := packetconn.MakePacket(buf)
		if err != nil {
			log.Errorf("TCP Err:%s", err)
			continue
		}

		switch packet.GetType() {
		case packetconn.PACKET_TYPE_TCP:
			t.recvch <- packet

		case packetconn.PACKET_TYPE_TCPAUTH:
			log.Fatal("Not implemented")

		case packetconn.PACKET_TYPE_PONG:
			log.Debug("Got Pong")
			ping := packet.GetPayload()
			t.pinger.Pongch <- ping

		case packetconn.PACKET_TYPE_PING:
			log.Debug("Got Ping")

			ping := packet.GetPayload()
			packet, err := packetconn.NewPacket(packetconn.PACKET_TYPE_PONG, ping)
			if err != nil {
				log.Fatalf("err:%s", err)
			}
			if !t.closed() {
				t.sendch <- packet
			}

		default:
			log.Errorf("Err:%s", err)
			continue
		}

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
