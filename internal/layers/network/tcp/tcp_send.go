package tcp

import (
	"io"

	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

//
// Send()
//
func (t *TcpStruct) Send(p *packet.PacketStruct) (err error) {
	if t == nil || t.sendch == nil {
		return errors.ErrNetNilPointerMethod(log.Errf(""))
	}

	log.Debugf("TCP Send:%v", p)
	select{
	case t.sendch <- p:
	default:
		return errors.ErrNetSendBufferFull(log.Errf(""))
	}

	return nil

}

//
// goSend() 
// handle the send buffer
//
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
	}
}

//
// sendpacket()
//
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

//
// sendtestpacket()
//
func (t *TcpStruct) sendtestpacket(raw []byte) {
	if t == nil {
		return
	}

	log.Debugf("TCP RAW Send:%v", raw)

	l, err := t.conn.Write(raw)
	if err != nil {
		if err != io.EOF {
			log.Errorf("TCP RAW Write() Error:%s, Closing Connection", err)
		}
		t.Cancel()
	}
	if l != len(raw) {
		log.Errorf("TCP RAW Write() Send length:%d, Closing Connection", l, len(raw))
		t.Cancel()
	}
}



// 
// Clean up sendch before exit
//
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
