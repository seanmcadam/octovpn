package auth

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"time"

	"github.com/seanmcadam/octovpn/internal/link"
	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type AuthErrChallengeGenaration error
type AuthErrAuthRejected error
type AuthErrRecvError error

const AuthChallengeMinSize = 1
const AuthChallengeMaxSize = 128
const AuthMaxPacketCount = 100
const AuthRetryDelay = 50 * time.Millisecond
const AuthTimeout = 5 * time.Second

type AuthStateType uint8

const AuthStateStart AuthStateType = 0x00
const AuthStateUnauthenticated AuthStateType = 0x01
const AuthStateChallenge AuthStateType = 0x02
const AuthStateAuthenticated AuthStateType = 0x04
const AuthStateError AuthStateType = 0xFF

//
// Auth Module is used to authenticate the same levels on each side of a connection
// Auth has a built in link object, and will toggle between DOWN -> CHALLENGE -> AUTH -> UP in a normal sequence
// The module sends a challenge packet to the far side and waits on a response
// The response is validated and the system either is Authenticated or Unauthenticated
// The module also handles the mirror response
// A Challeng packet is received and a response is sent back, and waits for either Accept or Rejcet
//
//

type AuthStruct struct {
	cx          *ctx.Ctx
	link        *link.LinkStateStruct
	localstate  AuthStateType
	remotestate AuthStateType
	secret      string
	random      string
	md5sum      string
	sendcount   uint8
	recvcount   uint8
	localtimer  *time.Timer
	remotetimer *time.Timer
	sendch      chan *packet.AuthPacket
	recvch      chan *packet.AuthPacket
}

// -
//
// -
func NewAuthStruct(ctx *ctx.Ctx, secret string) (as *AuthStruct, err error) {
	as = &AuthStruct{
		cx:          ctx,
		link:        link.NewLinkState(ctx),
		localstate:  AuthStateStart,
		remotestate: AuthStateStart,
		sendcount:   0,
		recvcount:   0,
		secret:      secret,
		random:      "",
		md5sum:      "",
		localtimer:  time.NewTimer(AuthTimeout),
		remotetimer: time.NewTimer(AuthTimeout),
		sendch:      make(chan *packet.AuthPacket),
		recvch:      make(chan *packet.AuthPacket),
	}

	as.link.Down()

	go as.goAuth()

	return as, nil
}

func (as *AuthStruct) GetRecvCh() (ch chan<- *packet.AuthPacket) {
	return as.recvch
}

func (as *AuthStruct) GetSendCh() (ch <-chan *packet.AuthPacket) {
	return as.sendch
}

func (as *AuthStruct) Link() *link.LinkStateStruct {
	return as.link
}


// -
//
// -
func (as *AuthStruct) goAuth() {

	as.resetLocalTimer(AuthRetryDelay)
	defer as.cx.Cancel()
	defer as.localtimer.Stop()
	defer as.remotetimer.Stop()
	defer log.Debugf("Auth Exit Sent:%d, Recv:%d", as.sendcount, as.recvcount)

MAINLOOP:
	for {

		//
		// Bail out with excessive Auth Traffic
		//
		if as.recvcount > AuthMaxPacketCount || as.sendcount > AuthMaxPacketCount {
			log.Errorf("Max Packets reached Send:%d, Recv:%d", as.recvcount, as.sendcount)
			return
		}

		log.Debug("Auth Loop Start")

		select {
		case <-as.cx.DoneChan():
			log.Debug("Auth Done, exiting")
			return
		case ap := <-as.recvch:
			log.Debug("Auth Recv")

			as.recvcount++

			switch ap.Action() {
			case packet.AuthChallenge:
				//
				// Remote Challange Packet - Link DN,
				// Set Link Down
				// Reply
				// Set Remote Timer
				//
				as.setRemoteState(AuthStateChallenge)
				as.recvChallengePacket(ap)
				as.resetRemoteTimer(AuthTimeout)

			case packet.AuthResponse:
				//
				// Receive Response Packet from a For a local Challange
				//
				valid := as.recvResponsePacket(ap)
				if valid {
					// Its Valid
					as.setLocalState(AuthStateAuthenticated)
					as.sendPacket(packet.AuthAccept, "")
					if as.remotestate == AuthStateAuthenticated {
						as.link.Up()
						as.resetCounters()
					} else {
						as.link.Auth()
					}
				} else {
					as.link.Down()
					as.setLocalState(AuthStateUnauthenticated)
					as.sendPacket(packet.AuthReject, "")
					as.resetLocalTimer(AuthRetryDelay)
				}

			case packet.AuthAccept:
				//
				// Recieve Remote Accept from a Remote Challenge
				//
				as.setRemoteState(AuthStateAuthenticated)
				if as.localstate == AuthStateAuthenticated {
					as.link.Up()
					as.resetCounters()
				}

			case packet.AuthReject:
				//
				// Recieve Remote Reject - Reset to initial conditions
				//
				as.link.Down()
				as.setRemoteState(AuthStateUnauthenticated)
				as.setLocalState(AuthStateStart)
				as.resetRemoteTimer(AuthTimeout)
				as.resetLocalTimer(AuthRetryDelay)

			case packet.AuthError:
				log.Errorf("Received AuthError:%s", ap.Text())
				fallthrough
			default:
				as.link.Down()
				return
			}

		case <-as.remotetimer.C:
			log.Debug("Auth Remote Timer")
			switch as.remotestate {
			case AuthStateUnauthenticated:
				fallthrough
			case AuthStateStart:
				//
				// No Challenge Packet Recieved after startup
				// We are just hanging...
				// Waited for timeout...
				//
				if as.localstate < AuthStateUnauthenticated {
					as.resetRemoteTimer(AuthTimeout)
					continue MAINLOOP
				}

				// No packets received from the other side...
				return

			case AuthStateChallenge:
				//
				// I receieved a challenge, but no follow up...
				//
				if as.localstate < AuthStateAuthenticated {
					as.resetRemoteTimer(AuthTimeout)
					continue MAINLOOP
				}

				// No more packets received from the other side...
				return

			case AuthStateAuthenticated:
				//
				// Remote side says good to go... what is the hold up?
				//
				if as.localstate < AuthStateAuthenticated {
					as.resetRemoteTimer(AuthTimeout)
					continue MAINLOOP
				}

				return

			default:
				as.link.Down()
				log.Errorf("Remote Unhandled Auth State:%02X, bailing out", as.remotestate)
				return
			}

		case <-as.localtimer.C:
			log.Debug("Auth Local Timer")
			switch as.localstate {
			case AuthStateStart:
				as.random = generateChallengePhrase()
				as.md5sum = generateMD5Sum(as.secret, as.random)
				as.setLocalState(AuthStateChallenge)
				as.link.Chal()
				as.sendPacket(packet.AuthChallenge, as.random)
				as.resetLocalTimer(AuthTimeout)

			case AuthStateChallenge:
				//
				// Timeout getting a response
				// and Have gotten no packets from the other side
				//
				if as.remotestate < AuthStateUnauthenticated {
					log.Warn("Local Auth Challenge Timeout, restarting process")
					return
				}
				// Else restart on this side
				as.setLocalState(AuthStateStart)
				as.resetLocalTimer(AuthRetryDelay)

			case AuthStateAuthenticated:
				// Do nothing - we are good

			case AuthStateUnauthenticated:
				log.Warn("Local Auth Unauthenticated, restarting process")
				as.link.Down()
				as.setLocalState(AuthStateStart)
				as.resetLocalTimer(AuthRetryDelay)

			default:
				as.link.Down()
				log.Errorf("Local Unhandled Auth State:%02X, bailing out", as.localstate)
				return
			}
		}
	}
}

// -
//
// -
func (as *AuthStruct) sendPacket(a packet.AuthPacketActionType, text string) {
	as.sendcount++
	p, err := packet.NewAuth(a, text)
	if err != nil {
		log.Errorf("NewAuth Err:%s", err)
		as.setLocalState(AuthStateError)
		as.setRemoteState(AuthStateError)
		as.cx.Cancel()
		return
	}
	as.sendch <- p
}

// -
//
// -
func (as *AuthStruct) resetLocalTimer(t time.Duration) {
	as.localtimer.Stop()
	as.localtimer.Reset(t)
}

// -
//
// -
func (as *AuthStruct) resetRemoteTimer(t time.Duration) {
	as.remotetimer.Stop()
	as.remotetimer.Reset(t)
}

// -
//
// -
func (as *AuthStruct) recvChallengePacket(recv *packet.AuthPacket) {
	challenge := recv.Text()
	md5sum := generateMD5Sum(as.secret, string(challenge))
	as.sendPacket(packet.AuthResponse, md5sum)
}

// -
//
// -
func (as *AuthStruct) recvResponsePacket(recv *packet.AuthPacket) bool {
	return string(recv.Text()) == as.md5sum
}

// -
// setLocalState()
// -
func (as *AuthStruct) setLocalState(s AuthStateType) {
	as.localstate = s
}

// -
// setRemoteState()
// -
func (as *AuthStruct) setRemoteState(s AuthStateType) {
	as.remotestate = s
}

// -
// resetCounters()
// -
func (as *AuthStruct) resetCounters() {
	as.sendcount = 0
	as.recvcount = 0
}

// -
// generateMD5Sum()
// Used at start up to create a unique challange phrase to combine with the pass phrase
// -
func generateMD5Sum(secret string, random string) (md5sum string) {

	buf := []byte(secret + random)
	hash := md5.Sum(buf)
	md5sum = hex.EncodeToString(hash[:])

	return md5sum
}

// -
// generateChannengePhrase()
// Used at start up to create a unique challange phrase to combine with the pass phrase
// -
func generateChallengePhrase() (result string) {
	min := AuthChallengeMinSize
	max := AuthChallengeMaxSize

	var sizeb [8]byte
	_, err := rand.Read(sizeb[:])
	if err != nil {
		log.Fatalf("ranf.Read() Size err:%s", err)
	}

	num := binary.BigEndian.Uint64(sizeb[:])
	size := int(num%(uint64(max-min+1))) + min

	random := make([]byte, size)
	_, err = rand.Read(random[:])
	if err != nil {
		log.Fatalf("ranf.Read() Random err:%s", err)
	}

	result = base64.StdEncoding.EncodeToString(random)

	return result
}
