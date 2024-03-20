package connection

import (
	"testing"

	"github.com/seanmcadam/bufferpool"
	"github.com/seanmcadam/ctx"
	"github.com/seanmcadam/loggy"
	"github.com/seanmcadam/testlib/testnetpair"
)

func TestCompile(t *testing.T) {

}

func TestConnection(t *testing.T) {
	cx := ctx.New()
	conna, connb, err := testnetpair.NewNetworkTCPConnPair()

	pool := bufferpool.New()

	if err != nil {
		t.Errorf("network pair error:%s", err)
		return
	}

	a := Connection(cx, conna)
	b := Connection(cx, connb)

	sendbufa := pool.Get()
	sendbufb := pool.Get()

	sendbufa.Append([]byte("buf a"))
	sendbufb.Append([]byte("buf b"))

	a.Send(sendbufa)
	b.Send(sendbufb)

	recvbufa := <-a.RecvCh()
	if recvbufa == nil {
		t.Errorf("A Conn closed")
		return
	}
	recvbufb := <-b.RecvCh()
	if recvbufb == nil {
		t.Errorf("B Conn closed")
		return
	}

	loggy.Debugf("A:'%s'", recvbufa.Data())
	loggy.Debugf("B:'%s'", recvbufb.Data())

	recvbufa.ReturnToPool()
	recvbufb.ReturnToPool()

	return
}
