package channel

import (
	"log"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/internal/packet"
)

func (cs *ChannelStruct) goRecv() {

	//defer

	for {
		select {
		case <-cs.cx.DoneChan():
			return
		case data := <-cs.channel.RecvChan():
			cs.recv(data)
		}
	}
}

func (cs *ChannelStruct) recv(data interfaces.PacketInterface) {
	t := data.Type()
	switch t {
	case packet.CHAN_TYPE_PARENT:

		cs.channel.Send(data.CopyAck())
		cs.tracker.Recv(data.Copy())
		cs.recvch <- data

	case packet.CHAN_TYPE_ACK:
		cs.tracker.Ack(counter.Counter32(data.Counter32()))

	case packet.CHAN_TYPE_NAK:
		cs.tracker.Nak(counter.Counter32(data.Counter32()))

	case packet.CHAN_TYPE_ERROR:
		log.Fatalf("Unhandled CHAN TYPE ERROR")
	case packet.CHAN_TYPE_PING64:
		log.Fatalf("Unhandled CHAN TYPE PING")
	case packet.CHAN_TYPE_PONG64:
		log.Fatalf("Unhandled CHAN TYPE PONG")
	default:
		log.Fatalf("Unhandled CHAN TYPE:%d", t)
	}
}
