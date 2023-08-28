package chanconn

import (
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/packet/packetchan"
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

func (cs *ChanconnStruct) Send(cp *packetchan.ChanPacket) error {

	if cs.Active() {
		packet, err := packetconn.NewPacket(packetconn.CONN_TYPE_CHAN, cp)
		if err != nil {
			return err
		}
		return cs.send(packet)
	}

	return errors.ErrNetChannelDown
}

func (cs *ChanconnStruct) send(packet *packetconn.ConnPacket) error {
	return cs.conn.Send(packet)
}
func (cs *ChanconnStruct) goSend() {

	for {
		select {
		case <-cs.cx.DoneChan():
			return

		case count := <-cs.pinger.Pingch:
			packet, err := packetconn.NewPacket(packetconn.CONN_TYPE_PING64, count)
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
