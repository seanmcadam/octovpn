package msg

import (
	log "github.com/seanmcadam/loggy"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/instance"
)

type PacketHandleTable struct {
	handler map[packet.PacketSigType]func(*PacketStruct)
}

func NewPacketHandler() (pht *PacketHandleTable) {
	pht = &PacketHandleTable{
		handler: map[packet.PacketSigType]func(*PacketStruct){},
	}
	pht.handler[packet.SIG_ERROR] = emptyPacketHandle
	pht.handler[packet.SIG_CONN_START] = emptyPacketHandle
	pht.handler[packet.SIG_CHAN_START] = emptyPacketHandle
	pht.handler[packet.SIG_SITE_START] = emptyPacketHandle
	pht.handler[packet.SIG_ROUTE_START] = emptyPacketHandle
	pht.handler[packet.SIG_CONN_CLOSE] = emptyPacketHandle
	pht.handler[packet.SIG_CHAN_CLOSE] = emptyPacketHandle
	pht.handler[packet.SIG_SITE_CLOSE] = emptyPacketHandle
	pht.handler[packet.SIG_ROUTE_CLOSE] = emptyPacketHandle
	pht.handler[packet.SIG_ROUTE] = emptyPacketHandle
	pht.handler[packet.SIG_SITE] = emptyPacketHandle
	pht.handler[packet.SIG_CHAN] = emptyPacketHandle
	pht.handler[packet.SIG_CONN] = emptyPacketHandle
	pht.handler[packet.SIG_ROUTE_AUTH] = emptyPacketHandle
	pht.handler[packet.SIG_ROUTE_RAW] = emptyPacketHandle
	pht.handler[packet.SIG_ROUTE_ETH] = emptyPacketHandle
	pht.handler[packet.SIG_ROUTE_IPV4] = emptyPacketHandle
	pht.handler[packet.SIG_ROUTE_IPV6] = emptyPacketHandle
	pht.handler[packet.SIG_ROUTE_ROUTER] = emptyPacketHandle
	pht.handler[packet.SIG_SITE_32_RAW] = emptyPacketHandle
	pht.handler[packet.SIG_SITE_32_PACKET] = emptyPacketHandle
	pht.handler[packet.SIG_SITE_32_AUTH] = emptyPacketHandle
	pht.handler[packet.SIG_SITE_32_ID] = emptyPacketHandle
	pht.handler[packet.SIG_SITE_32_ACK] = emptyPacketHandle
	pht.handler[packet.SIG_SITE_32_NAK] = emptyPacketHandle
	pht.handler[packet.SIG_SITE_64_RAW] = emptyPacketHandle
	pht.handler[packet.SIG_SITE_64_PACKET] = emptyPacketHandle
	pht.handler[packet.SIG_SITE_64_AUTH] = emptyPacketHandle
	pht.handler[packet.SIG_SITE_64_ID] = emptyPacketHandle
	pht.handler[packet.SIG_SITE_64_ACK] = emptyPacketHandle
	pht.handler[packet.SIG_SITE_64_NAK] = emptyPacketHandle
	pht.handler[packet.SIG_CHAN_32_RAW] = emptyPacketHandle
	pht.handler[packet.SIG_CHAN_32_PACKET] = emptyPacketHandle
	pht.handler[packet.SIG_CHAN_32_PING] = emptyPacketHandle
	pht.handler[packet.SIG_CHAN_32_PONG] = emptyPacketHandle
	pht.handler[packet.SIG_CHAN_32_ACK] = emptyPacketHandle
	pht.handler[packet.SIG_CHAN_32_NAK] = emptyPacketHandle
	pht.handler[packet.SIG_CHAN_64_RAW] = emptyPacketHandle
	pht.handler[packet.SIG_CHAN_64_PACKET] = emptyPacketHandle
	pht.handler[packet.SIG_CHAN_64_PING] = emptyPacketHandle
	pht.handler[packet.SIG_CHAN_64_PONG] = emptyPacketHandle
	pht.handler[packet.SIG_CHAN_64_ACK] = emptyPacketHandle
	pht.handler[packet.SIG_CHAN_64_NAK] = emptyPacketHandle
	pht.handler[packet.SIG_CONN_32_RAW] = emptyPacketHandle
	pht.handler[packet.SIG_CONN_32_PACKET] = emptyPacketHandle
	pht.handler[packet.SIG_CONN_32_PING] = emptyPacketHandle
	pht.handler[packet.SIG_CONN_32_PONG] = emptyPacketHandle
	pht.handler[packet.SIG_CONN_32_AUTH] = emptyPacketHandle
	pht.handler[packet.SIG_CONN_32_ID] = emptyPacketHandle
	pht.handler[packet.SIG_CONN_64_RAW] = emptyPacketHandle
	pht.handler[packet.SIG_CONN_64_PACKET] = emptyPacketHandle
	pht.handler[packet.SIG_CONN_64_PING] = emptyPacketHandle
	pht.handler[packet.SIG_CONN_64_PONG] = emptyPacketHandle
	pht.handler[packet.SIG_CONN_64_AUTH] = emptyPacketHandle
	pht.handler[packet.SIG_CONN_64_ID] = emptyPacketHandle

	return pht
}

func (nht *PacketHandleTable) Run(pst packet.PacketSigType, packet *PacketStruct) {
	if handlerFn, ok := nht.handler[pst]; ok {
		handlerFn(packet)
	} else {
		log.FatalfStack("Unknown Run Packet Type %08X", pst)
	}
}

func (nht *PacketHandleTable) AddHandle(pst packet.PacketSigType, fn func(*PacketStruct)) {
	nht.handler[pst] = fn
}

func (nht *PacketHandleTable) CallHandle(ps *PacketStruct) {
	nht.handler[ps.Packet.Sig()](ps)
}

type PacketStruct struct {
	Packet *packet.PacketStruct
	From   *instance.InstanceName
}

func NewPacket(from *instance.InstanceName, packet *packet.PacketStruct) (ps *PacketStruct) {
	ps = &PacketStruct{
		Packet: packet,
		From:   from,
	}
	return ps
}

func (p *PacketStruct) FromName() *instance.InstanceName {
	return p.From
}

func (n *PacketStruct) Data() interface{} {
	return n.Packet
}

func emptyPacketHandle(ps *PacketStruct) {
	log.ErrorfStack("Packet EmptyHandler From:%s Packet:%s", ps.From, ps.Packet)
}
