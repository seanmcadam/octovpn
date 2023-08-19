package netconn

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
)


func (nc *NetConnStruct)GetReadChan()(ch chan *NetPacketStruct, err error){
	if nc.isClosed(){
		return nil, fmt.Errorf("connection closed")
	}
	return nc.recvch, nil
}

func (nc *NetConnStruct)Read()(b []byte, err error){
	if nc.isClosed(){
		return nil, fmt.Errorf("connection closed")
	}

	b = (<- nc.recvch).Packet()
	return b, nil
}

// goRead()
// Asyncronously read the net.Conn
// If conn closes then return
// If
func (nc *NetConnStruct)goRead(){
	defer nc.close()


	for {
			var n int
			var err error
			lengthbyte := make([]byte,2)
			n, err = nc.conn.Read(lengthbyte)
			if err != nil {
				if err == io.EOF {
					log.Print("Network Connection EOF...")
				} else {
					log.Printf("Error reading:%s", err)
				}
				return
			}
			if n != 2 {
				log.Printf("read length wrong:2 != %d", n)
				return
			}
			length := binary.BigEndian.Uint16(lengthbyte)

			countbyte := make([]byte,4)
			n, err = nc.conn.Read(countbyte)
			if err != nil {
				if err == io.EOF {
					log.Print("Network Connection EOF...")
				} else {
					log.Printf("Error reading:%s", err)
				}
				return
			}
			if n != 4 {
				log.Printf("read count wrong:4 != %d", n)
				return
			}
			count := binary.BigEndian.Uint32(countbyte)

			data := make([]byte,int(length-4))
			n, err = nc.conn.Read(data)
			if err != nil {
				if err == io.EOF {
					log.Print("Network Connection EOF...")
				} else {
					log.Printf("Error reading:%s", err)
				}
				return
			}
			if n != int(length-4){
				log.Printf("read data wrong:%d != %d", n, int(length-4))
				return
			}

			buffer := &NetPacketStruct{
				length: length,
				count: count,
				packet: data,
			}
			
			nc.clock.Lock()
			if nc.isClosed() {
				return
			}
			nc.recvch <- buffer
			nc.clock.Unlock()
		}
	}
