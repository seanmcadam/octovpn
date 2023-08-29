package udp

import (
	"io"
	"net"
	"reflect"

	"github.com/seanmcadam/octovpn/interfaces"
	"github.com/seanmcadam/octovpn/internal/packet/packetconn"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Recv()
func (u *UdpStruct) RecvChan() <-chan interfaces.PacketInterface {
	return u.recvch
}

// Run while connection is running
// Exit when closed
func (u *UdpStruct) goRecv() {
	defer u.emptyrecv()

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
			log.Errorf("UDP MakePacket() Err:%s", err)
			continue
		}

		u.recvch <- packet

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
