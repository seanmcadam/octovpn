package udp

import (
	"io"
	"net"
	"reflect"

	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

// Recv()
func (u *UdpStruct) RecvChan() <-chan *packet.PacketStruct {
	if u == nil || u.recvch == nil {
		log.Debug("UPD RecvChan() Nil")
		return nil
	}

	return u.recvch
}

// Run while connection is running
// Exit when closed
func (u *UdpStruct) goRecv() {
	if u == nil {
		return
	}

	defer u.emptyrecv()
	defer u.link.Down()

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

		p, err := packet.MakePacket(buf)
		p.DebugPacket("UDP RECV")
		if err != nil {
			log.Errorf("UDP MakePacket() Err:%s", err)
			continue
		}

		u.recvch <- p

		if u.closed() {
			return
		}
	}
}

func (u *UdpStruct) emptyrecv() {
	if u == nil {
		return
	}

	for {
		select {
		case <-u.recvch:
		default:
			close(u.recvch)
			return
		}
	}
}
