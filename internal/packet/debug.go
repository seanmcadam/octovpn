package packet

import (
	"fmt"

	log "github.com/seanmcadam/loggy"
)

func (p *PacketStruct) DebugPacket(prefix string) {
	if p == nil {
		return
	}

	debug := fmt.Sprintf("\n\t----------- PACKET ------------------\n")
	debug += fmt.Sprintf("\t%s\n", prefix)
	debug += fmt.Sprintf("\tSIG:%s\n", p.Sig())
	debug += fmt.Sprintf("\tSIZE:%d WIDTH%d\n", p.Size(), p.Width())
	if p.Sig().Count() {
		debug += fmt.Sprintf("\tCOUNTER:%d\n", p.Count().Uint())
	}
	if p.Sig().Ping() {
		debug += fmt.Sprintf("\tPING:%d\n", p.Ping().Uint())
	}
	if p.Sig().Pong() {
		debug += fmt.Sprintf("\tPONG:%d\n", p.Pong().Uint())
	}
	if p.Sig().Data() {
		debug += fmt.Sprintf("\tPck:%v IP4:%v IP6:%v Eth:%v Rout:%v Auth:%v ID:%v Raw:%v\n",
			p.Sig().Packet(),
			p.Sig().IPV4(),
			p.Sig().IPV6(),
			p.Sig().Eth(),
			p.Sig().Router(),
			p.Sig().Auth(),
			p.Sig().ID(),
			p.Sig().Raw())
	}
	b, _ := p.ToByte()
	checksum := calculateCRC32Checksum(b)
	debug += fmt.Sprintf("\tCheckSum:%08x\n", checksum)

	log.Debug(debug)
}
