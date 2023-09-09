package tcp

import (
	"bytes"
	"io"

	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (t *TcpStruct) RecvChan() <-chan *packet.PacketStruct {
	if t == nil || t.recvch == nil {
		log.Debugf("TCP Recv Nil")
		return nil
	}

	log.Debugf("TCP Recv state:%s", t.link.GetState())
	return t.recvch
}

// Run while connection is running
// Exit when closed
func (t *TcpStruct) goRecv() {
	if t == nil {
		return
	}

	defer t.emptyrecv()
	defer t.link.Down()

	for {
		var buffer bytes.Buffer

		tmp := make([]byte, 2048)

		for {
			n, err := t.conn.Read(tmp)
			if err != nil {
				if err == io.EOF {
					log.Errorf("TCP Read() connection closed")
				} else {
					log.Errorf("TCP Read() Error:%s", err)
				}
				return
			}

			log.Debugf("TCP Raw Recv len:%d", n)
			buffer.Write(tmp[:n])

			if t.closed() {
				return
			}

			sig, length, err := packet.ReadPacketBuffer(buffer.Bytes()[:6])
			//
			// Error checking types here
			//
			if err != nil {
				log.Errorf("TCP MakePacket() Err:%s", err)
				continue
			}

			if !sig.ConnLayer() {
				log.Errorf("Bad Layer Received")
				continue
			}

			if buffer.Len() < int(length) {
				continue
			}

			p, err := packet.MakePacket(buffer.Next(int(length)))
			if err != nil {
				log.Errorf("MakePacket Err:%s", err)
				continue
			}

			if p == nil {
				log.FatalStack("Got Nil Packet")
			}

			log.Debugf("TCP Recv:%v", p)
			t.recvch <- p

			if t.closed() {
				return
			}
		}
	}
}

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
