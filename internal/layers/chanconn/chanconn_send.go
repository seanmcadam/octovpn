package chanconn

import (
	"fmt"

	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (cs *ChanconnStruct) Send(cp *packet.PacketStruct) error {

	if cs == nil {
		return errors.ErrNetNilMethodPointer(log.Errf(""))
	}

	if cs.link.IsUp() {
		packet, err := packet.NewPacket(packet.SIG_CONN_32_PACKET, cp, <-cs.counter.GetCountCh())
		if err != nil {
			return err
		}
		packet.DebugPacket("Chanconn Send")
		return cs.send(packet)
	}

	return errors.ErrNetChannelDown(log.Errf(""))
}

func (cs *ChanconnStruct) send(p *packet.PacketStruct) error {
	if cs == nil {
		return errors.ErrNetNilMethodPointer(log.Errf(""))
	}

	if cs.conn.Link().IsDown() {
		return fmt.Errorf("link shows down")
	}

	return cs.conn.Send(p)
}

func (cs *ChanconnStruct) goSend() {

	if cs == nil {
		return
	}

	defer cs.Cancel()

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
