package chanconn

import (
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (cs *ChanconnStruct) Send(cp *packet.PacketStruct) error {

	if cs.Active() {
		packet, err := packet.NewPacket(packet.SIG_CONN_32_PACKET, cp)
		if err != nil {
			return err
		}
		return cs.send(packet)
	}

	return errors.ErrNetChannelDown
}

func (cs *ChanconnStruct) send(p *packet.PacketStruct) error {
	return cs.conn.Send(p)
}

func (cs *ChanconnStruct) goSend() {

	for {
		select {
		case <-cs.cx.DoneChan():
			return

		case count := <-cs.pinger.GetPingChan():
			packet, err := packet.NewPacket(packet.SIG_CONN_32_PING, count)
			if err != nil {
				log.Errorf("NewPacket Err:%s", err)
				continue
			}

			err = cs.send(packet)
			if err != nil {
				log.Errorf("Send() PING64 Err:%s", err)
			}
		}
	}
}
