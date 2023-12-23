package pinger

import "github.com/seanmcadam/counter"

type Ping counter.Count
type Pong counter.Count
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
