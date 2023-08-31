package main

import (
	"math/rand"

	"github.com/seanmcadam/octovpn/internal/chanconn/loopconn"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/internal/packet/packetconn"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

const randomDataSize = 1

func main() {
	cx := ctx.NewContext()

	srv, cli, err := loopconn.NewTcpConnLoop(cx)
	if err != nil {
		log.FatalStack("NewTcpConnLoop Err:%s", err)
	}

	srvdatach := createDataGenerator(cx)
	clidatach := createDataGenerator(cx)

	for {
		select {
		case srvdata := <-srvdatach:
			p, err := packetconn.NewPacket(packet.CONN_TYPE_RAW, srvdata)
			if err != nil {
				log.FatalStack("NewPacket Err:%s", err)
			}
			srv.Send(p)
		case clidata := <-clidatach:
			p, err := packetconn.NewPacket(packet.CONN_TYPE_RAW, clidata)
			if err != nil {
				log.FatalStack("NewPacket Err:%s", err)
			}
			cli.Send(p)
		case srvrecv := <-srv.RecvChan():
			_ = srvrecv
		case clirecv := <-cli.RecvChan():
			_ = clirecv
		}

	}

}

func createDataGenerator(ctx *ctx.Ctx) (ch chan []byte) {
	ch = make(chan []byte)
	go func() {
		for {
			data := generateRandomData()
			ch <- data
		}
	}()
	return ch
}

func generateRandomData() []byte {
	size := 1 + rand.Intn(randomDataSize) // Generate random size between 1 and 1024
	data := make([]byte, size)
	_, err := rand.Read(data)
	if err != nil {
		log.Fatalf("Error generating random data:", err)
	}
	return data
}
