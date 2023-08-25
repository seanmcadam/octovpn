package tcp

import (
	"io"

	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

// Send()
func (t *TcpStruct) Send(packet *packetconn.ConnPacket) (err error) {

	go func(p *packetconn.ConnPacket) {
		t.sendch <- p
	}(packet)
	return err

}

func (t *TcpStruct) goSend() {

	defer t.emptysend()

	for {
		select {
		case ping := <-t.pinger.Pingch:

			packet, err := packetconn.NewPacket(packetconn.PACKET_TYPE_PING, ping.ToByte())
			if err != nil {
				log.Fatalf("err:%s", err)
			}
			t.sendpacket(packet)

		case packet := <-t.sendch:
			t.sendpacket(packet)

		case <-t.cx.DoneChan():
			return
		}

		if t.closed() {
			return
		}
	}
}

func (t *TcpStruct) sendpacket(packet *packetconn.ConnPacket) {
	var l int
	var err error
	packetlen := int(packet.GetLength()) + packetconn.PacketOverhead
	l, err = t.conn.Write(packet.ToByte())

	if err != nil {
		if err != io.EOF {
			log.Errorf("TCP Write() Error:%s", err)
		}
		//t.Close()
		t.cx.Done()
	}

	if l != packetlen {
		log.Errorf("TCP Write() Length Error:%d != %d", l, packetlen)
		//t.Close()
		t.cx.Done()
	}
}

func (t *TcpStruct) emptysend() {
	for {
		select {
		case <-t.sendch:
		default:
			close(t.sendch)
			return
		}
	}
}
