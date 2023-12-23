package tcpcli

import (
	"fmt"
	"net"
	"time"

	"github.com/seanmcadam/ctx"
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/internal/interfaces"
	"github.com/seanmcadam/octovpn/internal/layers/network/tcp"
	"github.com/seanmcadam/octovpn/internal/msg"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/internal/settings"
	"github.com/seanmcadam/octovpn/octolib/instance"
)


//
// func New(ctx *ctx.Ctx, config *settings.ConnectionStruct) (interfaces.Node, error) {
// func New(ctx *ctx.Ctx, config *settings.ConnectionStruct) (interfaces.Node, error) {
// func (tcpclient *TcpClientStruct) GetParentRecvCh() chan interfaces.MsgInterface {
// func (tcpclient *TcpClientStruct) GetChildRecvCh() chan interfaces.MsgInterface {
// func (tcpclient *TcpClientStruct) GetParentMsgHandlerFn() func(interfaces.MsgInterface) {
// func (tcpclient *TcpClientStruct) GetChildMsgHandlerFn() func(interfaces.MsgInterface) {
// func (tcpclient *TcpClientStruct) GetInstanceName() *instance.InstanceName {
//
// func (tcpclient *TcpClientStruct) setState(state msg.MsgState) {
// func (tcpclient *TcpClientStruct) notice(notice msg.MsgNotice) {
// func (tcpclient *TcpClientStruct) Cancel() {
// func (tcpclient *TcpClientStruct) closed() bool {
// func (tcpclient *TcpClientStruct) HandleParentPacket(m *msg.PacketStruct) {
// func (tcpclient *TcpClientStruct) HandleChildPacket(m *msg.PacketStruct) {
// func (tcpclient *TcpClientStruct) goNewClient() {
//
//

// -
// TcpClientStruct
// Maintains a connection to the specified address. If the connection drops, it is reestablished
// -
type TcpClientStruct struct {
	cx                  *ctx.Ctx
	me                  *instance.InstanceName
	config              *settings.ConnectionStruct
	address             string
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

func New(ctx *ctx.Ctx, config *settings.ConnectionStruct) (interfaces.Node, error) {
	return new(ctx, config)
}

func new(ctx *ctx.Ctx, config *settings.ConnectionStruct) (tcpclient *TcpClientStruct, err error) {

	tcpclient = &TcpClientStruct{
		cx:                  ctx,
		me:                  inst.Next(),
		config:              config,
		address:             fmt.Sprintf("%s:%d", config.Host, config.Port),
		tcpaddr:             nil,
		state:               msg.StateNONE,
		parentCh:            make(chan interfaces.MsgInterface),
		childCh:             make(chan interfaces.MsgInterface),
		childStateHandler:   msg.NewStateHandler(),
		childNoticeHandler:  msg.NewNoticeHandler(),
		parentPacketHandler: msg.NewPacketHandler(),
		childPacketHandler:  msg.NewPacketHandler(),
	}

	tcpclient.childPacketHandler.AddHandle(packet.SIG_CONN_32_RAW, tcpclient.HandleChildPacket)
	tcpclient.childPacketHandler.AddHandle(packet.SIG_CONN_32_PACKET, tcpclient.HandleChildPacket)
	tcpclient.childPacketHandler.AddHandle(packet.SIG_CONN_32_PING, tcpclient.HandleChildPacket)
	tcpclient.childPacketHandler.AddHandle(packet.SIG_CONN_32_PONG, tcpclient.HandleChildPacket)
	tcpclient.childPacketHandler.AddHandle(packet.SIG_CONN_32_AUTH, tcpclient.HandleChildPacket)
	tcpclient.childPacketHandler.AddHandle(packet.SIG_CONN_32_ID, tcpclient.HandleChildPacket)
	tcpclient.childPacketHandler.AddHandle(packet.SIG_CONN_64_RAW, tcpclient.HandleChildPacket)
	tcpclient.childPacketHandler.AddHandle(packet.SIG_CONN_64_PACKET, tcpclient.HandleChildPacket)
	tcpclient.childPacketHandler.AddHandle(packet.SIG_CONN_64_PING, tcpclient.HandleChildPacket)
	tcpclient.childPacketHandler.AddHandle(packet.SIG_CONN_64_PONG, tcpclient.HandleChildPacket)
	tcpclient.childPacketHandler.AddHandle(packet.SIG_CONN_64_AUTH, tcpclient.HandleChildPacket)
	tcpclient.childPacketHandler.AddHandle(packet.SIG_CONN_64_ID, tcpclient.HandleChildPacket)

	tcpclient.parentPacketHandler.AddHandle(packet.SIG_CONN_32_RAW, tcpclient.HandleParentPacket)
	tcpclient.parentPacketHandler.AddHandle(packet.SIG_CONN_32_PACKET, tcpclient.HandleParentPacket)
	tcpclient.parentPacketHandler.AddHandle(packet.SIG_CONN_32_PING, tcpclient.HandleParentPacket)
	tcpclient.parentPacketHandler.AddHandle(packet.SIG_CONN_32_PONG, tcpclient.HandleParentPacket)
	tcpclient.parentPacketHandler.AddHandle(packet.SIG_CONN_64_RAW, tcpclient.HandleParentPacket)
	tcpclient.parentPacketHandler.AddHandle(packet.SIG_CONN_64_PACKET, tcpclient.HandleParentPacket)
	tcpclient.parentPacketHandler.AddHandle(packet.SIG_CONN_64_PING, tcpclient.HandleParentPacket)
	tcpclient.parentPacketHandler.AddHandle(packet.SIG_CONN_64_PONG, tcpclient.HandleParentPacket)

	// Do an initial check and fail if it fails
	// Recheck this each time, the IP could change or rotate
	tcpclient.tcpaddr, err = net.ResolveTCPAddr(string(tcpclient.config.Proto), tcpclient.address)
	if err != nil {
		return nil, err
	}

	return tcpclient, err
}

// -
//
// -
func (tcpclient *TcpClientStruct) goNewClient() {

TCPFOR:
	for {
		var err error
		var conn *net.TCPConn

		// Dial it and keep trying forever
		conn, err = net.DialTCP(string(tcpclient.config.Proto), nil, tcpclient.tcpaddr)

		if err != nil || conn == nil {
			log.Warnf("connection failed %s: %s, wait", tcpclient.address, err)
			tcpclient.setState(msg.StateNOLINK)
			time.Sleep(1 * time.Second)
			continue TCPFOR
		}

		log.Info("New TCP Connection")

		tcpconn := tcp.New(tcpclient.cx.NewWithCancel(), conn)
		if tcpconn == nil {
			log.Fatal("tcpconn == nil")
		}
		return
	}
}

// -
//
// -
func (tcpclient *TcpClientStruct) GetParentRecvCh() chan interfaces.MsgInterface {
	if tcpclient == nil {
		return nil
	}
	return tcpclient.parentCh
}

// -
//
// -
func (tcpclient *TcpClientStruct) GetChildRecvCh() chan interfaces.MsgInterface {
	if tcpclient == nil {
		return nil
	}
	return tcpclient.childCh
}

// -
//
// -
func (tcpclient *TcpClientStruct) GetParentMsgHandlerFn() func(interfaces.MsgInterface) {
	if tcpclient == nil {
		log.FatalStack("Nil pointer")
	}

	fn := func(m interfaces.MsgInterface) {
		if tcpclient == nil {
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
func (tcpclient *TcpClientStruct) GetChildMsgHandlerFn() func(interfaces.MsgInterface) {
	if tcpclient == nil {
		log.FatalStack("Nil pointer")
	}

	fn := func(m interfaces.MsgInterface) {
		if tcpclient == nil {
			log.FatalStack("Nil pointer")
		}
		data := m.Data()
		switch p := data.(type) {
		case *msg.PacketStruct:
			tcpclient.childPacketHandler.Run(p.Packet.Sig(), p)
		case *msg.StateStruct:
			tcpclient.childStateHandler.Run(p)
		case *msg.NoticeStruct:
			tcpclient.childNoticeHandler.Run(p)
		default:
			log.FatalfStack("Default Reached on %v", p)
		}
	}

	return fn
}

// -
//
// -
func (tcpclient *TcpClientStruct) GetInstanceName() *instance.InstanceName {
	if tcpclient == nil {
		return nil
	}
	return tcpclient.me
}

// -
//
// -
func (tcpclient *TcpClientStruct) setState(state msg.MsgState) {
	if tcpclient == nil {
		log.FatalStack()
		return
	}
	tcpclient.state = state
	s := msg.NewState(tcpclient.me, state)
	tcpclient.parentCh <- s
}

// -
//
// -
func (tcpclient *TcpClientStruct) notice(notice msg.MsgNotice) {
	if tcpclient == nil {
		log.FatalStack()
		return
	}
	n := msg.NewNotice(tcpclient.me, notice)
	tcpclient.parentCh <- n

}

// -
//
// -
func (tcpclient *TcpClientStruct) Cancel() {
	if tcpclient == nil {
		return
	}
	tcpclient.setState(msg.StateNOLINK)
	tcpclient.notice(msg.NoticeCLOSED)
	tcpclient.cx.Cancel()
}

// -
//
// -
func (tcpclient *TcpClientStruct) closed() bool {
	if tcpclient == nil {
		return true
	}
	return tcpclient.cx.Done()
}

// -
//
// -
func (tcpclient *TcpClientStruct) HandleParentPacket(m *msg.PacketStruct) {
	tcpclient.childCh <- m
}

// -
//
// -
func (tcpclient *TcpClientStruct) HandleChildPacket(m *msg.PacketStruct) {
	tcpclient.parentCh <- m
}
