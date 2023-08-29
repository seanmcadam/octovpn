package udp

import (
	"io"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/packet/packetconn"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Send()
func (u *UdpStruct) Send(packet interfaces.PacketInterface) (err error) {

	go func(p interfaces.PacketInterface) {
		u.sendch <- p
	}(packet)
	return err

}

func (u *UdpStruct) goSend() {

	defer u.emptysend()

	for {
		select {
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

func (u *UdpStruct) sendpacket(packet interfaces.PacketInterface) {
	var l int
	var err error
	packetlen := int(packet.PayloadSize()) + packetconn.Overhead
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
