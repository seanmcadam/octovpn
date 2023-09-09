package main

import (
	"math/rand"
	"time"

	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/internal/layer/chanconn/loopconn"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

const randomDataSize = 1300

func main() {
	cx := ctx.NewContext()
	c32 := counter.NewCounter32(cx)

	srv, cli, err := loopconn.NewTcpConnLoop(cx)
	if err != nil {
		log.FatalStack("NewTcpConnLoop Err:%s", err)
	}

	time.Sleep(time.Second)

	//
	// Create Random Data and send it on each side
	//
	srvdatach := rawDataGenerator(cx)
	clidatach := rawDataGenerator(cx)
	//clidatach := make(chan []byte)

	//
	//
	for {
		select {
		case srvdata := <-srvdatach:
			p, err := packet.NewPacket(packet.SIG_CONN_32_RAW, srvdata, <-c32.GetCountCh())
			if err != nil {
				log.FatalStack("NewPacket Err:%s", err)
			}
			p.DebugPacket("CONN TCP Send RAW SRV: ")
			srv.Send(p)
		case clidata := <-clidatach:
			p, err := packet.NewPacket(packet.SIG_CONN_32_RAW, clidata, <-c32.GetCountCh())
			if err != nil {
				log.FatalStack("NewPacket Err:%s", err)
			}
			p.DebugPacket("CONN TCP Send RAW CLI: ")
			cli.Send(p)
		case srvrecv := <-srv.RecvChan():
			srvrecv.DebugPacket("CONN TCP RECV SRV: ")
			_ = srvrecv
		case clirecv := <-cli.RecvChan():
			clirecv.DebugPacket("CONN TCP RECV CLI: ")
			_ = clirecv
		}

	}

}

func rawDataGenerator(ctx *ctx.Ctx) (ch chan []byte) {
	ch = make(chan []byte)
	go func() {
		for {
			data := generateRandomData()
			ch <- data
			//return
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return ch
}

func generateRandomData() []byte {
	size := 1 + rand.Intn(randomDataSize) // Generate random size between 1 and 1024
	data := make([]byte, size)
	_, err := rand.Read(data)
	if err != nil {
		log.FatalfStack("Error generating random data:", err)
	}
	return data
}
