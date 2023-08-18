package netconn

import (
	"fmt"
	"log"
	"net"
	"sync"
)

const BufferLength = 5

type NetConnStruct struct {
	clock sync.Mutex
	conn net.Conn
	sendch chan []byte
	recvch chan []byte
	closech chan interface{}
}

func NewNetConn(conn net.Conn) (nc *NetConnStruct) {

	if conn == nil {
		log.Fatalf("NewNetConn() failed with conn = nil")
	}

	nc = &NetConnStruct{
		conn: conn,
		sendch: make(chan []byte, BufferLength),
		recvch: make(chan []byte, BufferLength),
		closech: make(chan interface{}),
	}

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
		case buf := <- nc.sendch:
			l, err := nc.conn.Write(buf)
			if err != nil{
				fmt.Printf("nc.conn.Write() error:%s", err)
				return
			}
			if l != len(buf) {
				fmt.Printf("nc.conn.Write() length mismatch: buf%d, sent:%d", len(buf), l)
				return
			}
		}
	}
}