package udp

import (
	"io"
	"net"
	"reflect"
	"time"

	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

//-
// Recv()
//-
func (u *UdpStruct) RecvChan() <-chan *packet.PacketStruct {
	if u == nil || u.recvch == nil {
		log.Debug("UPD RecvChan() Nil")
		return nil
	}

	u.link.Connected()

	return u.recvch
}

//-
// goRecv()
// Run while Listener is running
// Exit when closed
//-
func (u *UdpStruct) goRecv() {
	if u == nil {
		return
	}

	defer u.emptyrecv()
	defer u.Cancel()

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
			return
		}

		//
		// if addr is nil, new connection
		//
		if u.addr == nil {
			u.addr = addr
			u.link.Connected()
		} else if !reflect.DeepEqual(u.addr.IP, addr.IP) {
			log.Errorf("%s != %s - Dropping Packet", u.addr.IP, addr.IP)
			log.Debug("Need to look at reestablishing the connection here")
			continue
		}

		buf = buf[:l]
		log.Debugf("UDP Read %v", buf)

		p, err := packet.MakePacket(buf)
		if err != nil || p == nil{
			log.Errorf("MakePacket() Err:%s", err)
			return
		}

		//
		// Did I soft close packet?
		//
		if p.Sig().Close() {
			if u.srv {
				u.link.Listen()
			} else {
				return
			}
		}

		p.DebugPacket("UDP RECV")
		if err != nil {
			log.Errorf("UDP MakePacket() Err:%s", err)
			return
		}

		p.DebugPacket("UDP Recv")
		u.recvresettimeout <- nil
		u.recvch <- p

		if u.closed() {
			return
		}
	}
}

//-
//
//-
func (u *UdpStruct) goRecvTimeout() {
	if u == nil {
		return
	}

	defer u.Cancel()
	//
	// Event
	//  reset
	//   reset timer
	// timeout recv
	//  if up
	//   Reset timeout
	//  if Down
	//   Set Close time out
	//
	// timeout close
	//  if down
	//   if srv
	//    addr = nil
	//   if cli
	//    Cancel()
	//
	//

	timeouttime := UDPRecvTimeout

	for {
		select {
		case <-u.link.LinkUpDownCh():
			if u.link.IsUp() {
				timeouttime = UDPRecvTimeout
			} else {
				timeouttime = UDPCloseTimeout
			}

		case <-u.recvresettimeout:
			u.link.Link()
			timeouttime = UDPRecvTimeout

		case <-time.After(timeouttime):
			if timeouttime == UDPRecvTimeout {
				u.link.NoLink()
				timeouttime = UDPCloseTimeout
			} else {
				if u.srv {
					u.link.Listen() // Go back to just lisening
					u.addr = nil
				} else {
					return // Kill the connection, let the parent restart it
				}
			}

		}
	}
}

//-
//
//-
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
