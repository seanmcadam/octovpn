package packet

import (
	"fmt"

	"github.com/seanmcadam/octovpn/octolib/errors"
	log "github.com/seanmcadam/loggy"
)

type AuthPacketActionType uint8
type AuthErrNoTextProvided error
type AuthErrShortPacketSize error
type AuthErrPacketSizeMismatch error
type AuthErrChallengeGenaration error

const AuthChallenge AuthPacketActionType = 0x01
const AuthResponse AuthPacketActionType = 0x02
const AuthAccept AuthPacketActionType = 0x03
const AuthReject AuthPacketActionType = 0x04
const AuthError AuthPacketActionType = 0xFF
const AuthPacketActionSize = 1
const AuthPacketMinSize = PacketSize8Size + AuthPacketActionSize

type AuthPacket struct {
	//
	// The first byte represent the packet size
	// but it is converted to PacketSizeType with is a uint16
	pSize  PacketSizeType
	action AuthPacketActionType
	//
	// text is optional and only used for Challange and Response
	// The challange is a random string
	// The response is the MD5SUM of the random string + the secret
	text []byte
}

func NewAuth(action AuthPacketActionType, text ...string) (ap *AuthPacket, err error) {

	// TODO
	// Check MAX size too
	// Errors...

	var pSize = AuthPacketMinSize
	var t []byte

	switch action {
	case AuthChallenge:
		fallthrough
	case AuthResponse:
		fallthrough
	case AuthError:
		if (len(text) == 0) || (len(text[0]) == 0) {
			return nil, AuthErrNoTextProvided(fmt.Errorf("No Text Provided for %s", action))
		}
		pSize += PacketSizeType(len(text[0]))
		t = []byte(text[0])
	}

	ap = &AuthPacket{
		pSize:  pSize,
		action: action,
		text:   t,
	}
	return ap, err
}

// MakeAuth([]byte)(*AuthPacket, error)
// Takes a byte representation of the Auth Packet and converts it to AuthPacket
func MakeAuth(raw []byte) (p *AuthPacket, err error) {
	if raw == nil {
		log.ErrorStack("Nil Parameter")
		return nil, errors.ErrPacketBadParameter(log.Errf("Nil raw data"))
	}

	if len(raw) < int(AuthPacketMinSize) {
		return nil, AuthErrShortPacketSize(fmt.Errorf("Len:%d", len(raw)))
	}
	rawsize := len(raw)

	pSize := PacketSizeType(BtoU8(raw))
	raw = raw[1:]

	if pSize != PacketSizeType(rawsize) {
		return nil, AuthErrPacketSizeMismatch(fmt.Errorf("Recv: %d pSize:%d", rawsize, pSize))
	}

	action := AuthPacketActionType(BtoU8(raw))
	raw = raw[1:]

	switch action {
	case AuthChallenge:
		fallthrough
	case AuthResponse:
		fallthrough
	case AuthError:
		if pSize < 3 {
			return nil, AuthErrShortPacketSize(fmt.Errorf("No Text field - Len:%d", len(raw)))
		}
	}

	p = &AuthPacket{
		pSize:  pSize,
		action: action,
		text:   raw,
	}

	return p, err
}

func (ap *AuthPacket) Size() PacketSizeType {
	if ap == nil {
		return PacketSizeTypeERROR
	}

	return ap.pSize
}

func (ap *AuthPacket) Action() AuthPacketActionType {
	if ap == nil {
		return AuthError
	}

	return ap.action
}

func (ap *AuthPacket) Text() []byte {
	if ap == nil {
		return []byte("Nil Method Pointer, will this actually be seen?")
	}

	return ap.text
}

func (ap *AuthPacket) ToByte() (raw []byte, err error) {
	if ap == nil {
		return nil, errors.ErrPacketNilMethodPointer(log.Errf(""))
	}

	if ap.pSize > 256 {
		return nil, errors.ErrPacketBadParameter(log.Errf("psize:%d", ap.pSize))
	}

	raw = append(raw, uint8(ap.pSize))
	raw = append(raw, uint8(ap.action))
	if ap.pSize > AuthPacketMinSize {
		switch ap.action {
		case AuthChallenge:
			fallthrough
		case AuthResponse:
			fallthrough
		case AuthError:
			if len(ap.text) != 0 {
				raw = append(raw, ap.text...)
			} else {
				return nil, errors.ErrPacketBadParameter(log.Errf("no text provided"))
			}
		default:
			return nil, errors.ErrPacketBadParameter(log.Errf("default reached"))
		}
	}
	return raw, nil
}

func (a AuthPacketActionType) String() (ret string) {
	switch a {
	case AuthChallenge:
		return "CHALLENGE"
	case AuthResponse:
		return "RESPONSE"
	case AuthAccept:
		return "ACCEPT"
	case AuthReject:
		return "REJECT"
	case AuthError:
		return "ERROR"
	default:
		log.FatalfStack("Unhandled AuthPacketActionType:%d", uint8(a))
	}

	return ret
}
