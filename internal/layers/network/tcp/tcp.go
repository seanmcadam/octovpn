package tcp

import (
	"net"

	"github.com/seanmcadam/octovpn/internal/msgbus"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/instance"
	"github.com/seanmcadam/octovpn/octolib/log"
)

var inst *Instance

type TcpStruct struct {
	cx     *ctx.Ctx
	target msgbus.MsgTarget
	msgbus *msgbus.MsgBus
	conn   *net.TCPConn
	//sendch chan *packet.PacketStruct
	//recvch chan *packet.PacketStruct
}

func init() {
	inst = instance.New()
}

func NewTCP(ctx *ctx.Ctx, mb *msgbus.MsgBus, parent msgbus.MsgTarget, conn *net.TCPConn) (tcp *TcpStruct) {
	if ctx == nil || conn == nil {
		return nil
	}

	//name := fmt.Sprintf(">>TCP<<[%s<->%s]", conn.LocalAddr(), conn.RemoteAddr())
	log.Debugf("Local Addr %s <-> %s", conn.LocalAddr(), conn.RemoteAddr())

	tcp = &TcpStruct{
		cx:     ctx,
		target: msgbus.MsgTarget(inst.Next()),
		msgbus: mb,
		conn:   conn,
	}

	return tcp
}

func (t *TcpStruct) Run() {
	if t == nil {
		return
	}
	//t.msgbus.SetState(msgbus.StateCONNECTED)
	go t.goRecv()
	go t.goSend()
}

func (t *TcpStruct) RemoteAddrString() string {
	return t.conn.RemoteAddr().String()
}
