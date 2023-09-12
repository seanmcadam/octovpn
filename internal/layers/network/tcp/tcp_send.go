package tcp

import (
	"io"

	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Send()
func (t *TcpStruct) Send(p *packet.PacketStruct) (err error) {
	if t == nil || t.sendch == nil {
		return errors.ErrNetNilPointerMethod(log.Errf(""))
	}

	log.Debugf("TCP Send:%v", p)

	go func(p *packet.PacketStruct) {
		t.sendch <- p
	}(p)
	return err

}

func (t *TcpStruct) goSend() {
	if t == nil {
		return
	}

	defer t.emptysend()

	for {
		select {
		case packet := <-t.sendch:
			t.sendpacket(packet)

		case <-t.doneChan():
			return
		}

		if t.closed() {
			return
		}
	}
}

func (t *TcpStruct) sendpacket(p *packet.PacketStruct) {
	if t == nil {
		return
	}

	p.DebugPacket("TCP Send")

	raw := p.ToByte()
	l, err := t.conn.Write(raw)
	if err != nil {
		if err != io.EOF {
			log.Errorf("TCP Write() Error:%s, Closing Connection", err)
		}
		t.Cancel()
	}
	if l != len(raw) {
		log.Errorf("TCP Write() Send length:%d, Closing Connection", l, len(raw))
		t.Cancel()
	}
}

func (t *TcpStruct) emptysend() {
	if t == nil {
		return
	}

	for {
		select {
		case <-t.sendch:
		default:
			close(t.sendch)
			return
		}
	}
}
