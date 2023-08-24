package tcp

import (
	"io"

	"github.com/seanmcadam/octovpn/octolib/log"
)

// Send()
func (t *TcpStruct) Send(buf []byte) (err error) {

	go func(buf []byte) {
		t.sendch <- buf
	}(buf)
	return err

}

func (t *TcpStruct) goSend() {

	defer t.emptysend()

	for {
		select {
		case buf := <-t.sendch:
			l, err := t.conn.Write(buf)
			if err != nil {
				if err != io.EOF {
					log.Errorf("TCPCli Write():%s", err)
				}
				return
			}

			if l != len(buf) {
				log.Errorf("TCP Write() Length Error:%d != %d", l, len(buf))
				return
			}

		case <-t.Closech:
			return
		default:
		}
	}

}

func (t *TcpStruct) emptysend() {
	for {
		select {
		case <-t.sendch:
		default:
			return
		}
	}
	close(t.sendch)
}
