package tcpsrv

import (
	"fmt"
	"net"
	"time"

	"github.com/seanmcadam/ctx"
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/internal/interfaces"
	"github.com/seanmcadam/octovpn/internal/layers/network/tcp"
	"github.com/seanmcadam/octovpn/internal/msg"
	"github.com/seanmcadam/octovpn/internal/msgnode"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/instance"
)

//
// func New(ctx *ctx.Ctx, config *settings.ConnectionStruct) (interfaces.Node, error) {
// func New(ctx *ctx.Ctx, config *settings.ConnectionStruct) (interfaces.Node, error) {
// func (tcpserver *TcpServerStruct) GetParentRecvCh() chan interfaces.MsgInterface {
// func (tcpclient *TcpServerStruct) GetChildRecvCh() chan interfaces.MsgInterface {
// func (tcpclient *TcpServerStruct) GetParentMsgHandlerFn() func(interfaces.MsgInterface) {
// func (tcpclient *TcpServerStruct) GetChildMsgHandlerFn() func(interfaces.MsgInterface) {
// func (tcpclient *TcpServerStruct) GetInstanceName() *instance.InstanceName {
//
// func (tcpclient *TcpServerStruct) setState(state msg.MsgState) {
// func (tcpclient *TcpServerStruct) notice(notice msg.MsgNotice) {
// func (tcpclient *TcpServerStruct) Cancel() {
// func (tcpclient *TcpServerStruct) closed() bool {
// func (tcpclient *TcpServerStruct) HandleParentPacket(m *msg.PacketStruct) {
// func (tcpclient *TcpServerStruct) HandleChildPacket(m *msg.PacketStruct) {
// func (tcpclient *TcpServerStruct) goNewClient() {
//

// -
// TcpServerStruct
// -
type TcpServerStruct struct {
	cx                  *ctx.Ctx
	me                  *instance.InstanceName
	config              *settings.ConnectionStruct
	address             string
	tcplistener         *net.TCPListener
	tcpaddr             *net.TCPAddr
	state               msg.MsgState
	parentCh            chan interfaces.MsgInterface
	parentPacketHandler *msg.PacketHandleTable
	childCh             chan interfaces.MsgInterface
	childPacketHandler  *msg.PacketHandleTable
	childStateHandler   *msg.StateHandleTable
	childNoticeHandler  *msg.NoticeHandleTable
}

var inst *instance.Instance

func init() {
	inst = instance.New()
}

func New(ctx *ctx.Ctx, config *settings.ConnectionStruct) (msgnode.Node, error) {
	return new(ctx, config)
}

func new(ctx *ctx.Ctx, config *settings.ConnectionStruct) (tcpserver *TcpServerStruct, err error) {

	tcpserver = &TcpServerStruct{
		cx:                  ctx,
		me:                  inst.Next(),
		config:              config,
		address:             fmt.Sprintf("%s:%d", config.Host, config.Port),
		tcplistener:         nil,
		tcpaddr:             nil,
		state:               msg.StateNONE,
		parentCh:            make(chan interfaces.MsgInterface),
		childCh:             make(chan interfaces.MsgInterface),
		childStateHandler:   msg.NewStateHandler(),
		childNoticeHandler:  msg.NewNoticeHandler(),
		parentPacketHandler: msg.NewPacketHandler(),
		childPacketHandler:  msg.NewPacketHandler(),
	}

	tcpserver.childPacketHandler.AddHandle(packet.SIG_CONN_32_RAW, tcpserver.HandleChildPacket)
	tcpserver.childPacketHandler.AddHandle(packet.SIG_CONN_32_PACKET, tcpserver.HandleChildPacket)
	tcpserver.childPacketHandler.AddHandle(packet.SIG_CONN_32_PING, tcpserver.HandleChildPacket)
	tcpserver.childPacketHandler.AddHandle(packet.SIG_CONN_32_PONG, tcpserver.HandleChildPacket)
	tcpserver.childPacketHandler.AddHandle(packet.SIG_CONN_32_AUTH, tcpserver.HandleChildPacket)
	tcpserver.childPacketHandler.AddHandle(packet.SIG_CONN_32_ID, tcpserver.HandleChildPacket)
	tcpserver.childPacketHandler.AddHandle(packet.SIG_CONN_64_RAW, tcpserver.HandleChildPacket)
	tcpserver.childPacketHandler.AddHandle(packet.SIG_CONN_64_PACKET, tcpserver.HandleChildPacket)
	tcpserver.childPacketHandler.AddHandle(packet.SIG_CONN_64_PING, tcpserver.HandleChildPacket)
	tcpserver.childPacketHandler.AddHandle(packet.SIG_CONN_64_PONG, tcpserver.HandleChildPacket)
	tcpserver.childPacketHandler.AddHandle(packet.SIG_CONN_64_AUTH, tcpserver.HandleChildPacket)
	tcpserver.childPacketHandler.AddHandle(packet.SIG_CONN_64_ID, tcpserver.HandleChildPacket)

	tcpserver.parentPacketHandler.AddHandle(packet.SIG_CONN_32_RAW, tcpserver.HandleParentPacket)
	tcpserver.parentPacketHandler.AddHandle(packet.SIG_CONN_32_PACKET, tcpserver.HandleParentPacket)
	tcpserver.parentPacketHandler.AddHandle(packet.SIG_CONN_32_PING, tcpserver.HandleParentPacket)
	tcpserver.parentPacketHandler.AddHandle(packet.SIG_CONN_32_PONG, tcpserver.HandleParentPacket)
	tcpserver.parentPacketHandler.AddHandle(packet.SIG_CONN_64_RAW, tcpserver.HandleParentPacket)
	tcpserver.parentPacketHandler.AddHandle(packet.SIG_CONN_64_PACKET, tcpserver.HandleParentPacket)
	tcpserver.parentPacketHandler.AddHandle(packet.SIG_CONN_64_PING, tcpserver.HandleParentPacket)
	tcpserver.parentPacketHandler.AddHandle(packet.SIG_CONN_64_PONG, tcpserver.HandleParentPacket)

	// Recheck this each time, the IP could change or rotate
	tcpserver.tcpaddr, err = net.ResolveTCPAddr(string(tcpserver.config.Proto), tcpserver.address)
	if err != nil {
		return nil, fmt.Errorf("ResolveTCPAddr Failed:%s", err)
	}

	tcpserver.tcplistener, err = net.ListenTCP(string(tcpserver.config.Proto), tcpserver.tcpaddr)
	if err != nil {
		return nil, fmt.Errorf("ListenTCP Failed:%s", err)
	}

	return tcpserver, err

}

// -
// goRun()
// Loop on
// 	Establish Connection
// 	Start Send and Recv Goroutines
// 	Monitor reset request
// -

func (tcpserver *TcpServerStruct) goRun() {

	if tcpserver == nil {
		log.ErrorStack("Nil Method Pointer")
		return
	}

	defer tcpserver.Cancel()

	for {
		var tcpconnclosech chan interface{}

		select {
		case conn := <-t.tcpconnch:
			log.Debugf("New incoming TCP Server Connection")
			t.msgnode.AddChildNode(conn)
			// t.link.AddLinkStateCh(conn.Link())
			// go t.goTcpStart(conn)

		case <-tcpconnclosech:
			continue

		case <-tcpserver.cx.DoneChan():
			return

		}
	}
}

// -
//
// -
func (tcpserver *TcpServerStruct) goNewClient() {

TCPFOR:
	for {
		var err error
		var conn *net.TCPConn

		// Dial it and keep trying forever
		conn, err = net.DialTCP(string(tcpserver.config.Proto), nil, tcpserver.tcpaddr)

		if err != nil || conn == nil {
			log.Warnf("connection failed %s: %s, wait", tcpserver.address, err)
			tcpserver.setState(msg.StateNOLINK)
			time.Sleep(1 * time.Second)
			continue TCPFOR
		}

		log.Info("New TCP Connection")

		tcpconn := tcp.New(tcpserver.cx.NewWithCancel(), conn)
		if tcpconn == nil {
			log.Fatal("tcpconn == nil")
		}
		return
	}
}

// -
//
// -
func (tcpserver *TcpServerStruct) GetParentRecvCh() chan interfaces.MsgInterface {
	if tcpserver == nil {
		return nil
	}
	return tcpserver.parentCh
}

// -
//
// -
func (tcpserver *TcpServerStruct) GetChildRecvCh() chan interfaces.MsgInterface {
	if tcpserver == nil {
		return nil
	}
	return tcpserver.childCh
}

// -
//
// -
func (tcpserver *TcpServerStruct) GetParentMsgHandlerFn() func(interfaces.MsgInterface) {
	if tcpserver == nil {
		log.FatalStack("Nil pointer")
	}

	fn := func(m interfaces.MsgInterface) {
		if tcpserver == nil {
			log.FatalStack("Nil pointer")
		}
		data := m.Data()
		switch p := data.(type) {
		case *msg.PacketStruct:
			tcpclient.parentPacketHandler.Run(p.Packet.Sig(), p)
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
//
// -
func (tcpserver *TcpServerStruct) GetChildMsgHandlerFn() func(interfaces.MsgInterface) {
	if tcpserver == nil {
		log.FatalStack("Nil pointer")
	}

	fn := func(m interfaces.MsgInterface) {
		if tcpserver == nil {
			log.FatalStack("Nil pointer")
		}
		data := m.Data()
		switch p := data.(type) {
		case *msg.PacketStruct:
			tcpserver.childPacketHandler.Run(p.Packet.Sig(), p)
		case *msg.StateStruct:
			tcpserver.childStateHandler.Run(p)
		case *msg.NoticeStruct:
			tcpserver.childNoticeHandler.Run(p)
		default:
			log.FatalfStack("Default Reached on %v", p)
		}
	}

	return fn
}

// -
//
// -
func (tcpserver *TcpServerStruct) GetInstanceName() *instance.InstanceName {
	if tcpserver == nil {
		return nil
	}
	return tcpserver.me
}

// -
//
// -
func (tcp *TcpServerStruct) setState(state msg.MsgState) {
	if tcp == nil {
		log.FatalStack()
		return
	}
	s := msg.NewState(tcp.me, state)
	tcp.state = state
	tcpserver.parentCh <- s
}

// -
//
// -
func (tcpserver *TcpServerStruct) notice(notice msg.MsgNotice) {
	if tcpserver == nil {
		log.FatalStack()
		return
	}
	n := msg.NewNotice(tcpserver.me, notice)
	tcpserver.parentCh <- n

}

// -
//
// -
func (tcpserver *TcpServerStruct) Cancel() {
	if tcpserver == nil {
		return
	}
	tcpserver.setState(msg.StateNOLINK)
	tcpserver.notice(msg.NoticeCLOSED)
	tcpserver.cx.Cancel()
}

// -
//
// -
func (tcpserver *TcpServerStruct) closed() bool {
	if tcpserver == nil {
		return true
	}
	return tcpserver.cx.Done()
}

// -
//
// -
func (tcpserver *TcpServerStruct) HandleParentPacket(m *msg.PacketStruct) {
	tcpserver.childCh <- m
}

// -
//
// -
func (tcpserver *TcpServerStruct) HandleChildPacket(m *msg.PacketStruct) {
	tcpserver.parentCh <- m
}
