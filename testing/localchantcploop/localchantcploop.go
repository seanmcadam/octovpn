package main

import (
	"math/rand"

	"github.com/seanmcadam/octovpn/internal/channel/loopchan"
	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func main() {
	cx := ctx.NewContext()

	srv, cli, err := loopchan.NewTcpChanLoop(cx)
	if err != nil {
		log.FatalfStack("NewTcpChanLoop Err:%s", err)
	}

	srvdatach := createDataGenerator(cx)
	clidatach := createDataGenerator(cx)

	for {
		select {
		case srvdata := <-srvdatach:
			p, err := packet.NewPacket(packet.SIG_CHAN_32_RAW, srvdata, counter.MakeCounter32(uint32(1)))
			if err != nil {
				log.FatalfStack("NewPacket Err:%s", err)
			}
			srv.Send(p)
		case clidata := <-clidatach:
			p, err := packet.NewPacket(packet.SIG_CHAN_32_RAW, clidata, counter.MakeCounter32(uint32(2)))
			if err != nil {
				log.FatalfStack("NewPacket Err:%s", err)
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
	size := 1 + rand.Intn(1024) // Generate random size between 1 and 1024
	data := make([]byte, size)
	_, err := rand.Read(data)
	if err != nil {
		log.Fatalf("Error generating random data:", err)
	}
	return data
}
