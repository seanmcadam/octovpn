package chanconn

import (
	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/internal/packet/packetconn"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (cs *ChanconnStruct) Send(cp interfaces.PacketInterface) error {

	if cs.Active() {
		packet, err := packetconn.NewPacket(packet.CONN_TYPE_PARENT, cp)
		if err != nil {
			return err
		}
		return cs.send(packet)
	}

	return errors.ErrNetChannelDown
}

func (cs *ChanconnStruct) send(packet interfaces.PacketInterface) error {
	return cs.conn.Send(packet)
}
func (cs *ChanconnStruct) goSend() {

	for {
		select {
		case <-cs.cx.DoneChan():
			return

		case count := <-cs.pinger.Pingch:
			packet, err := packetconn.NewPacket(packet.CONN_TYPE_PING64, count)
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
