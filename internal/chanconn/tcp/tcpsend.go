package tcp

import (
	"io"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/packet/packetconn"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Send()
func (t *TcpStruct) Send(packet interfaces.PacketInterface) (err error) {

	go func(p interfaces.PacketInterface) {
		t.sendch <- p
	}(packet)
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

func (t *TcpStruct) sendpacket(packet interfaces.PacketInterface) {
	var l int
	var err error
	packetlen := int(packet.PayloadSize()) + packetconn.Overhead
	l, err = t.conn.Write(packet.ToByte())

	if err != nil {
		if err != io.EOF {
			log.Errorf("TCP Write() Error:%s", err)
		}
		t.cx.Done()
	}

	if l != packetlen {
		log.Errorf("TCP Write() Length Error:%d != %d", l, packetlen)
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
