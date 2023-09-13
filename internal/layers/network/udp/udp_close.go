package udp

import (
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

var closePacketByte []byte

//-
//
//-
func init() {
	if closePacket, err := packet.NewPacket(packet.SIG_CONN_CLOSE); err != nil {
		log.Fatalf("UDP Initalization Err for closePacket:%s", err)
	} else {
		if closePacketByte, err = closePacket.ToByte(); err != nil {
			log.Fatalf("UDP Initalization Err for closePacketByte:%s", err)
		}
	}
}

//-
//
//-
func (u *UdpStruct) doneChan() <-chan struct{} {
	if u == nil {
		return nil
	}

	return u.cx.DoneChan()
}

//-
//
//-
func (u *UdpStruct) Cancel() {
	if u == nil {
		return
	}
	u.sendclose()
	u.link.Close()
	u.cx.Cancel()
}

//-
//
//-
func (u *UdpStruct) closed() bool {
	if u == nil {
		return true
	}
	return u.cx.Done()
}

//-
//
//-
func (u *UdpStruct) sendclose() {
	if u == nil {
		return
	}

	var err error
	if u.conn != nil {
		if u.srv {
			_, err = u.conn.WriteToUDP(closePacketByte, u.addr)
		} else {
			_, err = u.conn.Write(closePacketByte)
		}

		if err != nil {
			log.Warnf("Err:%s", err)
		}
	}

}
