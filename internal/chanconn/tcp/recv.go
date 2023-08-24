package tcp

import (
	"io"

	"github.com/seanmcadam/octovpn/octolib/log"
)

// Recv()
func (t *TcpStruct) Recv() (buf []byte, err error) {

	buf = <-t.recvch

	return buf, err
}


// Run while connection is running
// Exit when closed
func (t *TcpStruct) goRecv() {
	defer t.emptyrecv()

	for{
		buf := make([]byte, 2048)

		l, err := t.conn.Read(buf)
		if err != nil {
			if err != io.EOF{
				log.Errorf("TCP Read() Error:%s", err)
			}
			return 
		}

		buf = buf[:l]
	
		select {
		case t.recvch <- buf:
		case <-t.Closech:
			return
		default:
		}
	}

}

func (t *TcpStruct) emptyrecv() {
	for {
		select {
		case <-t.recvch:
		default:
			return
		}
	}
	close(t.recvch)
}
