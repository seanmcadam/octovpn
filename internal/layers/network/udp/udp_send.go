package udp

import (
	"io"

	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Send()
func (u *UdpStruct) Send(p *packet.PacketStruct) (err error) {
	if u == nil || u.sendch == nil {
		return errors.ErrNetNilPointerMethod(log.Errf(""))
	}

	go func(p *packet.PacketStruct) {
		u.sendch <- p
	}(p)
	return err

}

func (u *UdpStruct) goSend() {
	if u == nil {
		return
	}

	defer u.emptysend()
	defer u.Cancel()

	for {
		select {
		case packet := <-u.sendch:
			u.sendpacket(packet)

		case <-u.doneChan():
			return
		}

		if u.closed() {
			return
		}
	}

}

func (u *UdpStruct) sendpacket(p *packet.PacketStruct) {
	if u == nil {
		return
	}

	var l int
	var err error
	p.DebugPacket("UDP Send")
	raw := p.ToByte()
	if u.srv {
		l, err = u.conn.WriteToUDP(raw, u.addr)
		log.Debugf("UDP WriteToUDP():%v", raw)
	} else {
		l, err = u.conn.Write(raw)
		log.Debugf("UDP Write():%v", raw)
	}

	if err != nil {
		if err != io.EOF {
			log.Errorf("UDP %s Write() Error:%s", u.endpoint(), err)
		}
		u.Cancel()
	}

	if l != len(raw) {
		log.Errorf("UDP Write() Length Error:%d != %d", l, len(raw))
		u.Cancel()
	}
}

func (u *UdpStruct) emptysend() {
	if u == nil {
		return
	}

	for {
		select {
		case <-u.sendch:
		default:
			close(u.sendch)
			return
		}
	}
}
