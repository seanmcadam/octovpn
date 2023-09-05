package tcp

import (
	"io"

	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Send()
func (t *TcpStruct) Send(p *packet.PacketStruct) (err error) {

	go func(p *packet.PacketStruct) {
		t.sendch <- p
	}(p)
	return err

}

func (t *TcpStruct) goSend() {

	defer t.emptysend()

	for {
		select {
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

func (t *TcpStruct) sendpacket(p *packet.PacketStruct) {
	raw := p.ToByte()
	l, err := t.conn.Write(raw)
	if err != nil {
		if err != io.EOF {
			log.Errorf("TCP Write() Error:%s, Closing Connection", err)
		}
		t.cx.Done()
	}
	if l != len(raw) {
		log.Errorf("TCP Write() Send length:%d, Closing Connection", l, len(raw))
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
