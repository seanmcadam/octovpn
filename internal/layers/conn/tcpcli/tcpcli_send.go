package tcpcli

import (
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Send()
func (t *TcpClientStruct) Send(co *packet.PacketStruct) (err error) {
	if t == nil {
		return errors.ErrNetNilMethodPointer(log.Errf(""))
	}

	log.Debugf("TCPCli Send:%v", co)

	if uint16(co.Size()) > uint16(t.config.Mtu) {
		return errors.ErrNetPacketTooBig(log.Errf(""))
	}

	if t.tcpconn != nil {
		if err = t.tcpconn.Send(co); err != nil{
			log.Errorf("Send() Err:%s", err)
		}
		return err
	}

	return errors.ErrNetChannelDown(log.Errf(""))

}
