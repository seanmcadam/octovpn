package netconn

import (
	"fmt"
	"io"
)


func (nc *NetConnStruct)GetReadChan()(ch chan []byte, err error){
	if nc.isClosed(){
		return nil, fmt.Errorf("connection closed")
	}
	return nc.recvch, nil
}

func (nc *NetConnStruct)Read()(b []byte, err error){
	if nc.isClosed(){
		return nil, fmt.Errorf("connection closed")
	}

	b = <- nc.recvch
	return b, nil
}

// goRead()
// Asyncronously read the net.Conn
// If conn closes then return
// If
func (nc *NetConnStruct)goRead(){
	defer nc.close()

	for {
			buffer := make([]byte, 2048) // Create a buffer to read data
			n, err := nc.conn.Read(buffer)
			if err != nil {
				if err != io.EOF{
				fmt.Println("Error reading:", err)
				}
				return
			}

			buffer = buffer[:n]

			nc.clock.Lock()
			if nc.isClosed() {
				return
			}
			nc.recvch <- buffer
			nc.clock.Unlock()
		}
	}
