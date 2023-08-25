package udp

import (
	"io"

	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

// Send()
func (u *UdpStruct) Send(packet *packetconn.ConnPacket) (err error) {

	go func(p *packetconn.ConnPacket) {
		u.sendch <- p
	}(packet)
	return err

}

func (u *UdpStruct) goSend() {

	defer u.emptysend()

	for {
		select {
		case ping := <-u.pinger.Pingch:

			packet, err := packetconn.NewPacket(packetconn.PACKET_TYPE_PING, ping.ToByte())
			if err != nil {
				log.Fatalf("err:%s", err)
			}
			u.sendpacket(packet)

		case packet := <-u.sendch:
			u.sendpacket(packet)

		case <-u.cx.DoneChan():
			return
		}

		if u.closed() {
			return
		}
	}

}

func (u *UdpStruct) sendpacket(packet *packetconn.ConnPacket) {
	var l int
	var err error
	packetlen := int(packet.GetLength()) + packetconn.PacketOverhead
	if u.srv {
		l, err = u.conn.WriteToUDP(packet.ToByte(), u.addr)
	} else {
		l, err = u.conn.Write(packet.ToByte())
	}

	if err != nil {
		if err != io.EOF {
			log.Errorf("UDP %s Write() Error:%s", u.endpoint(), err)
		}
		u.cx.Cancel()
	}

	if l != packetlen {
		log.Errorf("UDP Write() Length Error:%d != %d", l, packetlen)
		u.cx.Cancel()
	}
}

func (u *UdpStruct) emptysend() {
	for {
		select {
		case <-u.sendch:
		default:
			close(u.sendch)
			return
		}
	}
}
