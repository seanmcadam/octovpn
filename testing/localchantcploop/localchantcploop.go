package main

import (
	"math/rand"
	"time"

	"github.com/seanmcadam/octovpn/internal/channel/loopchan"
	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func main() {
	cx := ctx.NewContext()
	c32 := counter.NewCounter32(cx)

	srv, cli, err := loopchan.NewTcpChanLoop(cx)
	if err != nil {
		log.FatalfStack("NewTcpChanLoop Err:%s", err)
	}

	srvdatach := createDataGenerator(cx)
	clidatach := createDataGenerator(cx)

	for {
		select {
		case srvdata := <-srvdatach:
			p, err := packet.NewPacket(packet.SIG_CHAN_32_RAW, srvdata, <-c32.GetCountCh())
			//p.DebugPacket("CHAN TCP SEND SRV: ")
			if err != nil {
				log.FatalfStack("NewPacket Err:%s", err)
			}
			srv.Send(p)
		case clidata := <-clidatach:
			p, err := packet.NewPacket(packet.SIG_CHAN_32_RAW, clidata, <-c32.GetCountCh())
			//p.DebugPacket("CHAN TCP RECV CLI: ")
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
	time.Sleep(10 * time.Millisecond)
	ch = make(chan []byte)
	go func() {
		for {
			data := generateRandomData()
			ch <- data
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return ch
}

func generateRandomData() []byte {
	size := 1 + rand.Intn(1500) // Generate random size between 1 and 1024
	data := make([]byte, size)
	_, err := rand.Read(data)
	if err != nil {
		log.Fatalf("Error generating random data:", err)
	}
	return data
}
