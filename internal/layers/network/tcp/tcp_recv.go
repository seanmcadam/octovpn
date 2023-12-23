package tcp

import (
	"bytes"
	"io"

	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/internal/msg"
	"github.com/seanmcadam/octovpn/internal/packet"
)

// -
// Run while connection is running
// Exit when closed
// -
func (t *TcpStruct) goRecv() {
	if t == nil {
		return
	}

	defer t.close()

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

			buffer.Write(tmp[:n])

			//
			// Does the buffer have enough data to assemble a packet?
			//
			sig, length, err := packet.ReadPacketBuffer(buffer.Bytes()[:6])

			// log.Debugf("RecvBuffer:%v",buffer.Bytes()[:n])
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
			newpacketbuf := buffer.Next(int(length))
			// log.Debugf("Raw TCP Recv:%v", newpacketbuf)
			p, err := packet.MakePacket(newpacketbuf)
			if err != nil {
				log.Errorf("MakePacket() Err:%s on %s", err, t.conn.RemoteAddr())
				return
			}

			if p == nil {
				log.Errorf("MakePacket() returned Nil Packet")
			}

			if p.Sig().Close() {
				log.Debug("TCP received SOFT CLOSE")
				return
			}

			packet := msg.NewPacket(t.me, p)
			log.Debug("TCP Recv PAcket %v", packet)
			t.parentCh <- packet
		}
	}
}
