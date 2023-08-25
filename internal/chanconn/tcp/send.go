package tcp

import (
	"io"

	"github.com/seanmcadam/octovpn/internal/chanconn"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Send()
func (t *TcpStruct) Send(packet *chanconn.ConnPacket) (err error) {

	go func(p *chanconn.ConnPacket) {
		t.sendch <- p
	}(packet)
	return err

}

func (t *TcpStruct) goSend() {

	defer t.emptysend()

	for {
		select {
		case packet := <-t.sendch:
			packetlen := int(packet.GetLength()) + chanconn.PacketOverhead
			l, err := t.conn.Write(packet.ToByte())
			if err != nil {
				if err != io.EOF {
					log.Errorf("TCP Write():%s", err)
				}
				return
			}

			if l != packetlen {
				log.Errorf("TCP Write() Length Error:%d != %d", l, packetlen)
				return
			}

		case <-t.Closech:
			return
		}
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
