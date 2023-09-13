package pinger

import "github.com/seanmcadam/octovpn/internal/counter"

type Ping counter.Counter
type Pong counter.Counter
type PingWidth uint8

const PingWidth32 PingWidth = 32
const PingWidth64 PingWidth = 64

type PingerStruct interface {
	NewPong([]byte) (Pong, error)
	TurnOn()
	TurnOff()
	Width() PingWidth
	RecvPong(Pong)
	GetPingChan() <- chan Ping
}
