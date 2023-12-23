package tcp

import (
	"net"
	"sync"

	"github.com/seanmcadam/ctx"
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/internal/interfaces"
	"github.com/seanmcadam/octovpn/internal/msg"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/instance"
)

var inst *instance.Instance

// -
//
// -
type TcpStruct struct {
	cx                  *ctx.Ctx
	closeonce           sync.Once
	me                  *instance.InstanceName
	conn                *net.TCPConn
	parentCh            chan interfaces.MsgInterface
	sendCh              chan *packet.PacketStruct
	state               msg.MsgState
	parentPacketHandler *msg.PacketHandleTable
}

// -
//
// -
func init() {
	inst = instance.New()
}

// -
//
// -
func New(ctx *ctx.Ctx, conn *net.TCPConn) interfaces.Node {
	return new(ctx, conn)
}

// -
//
// -
func new(ctx *ctx.Ctx, conn *net.TCPConn) (tcp *TcpStruct) {
	if ctx == nil || conn == nil {
		return nil
	}

	//name := fmt.Sprintf(">>TCP<<[%s<->%s]", conn.LocalAddr(), conn.RemoteAddr())
	log.Debugf("Local Addr %s <-> %s", conn.LocalAddr(), conn.RemoteAddr())

	tcp = &TcpStruct{
		cx: ctx,
		//closeonce:           new(sync.Once),
		me:                  inst.Next(),
		conn:                conn,
		parentCh:            make(chan interfaces.MsgInterface),
		sendCh:              make(chan *packet.PacketStruct),
		state:               msg.StateNONE,
		parentPacketHandler: msg.NewPacketHandler(),
	}

	tcp.parentPacketHandler.AddHandle(packet.SIG_CONN_32_RAW, tcp.HandleParentPacket)
	tcp.parentPacketHandler.AddHandle(packet.SIG_CONN_32_PACKET, tcp.HandleParentPacket)
	tcp.parentPacketHandler.AddHandle(packet.SIG_CONN_32_PING, tcp.HandleParentPacket)
	tcp.parentPacketHandler.AddHandle(packet.SIG_CONN_32_PONG, tcp.HandleParentPacket)
	tcp.parentPacketHandler.AddHandle(packet.SIG_CONN_32_AUTH, tcp.HandleParentPacket)
	tcp.parentPacketHandler.AddHandle(packet.SIG_CONN_32_ID, tcp.HandleParentPacket)
	tcp.parentPacketHandler.AddHandle(packet.SIG_CONN_64_RAW, tcp.HandleParentPacket)
	tcp.parentPacketHandler.AddHandle(packet.SIG_CONN_64_PACKET, tcp.HandleParentPacket)
	tcp.parentPacketHandler.AddHandle(packet.SIG_CONN_64_PING, tcp.HandleParentPacket)
	tcp.parentPacketHandler.AddHandle(packet.SIG_CONN_64_PONG, tcp.HandleParentPacket)
	tcp.parentPacketHandler.AddHandle(packet.SIG_CONN_64_AUTH, tcp.HandleParentPacket)
	tcp.parentPacketHandler.AddHandle(packet.SIG_CONN_64_ID, tcp.HandleParentPacket)

	go tcp.goSend()
	go tcp.goRecv()
	go tcp.setState(msg.StateCONNECTED)
	return tcp
}

// -
//
// -
func (tcp *TcpStruct) GetInstanceName() *instance.InstanceName {
	if tcp == nil {
		log.FatalStack("Nil pointer")
	}
	return tcp.me
}

// -
//
// -
func (tcp *TcpStruct) GetParentRecvCh() chan interfaces.MsgInterface {
	if tcp == nil {
		log.FatalStack("Nil pointer")
	}
	return tcp.parentCh
}

// -
//
// -
func (tcp *TcpStruct) HandleParentPacket(packet *msg.PacketStruct) {
	if tcp == nil {
		log.FatalStack("Nil pointer")
	}
	tcp.sendCh <- packet.Packet
}

// -
//
// -
func (tcp *TcpStruct) GetParentMsgHandlerFn() func(interfaces.MsgInterface) {
	if tcp == nil {
		log.FatalStack("Nil pointer")
	}

	fn := func(m interfaces.MsgInterface) {
		if tcp == nil {
			log.FatalStack("Nil pointer")
		}
		data := m.Data()
		switch p := data.(type) {
		case *msg.PacketStruct:
			tcp.parentPacketHandler.Run(p.Packet.Sig(), p)
		case *msg.StateStruct:
			log.FatalfStack("Parent State Msg recvived %v", p)
		case *msg.NoticeStruct:
			log.FatalfStack("Parent Notice Msg recvived %v", p)
		default:
			log.FatalfStack("Default Reached on %v", p)
		}
	}

	return fn

}

// -
// No children
// -
func (tcp *TcpStruct) GetChildRecvCh() chan interfaces.MsgInterface {
	log.FatalStack("TCP does not have children")
	return nil
}

// -
// No children
// -
func (tcp *TcpStruct) GetChildMsgHandlerFn() func(interfaces.MsgInterface) {
	log.FatalStack("TCP does not have children")
	return nil
}

// -
//
// -
func (tcp *TcpStruct) ReceiveHandler(data interfaces.MsgInterface) {
	if tcp == nil {
		return
	}

	if data.FromName() == tcp.me {

	} else {
		log.FatalStack("Got non-parent message")
	}

	switch msg := data.(type) {
	case *msg.NoticeStruct:
		log.Fatalf("[%s] received Notice:%s", *tcp.me, msg)
	case *msg.StateStruct:
		log.Fatalf("[%s] received State:%s", *tcp.me, msg)
	case *msg.PacketStruct:
		tcp.sendCh <- msg.Packet
	default:
		log.Fatalf("default reached %T", data)
		return
	}

}

// -
//
// -
func (tcp *TcpStruct) RemoteAddrString() string {
	if tcp == nil {
		log.FatalStack()
	}
	return tcp.conn.RemoteAddr().String()
}

// -
//
// -
func (tcp *TcpStruct) setState(state msg.MsgState) {
	if tcp == nil {
		log.FatalStack()
	}

	tcp.state = state
	s := msg.NewState(tcp.me, state)
	tcp.parentCh <- s
}

// -
//
// -
func (tcp *TcpStruct) notice(notice msg.MsgNotice) {
	if tcp == nil {
		log.FatalStack()
	}
	n := msg.NewNotice(tcp.me, notice)
	tcp.parentCh <- n
}

// -
//
// -
func (t *TcpStruct) doneChan() <-chan struct{} {
	if t == nil {
		return nil
	}
	return t.cx.DoneChan()
}

// -
// Just call it one time, but multiple places can call it
// -
func (tcp *TcpStruct) close() {
	if tcp == nil {
		log.Errorf("TCP Close() called with nill pointer")
		return
	}

	c := func() {
		log.Debugf("MsgBus Closeing: %s", tcp.me)

		tcp.setState(msg.StateNOLINK)
		tcp.notice(msg.NoticeCLOSED)

		tcp.cx.Cancel()
		close(tcp.parentCh)
		close(tcp.sendCh)
	}

	tcp.closeonce.Do(c)
}
