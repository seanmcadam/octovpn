package tcp

import (
	"bytes"
	"io"

	"github.com/seanmcadam/octovpn/internal/link"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (t *TcpStruct) RecvChan() <-chan *packet.PacketStruct {

	if t == nil {
		log.FatalStack("nil TcpStruct")
		return nil
	}
	if t.recvch == nil {
		log.Error("Nil recvch pointer")
		return nil
	}

	return t.recvch
}

// Run while connection is running
// Exit when closed
func (t *TcpStruct) goRecv() {
	defer t.emptyrecv()
	defer t.link.ToggleState(link.LinkStateDown)

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

			if !sig.ConnLayer(){
				log.Errorf("Bad Layer Received")
				continue
			}

			if buffer.Len() < int(length){
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

			t.recvch <- p

			if t.closed() {
				return
			}
		}
	}
}

func (t *TcpStruct) emptyrecv() {
	for {
		select {
		case <-t.recvch:
		default:
			close(t.recvch)
			return
		}
	}
}
