package udp

import (
	"io"
	"net"
	"reflect"

	"github.com/seanmcadam/octovpn/octolib/log"
	"github.com/seanmcadam/octovpn/octolib/packet/packetconn"
)

// Recv()
func (u *UdpStruct) RecvChan() <-chan *packetconn.ConnPacket {
	return u.recvch
}

// Run while connection is running
// Exit when closed
func (u *UdpStruct) goRecv() {
	defer u.emptyrecv()

	if !u.srv {
		u.pinger.TurnOn() // If this is a client turn Ping on
	}

	for {
		buf := make([]byte, 2048)
		var l int
		var err error
		var addr *net.UDPAddr

		if u.srv {
			l, addr, err = u.conn.ReadFromUDP(buf)
		} else {
			l, err = u.conn.Read(buf)
		}

		if err != nil {
			if err != io.EOF {
				log.Errorf("UDP %s Read() Error:%s", u.endpoint(), err)
			}
			u.cx.Cancel()
			return
		}

		if u.addr == nil {
			u.addr = addr
		} else if !reflect.DeepEqual(u.addr.IP, addr.IP) {
			log.Errorf("%s != %s - Dropping Packet", u.addr.IP, addr.IP)
			continue
		}

		buf = buf[:l]

		packet, err := packetconn.MakePacket(buf)
		if err != nil {
			log.Errorf("Err:%s", err)
			continue
		}

		switch packet.GetType() {
		case packetconn.PACKET_TYPE_UDP:
			u.recvch <- packet

		case packetconn.PACKET_TYPE_UDPAUTH:
			log.Fatal("Not Implemented")

		case packetconn.PACKET_TYPE_PONG:
			log.Debug("Got Pong")
			ping := packet.GetPayload()
			u.pinger.Pongch <- ping.([]byte)

		case packetconn.PACKET_TYPE_PING:
			log.Debug("Got Ping")

			if u.srv {
				u.pinger.TurnOn() // Turn this pinger on once the first PING is receieved
			}

			ping := packet.GetPayload()
			packet, err := packetconn.NewPacket(packetconn.PACKET_TYPE_PONG, ping)
			if err != nil {
				log.Fatalf("err:%s", err)
			}
			u.sendch <- packet

		default:
			log.Errorf("Err:%s", err)
			continue
		}

		if u.closed() {
			return
		}

	}
}

func (u *UdpStruct) emptyrecv() {
	for {
		select {
		case <-u.recvch:
		default:
			close(u.recvch)
			return
		}
	}
}
