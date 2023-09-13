package tcp

import (
	"io"

	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

//-
// Send()
//-
func (t *TcpStruct) Send(p *packet.PacketStruct) (err error) {
	if t == nil || t.sendch == nil {
		return errors.ErrNetNilMethodPointer(log.Errf(""))
	}

	log.Debugf("TCP Send:%v", p)
	select {
	case t.sendch <- p:
	default:
		return errors.ErrNetSendBufferFull(log.Errf(""))
	}

	return nil

}

//-
// goSend()
// handle the send buffer
//-
func (t *TcpStruct) goSend() {
	if t == nil {
		return
	}

	defer t.emptysend()
	defer t.Cancel()

	for {
		select {
		case packet := <-t.sendch:
			if err := t.sendpacket(packet); err != nil {
				log.Warnf("sendpacket() Err:%s", err)
				return
			}

		case <-t.doneChan():
			return
		}
	}
}

//-
// sendpacket()
//-
func (t *TcpStruct) sendpacket(p *packet.PacketStruct) (err error) {
	if t == nil {
		return errors.ErrNetNilMethodPointer(log.Err(""))
	}

	p.DebugPacket("TCP Send")

	if raw, err := p.ToByte(); err != nil{
		return errors.ErrNetParameter(log.Errf("Err:%s", err))

	}else if l, err := t.conn.Write(raw); err != nil {
		if err != io.EOF {
			return errors.ErrNetChannelError(log.Errf("TCP Write() Error:%s", err))
		}
		return errors.ErrNetChannelDown(log.Errf("TCP Write() Channel Closed"))
	} else if l != len(raw) {
		return errors.ErrNetChannelError(log.Errf("TCP Write() Lenth Error:%d != %d", l, len(raw)))
	}

	return nil
}

//-
// sendtestpacket()
//-
func (t *TcpStruct) sendtestpacket(raw []byte)(err error) {
	if t == nil {
		return
	}

	log.Debugf("TCP RAW Send:%v", raw)

	if l, err := t.conn.Write(raw); err != nil{
		if err != io.EOF {
			return errors.ErrNetChannelError(log.Errf("TCP RAW Write() Error:%s, Closing Connection", err))
		}
	}else if l != len(raw) {
		return errors.ErrNetChannelError(log.Errf("TCP RAW Write() length Error:%d != %d", l, len(raw)))
	}

	return nil
}

//-
// Clean up sendch before exit
//-
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
