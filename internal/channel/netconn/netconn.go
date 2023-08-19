package netconn

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
)

const BufferLength = 5


type NetConnStruct struct {
	count chan uint32
	clock sync.Mutex
	conn net.Conn
	sendch chan *NetPacketStruct
	recvch chan *NetPacketStruct
	closech chan interface{}
}

func NewNetConn(conn net.Conn) (nc *NetConnStruct) {

	if conn == nil {
		log.Fatalf("NewNetConn() failed with conn = nil")
	}

	nc = &NetConnStruct{
		conn: conn,
		sendch: make(chan *NetPacketStruct, BufferLength),
		recvch: make(chan *NetPacketStruct, BufferLength),
		closech: make(chan interface{}),
	}
	nc.startCounter()

	return nc
}

func (nc *NetConnStruct)Run(){
	go nc.goRun()
	go nc.goRead()
}

func (nc *NetConnStruct)isClosed()bool{
	select{
	case x := <-nc.closech:
		return x == nil
	default:
		return false
	}
}

func (nc *NetConnStruct)goRun(){

	defer nc.close()

	for {
		select {
		case nps := <- nc.sendch:
			buf := make([]byte,nps.Length()+6)
			
			binary.BigEndian.PutUint16(buf[:2], nps.length)
			binary.BigEndian.PutUint32(buf[2:6], nps.count)
			buf = append(buf,nps.Packet()...)

			l, err := nc.conn.Write(buf)
			if err != nil{
				fmt.Printf("nc.conn.Write() error:%s", err)
				return
			}
			if l != len(buf) {
				fmt.Printf("nc.conn.Write() length mismatch: buf%d, sent:%d", len(buf), l)
				return
			}
		case <-nc.closech:
			return
		}
	}
}

func (nc *NetConnStruct)startCounter() {
	if nc == nil {
		panic("bad pointer")
	}

	c := make(chan uint32, 10)
	nc.count = c

	go func(ch chan uint32) {
		var counter uint32 = 1
		for {
			select {
			case ch <- counter:
				counter += 1
			case <-nc.closech:
				close(ch)
				return
			}
		}
	}(c)
}
