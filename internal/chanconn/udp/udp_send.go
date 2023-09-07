package udp

import (
	"io"

	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Send()
func (u *UdpStruct) Send(p *packet.PacketStruct) (err error) {

	go func(p *packet.PacketStruct) {
		u.sendch <- p
	}(p)
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

func (u *UdpStruct) sendpacket(p *packet.PacketStruct) {
	var l int
	var err error
	raw := p.ToByte()
	if u.srv {
		l, err = u.conn.WriteToUDP(raw, u.addr)
	} else {
		l, err = u.conn.Write(raw)
	}

	if err != nil {
		if err != io.EOF {
			log.Errorf("UDP %s Write() Error:%s", u.endpoint(), err)
		}
		u.cx.Cancel()
	}

	if l != len(raw) {
		log.Errorf("UDP Write() Length Error:%d != %d", l, len(raw))
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
