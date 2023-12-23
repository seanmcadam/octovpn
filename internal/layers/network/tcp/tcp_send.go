package tcp

import (
	"io"

	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/errors"
)

var closePacketByte []byte

// -
//
// -
func init() {
	if closePacket, err := packet.NewPacket(packet.SIG_CONN_CLOSE); err != nil {
		log.Fatalf("TCP Initalization Err for closePacket:%s", err)
	} else {
		if closePacketByte, err = closePacket.ToByte(); err != nil {
			log.Fatalf("TCP Initalization Err for closePacket:%s", err)
		}
	}
}

// -
// goSend()
// handle the send buffer
// -
func (t *TcpStruct) goSend() {
	if t == nil {
		return
	}

	defer t.close()

	for {
		select {
		case packet := <-t.sendCh:
			if err := t.sendpacket(packet); err != nil {
				log.Warnf("sendpacket() Err:%s", err)
				return
			}

		case <-t.doneChan():
			return
		}
	}
}

// -
// sendpacket()
// -
func (t *TcpStruct) sendpacket(p *packet.PacketStruct) (err error) {
	if t == nil {
		return errors.ErrNetNilMethodPointer(log.Err(""))
	}

	var raw []byte
	if raw, err = p.ToByte(); err != nil {
		return errors.ErrNetParameter(log.Errf("Err:%s", err))

	} else if l, err := t.conn.Write(raw); err != nil {
		if err != io.EOF {
			return errors.ErrNetChannelError(log.Errf("TCP Write() Error:%s", err))
		}
		return errors.ErrNetChannelDown(log.Errf("TCP Write() Channel Closed"))
	} else if l != len(raw) {
		return errors.ErrNetChannelError(log.Errf("TCP Write() Lenth Error:%d != %d", l, len(raw)))
	}

	p.DebugPacket("Send()")
	log.Debugf("TCP Send %v", raw)

	return nil
}

// -
// sendclose()
// tries to send a last ditch effort close packet (soft close)
// letting the other side know that the link is closing out side of the IP protocol
// -
func (t *TcpStruct) sendclose() {
	if t == nil {
		return
	}

	if t.conn != nil {
		_, err := t.conn.Write(closePacketByte)

		if err != nil {
			log.Warnf("Err:%s", err)
		}
	}
}

// -
// sendtestpacket()
// -
func (t *TcpStruct) sendtestpacket(raw []byte) (err error) {
	if t == nil {
		return
	}

	log.Debugf("TCP RAW Send:%v", raw)

	if l, err := t.conn.Write(raw); err != nil {
		if err != io.EOF {
			return errors.ErrNetChannelError(log.Errf("TCP RAW Write() Error:%s, Closing Connection", err))
		}
	} else if l != len(raw) {
		return errors.ErrNetChannelError(log.Errf("TCP RAW Write() length Error:%d != %d", l, len(raw)))
	}

	return nil
}
