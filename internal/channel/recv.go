package channel

import (
	"log"

	"github.com/seanmcadam/octovpn/octolib/counter"
	"github.com/seanmcadam/octovpn/octolib/packet/packetchan"
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

func (cs *ChannelStruct) recv(data *packetchan.ChanPacket) {
	t := data.GetType()
	switch t {
	case packetchan.CHAN_TYPE_DATA:

		cs.channel.Send(data.CopyDataToAck())
		cs.tracker.Recv(data.Copy())
		cs.recvch <- data

	case packetchan.CHAN_TYPE_ACK:
		cs.tracker.Ack(counter.Counter64(data.GetCounter()))

	case packetchan.CHAN_TYPE_NAK:
		cs.tracker.Nak(counter.Counter64(data.GetCounter()))

	case packetchan.CHAN_TYPE_ERROR:
		log.Fatalf("Unhandled CHAN TYPE ERROR")
	case packetchan.CHAN_TYPE_PING:
		log.Fatalf("Unhandled CHAN TYPE PING")
	case packetchan.CHAN_TYPE_PONG:
		log.Fatalf("Unhandled CHAN TYPE PONG")
	default:
		log.Fatalf("Unhandled CHAN TYPE:%d", t)
	}
}
