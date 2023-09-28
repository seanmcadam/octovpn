package tcp

import (
	"net"

	"github.com/seanmcadam/octovpn/internal/msgbus"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/instance"
	"github.com/seanmcadam/octovpn/octolib/log"
)

var inst *instance.Instance

type TcpStruct struct {
	cx     *ctx.Ctx
	me     msgbus.MsgTarget
	parent msgbus.MsgTarget
	state  msgbus.MsgState
	msgbus *msgbus.MsgBus
	conn   *net.TCPConn
	sendch chan *packet.PacketStruct
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

	me := msgbus.MsgTarget(inst.Next())
	tcp = &TcpStruct{
		cx:     ctx,
		msgbus: mb,
		me:     me,
		parent: parent,
		conn:   conn,
		sendch: make(chan *packet.PacketStruct),
		state:  msgbus.StateNONE,
	}

	tcp.msgbus.ReceiveHandler(tcp.me, tcp.receiveHandler)

	tcp.setState(msgbus.StateCONNECTED)
	go tcp.goRecv()
	go tcp.goSend()
	return tcp
}

func (tcp *TcpStruct) InstanceName() msgbus.MsgTarget {
	return tcp.me
}

func (tcp *TcpStruct) RemoteAddrString() string {
	return tcp.conn.RemoteAddr().String()
}

// -
//
// -
func (tcp *TcpStruct) setState(state msgbus.MsgState) {
	tcp.state = state
	tcp.msgbus.SendState(tcp.me, tcp.parent, tcp.state)
}

// -
//
// -
func (t *TcpStruct) receiveHandler(data ...interface{}) {

	if len(data) == 0 {
		log.Errorf("[%s]:no data", t.me)
	}

	switch datatype := data[0].(type) {
	case *packet.PacketStruct:
		t.sendch <- data[0].(*packet.PacketStruct)
		//t.sendpacket(data[0].(*packet.PacketStruct))
	case *msgbus.MsgNotice:
	case *msgbus.MsgState:
	default:
		log.Fatalf("[%s]:default reached: %T", t.me, datatype)
	}
}
