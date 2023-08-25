package udp

import (
	"io"

	"github.com/seanmcadam/octovpn/internal/chanconn"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Send()
func (u *UdpStruct) Send(packet *chanconn.ConnPacket) (err error) {

	go func(p *chanconn.ConnPacket) {
		u.sendch <- p
	}(packet)
	return err

}

func (u *UdpStruct) goSend() {

	defer u.emptysend()

	for {
		select {
		case packet := <-u.sendch:
			var l int
			var err error
			packetlen := int(packet.GetLength()) + chanconn.PacketOverhead
			if u.srv {
				l, err = u.conn.WriteToUDP(packet.ToByte(), u.addr)
			} else {
				l, err = u.conn.Write(packet.ToByte())
			}
			if err != nil {
				if err != io.EOF {
					log.Errorf("UDP %s Write() Error:%s", u.endpoint(), err)
				}
				u.Close()
				return
			}

			if l != packetlen {
				log.Errorf("UDP Write() Length Error:%d != %d", l, packetlen)
				return
			}

		case <-u.Closech:
			return
		}
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
