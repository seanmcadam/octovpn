package tcp

import (
	"bytes"
	"io"

	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

//-
//
//-
func (t *TcpStruct) RecvChan() <-chan *packet.PacketStruct {
	if t == nil || t.recvch == nil {
		return nil
	}

	return t.recvch
}

//-
// Run while connection is running
// Exit when closed
//-
func (t *TcpStruct) goRecv() {
	if t == nil {
		return
	}

	defer t.emptyrecv()
	defer t.Cancel()

	for {
		var buffer bytes.Buffer

		tmp := make([]byte, 2048)

		for {
			//
			// Load the receive buffer
			//
			n, err := t.conn.Read(tmp)
			if err != nil {
				if err == io.EOF {
					log.Errorf("Read() connection closed %s", t.conn.RemoteAddr())
				} else {
					log.Errorf("Read() Error:%s on %s", err, t.conn.RemoteAddr())
				}
				return
			}

			log.Debugf("Raw Recv len:%d", n)
			buffer.Write(tmp[:n])

			//
			// Does the buffer have enough data to assemble a packet?
			//
			sig, length, err := packet.ReadPacketBuffer(buffer.Bytes()[:6])

			//
			// Error checking types here
			//
			if err != nil {
				log.Errorf("MakePacket() Err:%s on %s", err, t.conn.RemoteAddr())
				return
			}

			//
			// Only receive CONN layer packets here
			//
			if !sig.ConnLayer() {
				log.Errorf("Bad SIG Layer Received:%s, on %s", sig, t.conn.RemoteAddr())
				return
			}

			//
			// Is there enough data?
			//
			if buffer.Len() < int(length) {
				log.Warnf("Not Enough Buffer Data %d < %d", buffer.Len(), int(length))
				continue
			}

			//
			// Extract a packet
			//
			p, err := packet.MakePacket(buffer.Next(int(length)))
			if err != nil {
				log.Errorf("MakePacket() Err:%s on %s", err, t.conn.RemoteAddr())
				return
			}

			if p == nil {
				log.Errorf("MakePacket() returned Nil Packet")
			}

			t.recvch <- p
		}
	}
}

// -
// Clean up the recvch before closing
// -
func (t *TcpStruct) emptyrecv() {
	if t == nil {
		return
	}

	for {
		select {
		case <-t.recvch:
		default:
			close(t.recvch)
			return
		}
	}
}
