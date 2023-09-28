package tcpsrv

import (
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Send()
func (t *TcpServerStruct) Send(co *packet.PacketStruct) (err error) {
	if t == nil {
		return errors.ErrNetNilMethodPointer(log.Errf(""))
	}

	log.Debugf("TCPSrc Send:%v", co)

	if uint16(co.Size()) > uint16(t.config.Mtu) {
		return errors.ErrNetPacketTooBig(log.Errf(" size:%d > %d", uint16(co.Size()), uint16(t.config.Mtu)))
	}

	for i, conn := range t.tcpconn {
		if conn.Link().IsUp() {
			if err = conn.Send(co); err != nil {
				log.Errorf("Send() on %s Err", i, err)
				return err
			}
			return nil
		}
	}
	//if t.tcpconn != nil {
	//	if err = t.tcpconn.Send(co); err != nil {
	//		log.Errorf("Send() Err", err)
	//	}
	//	return err
	//}

	return errors.ErrNetChannelDown(log.Errf("No open chanels in tcp server"))

}
