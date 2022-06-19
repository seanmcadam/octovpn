package packet

import (
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/seanmcadam/octovpn/octolib"
	"github.com/seanmcadam/octovpn/pinger"
)

//type headerString string
//const UDPHeaderSignature headerString = "-OctoVPN-UDP-1"
//const TCPHeaderSignature headerString = "-OctoVPN-TCP-1"

type ProtoHeader struct {
	//	Signature headerString
	ID      uint64
	Payload interface{}
}

var counterHeaderIDChan chan uint64

//
//
//
func init() {
	counterHeaderIDChan = octolib.RunGoCounter64()
	gob.Register(ProtoHeader{})
}

// func NewProtoHeader(sig headerString, payload interface{}) (p *ProtoHeader, e error) {
func NewProtoHeader(payload interface{}) (p *ProtoHeader, e error) {

	switch payload.(type) {
	case *ConnFrame:
	case ConnFrame:
	case *pinger.Ping:
	case pinger.Ping:
	case *pinger.Pong:
	case pinger.Pong:
	default:
		return nil, errors.New(fmt.Sprintf("Invalid payload type:%t", payload))
	}

	p = &ProtoHeader{
		//		Signature: sig,
		ID:      <-counterHeaderIDChan,
		Payload: payload,
	}
	return p, e
}
