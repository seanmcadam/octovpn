package server

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"testing"

	"github.com/seanmcadam/octovpn/connection"
	"github.com/seanmcadam/octovpn/ctx"
	"github.com/seanmcadam/octovpn/octoconfig"
)

func TestTCPSetupTestCancel(t *testing.T) {

	ctx := ctx.NewContext()
	config := &octoconfig.ConfigListen{
		Protocol: octoconfig.UDP,
		IP:       "127.0.0.1",
		Port:     33333,
		MTU:      1400,
	}

	readChan := make(chan *connection.ConnFrame)
	listen, e := newListenUDP(ctx, config, readChan)

	if e != nil {
		t.Fatalf(fmt.Sprintf("newListenUDP() Test Failed: %s", e))
	}

	rUDPaddr, e := net.ResolveUDPAddr("udp", "127.0.0.1:33333")
	lUDPaddr, e := net.ResolveUDPAddr("udp", "127.0.0.1")
	udpConn, e := net.DialUDP("udp", lUDPaddr, rUDPaddr)

	buf := bufio.NewWriter(udpConn)

	encoder := gob.NewEncoder(buf)

	ping := NewPing()
	header := NewProtoHeader(udpHeaderSignature, ping)

	e = encoder.Encode(&header)
	if e != nil {
		panic("")
	}
	sz := buf.Size()
	t.Logf("buf sie:%d", sz)

	buf.Flush()

	test := <-readChan
	_ = test

	<-listen.ctx.DoneChan()
	listen.ctx.Cancel()
}
