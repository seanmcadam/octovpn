package udp

import (
	"io"

	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// -
//
// Send()
// -
func (u *UdpStruct) Send(p *packet.PacketStruct) (err error) {
	if u == nil || u.sendch == nil {
		return errors.ErrNetNilMethodPointer(log.Errf(""))
	}

	log.Debugf("UDP Send:%v", p)
	select {
	case u.sendch <- p:
	default:
		return errors.ErrNetSendBufferFull(log.Errf(""))
	}

	return nil

}

// -
//
// -
func (u *UdpStruct) goSend() {
	if u == nil {
		return
	}

	defer u.emptysend()
	defer u.Cancel()

	for {
		select {
		case packet := <-u.sendch:
			if err := u.sendpacket(packet); err != nil {
				log.Warnf("sendpacket() Err:%s", err)
				return
			}

		case <-u.doneChan():
			return
		}
	}
}

// -
//
// -
func (u *UdpStruct) sendpacket(p *packet.PacketStruct) (err error) {
	if u == nil {
		return errors.ErrNetNilMethodPointer(log.Err(""))
	}

	var l int
	var raw []byte
	p.DebugPacket("UDP Send")
	if raw, err = p.ToByte(); err != nil {
		return errors.ErrNetParameter(log.Errf("Err:%s", err))
	}

	if u.srv {
		if u.link.IsUp() {
			if l, err = u.conn.WriteToUDP(raw, u.addr); err != nil {
				if err != io.EOF {
					return errors.ErrNetChannelError(log.Errf("UDP Srv Write() Error:%s", err))
				}
				return errors.ErrNetChannelDown(log.Errf("UDP Srv Write() Channel Closed"))
			}

		} else {
			return errors.ErrNetChannelDown(log.Errf("UDP Write() Channel Down"))
		}
	} else {
		if l, err = u.conn.Write(raw); err != nil {
			if err != io.EOF {
				return errors.ErrNetChannelError(log.Errf("UDP Srv Write() Error:%s", err))
			}
			return errors.ErrNetChannelDown(log.Errf("UDP Srv Write() Channel Closed"))
		}
	}

	if l != len(raw) {
		return errors.ErrNetChannelError(log.Errf("UDP Write() Lenth %d != %d", l, len(raw)))
	}

	return nil
}

// -
//
// -
func (u *UdpStruct) sendtestpacket(raw []byte) (err error) {
	if u == nil {
		return
	}

	var l int
	if u.srv {
		if l, err = u.conn.WriteToUDP(raw, u.addr); err != nil {
			return errors.ErrNetSend(log.Errf("UDP Srv WriteToUDP():%v", raw))
		} else {
			if l, err = u.conn.Write(raw); err != nil {
				return errors.ErrNetSend(log.Errf("UDP Cli WriteToUDP():%v", raw))
			}
		}

		if l != len(raw) {
			return errors.ErrNetParameter(log.Errf("UDP Write() Length Error:%d != %d", l, len(raw)))
		}
	}
	return nil
}

// -
//
// -
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
