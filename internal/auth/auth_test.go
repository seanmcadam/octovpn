package auth

import (
	"crypto/rand"
	"encoding/base64"
	"testing"
	"time"

	"github.com/seanmcadam/octovpn/internal/packet"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

//
// TODO - Add link state testing
//

func TestNewAuth_compile(t *testing.T) {
	cx := ctx.NewContext()
	_, _ = NewAuthStruct(cx, "")
	cx.Cancel()
}

// Auth Get to Local Authenticated
func TestNewAuth_test_local_auth(t *testing.T) {

	cx := ctx.NewContext()
	secret := generateRandomSecret()

	as, err := NewAuthStruct(cx, secret)
	if err != nil {
		t.Errorf("NewAuthStruct() Err:%s", err)
	}

	challenge := <-as.GetSendCh()
	if challenge.Action() != packet.AuthChallenge {
		t.Error("Expected Challenge Packet")
	}
	log.Debug("Challenged Recv")

	md5sum := generateMD5Sum(secret, string(challenge.Text()))

	p, err := packet.NewAuth(packet.AuthResponse, md5sum)
	if err != nil {
		t.Errorf("NewAuth() Err:%s", err)
	}
	log.Debugf("Auth Packet: %v", p)

	as.GetRecvCh() <- p

	reply := <-as.GetSendCh()
	if reply.Action() != packet.AuthAccept {
		t.Error("Expected Accept Packet")
	}
	log.Debugf("Result Packet: %v", reply)

	cx.Cancel()
}

// Auth Get to Remote Authenticated
func TestNewAuth_test_remote_auth(t *testing.T) {

	cx := ctx.NewContext()
	secret := generateRandomSecret()
	phrase := generateChallengePhrase()
	md5sum := generateMD5Sum(secret, phrase)

	as, err := NewAuthStruct(cx, secret)
	if err != nil {
		t.Errorf("NewAuthStruct() Err:%s", err)
	}

	// Ignore the first packet
	challenge := <-as.GetSendCh()
	if challenge.Action() != packet.AuthChallenge {
		t.Error("Expected Challenge Packet")
	}

	//
	// Send Challenge
	//
	p, err := packet.NewAuth(packet.AuthChallenge, phrase)
	if err != nil {
		t.Errorf("NewAuth() Err:%s", err)
	}
	as.GetRecvCh() <- p

	//
	// Get Response
	//
	response := <-as.GetSendCh()
	if response.Action() != packet.AuthResponse {
		t.Error("Expected Accpte Packet")
	}

	if string(response.Text()) != md5sum {
		t.Error("Expected matching md5sums")
	}

	cx.Cancel()
}

// Local Auth First, then Remote Auth
func TestNewAuth_test_local_remote_auth(t *testing.T) {

	cx := ctx.NewContext()
	secret := generateRandomSecret()

	remotephrase := generateChallengePhrase()
	remotemd5sum := generateMD5Sum(secret, remotephrase)

	as, err := NewAuthStruct(cx, secret)
	if err != nil {
		t.Errorf("NewAuthStruct() Err:%s", err)
	}

	//
	// Get Challenge
	//
	challenge := <-as.GetSendCh()
	if challenge.Action() != packet.AuthChallenge {
		t.Error("Expected Challenge Packet")
	}
	md5sum := generateMD5Sum(secret, string(challenge.Text()))

	//
	// Send Response
	//
	p, err := packet.NewAuth(packet.AuthResponse, md5sum)
	if err != nil {
		t.Errorf("NewAuth() Err:%s", err)
	}
	as.GetRecvCh() <- p

	//
	// Get Accept
	//
	reply := <-as.GetSendCh()
	if reply.Action() != packet.AuthAccept {
		t.Error("Expected Accept Packet")
	}

	//
	// Send Challenge
	//
	p, err = packet.NewAuth(packet.AuthChallenge, remotephrase)
	if err != nil {
		t.Errorf("NewAuth() Err:%s", err)
	}
	as.GetRecvCh() <- p

	//
	// Get Response
	//
	response := <-as.GetSendCh()
	if response.Action() != packet.AuthResponse {
		t.Error("Expected Response Packet")
	}

	if string(response.Text()) != remotemd5sum {
		t.Error("Expected matching md5sums")
	}

	//
	// Send Accept
	//
	p, err = packet.NewAuth(packet.AuthAccept)
	if err != nil {
		t.Errorf("NewAuth() Err:%s", err)
	}
	as.GetRecvCh() <- p

	cx.Cancel()
}

// Remote Auth First, then Local Auth
func TestNewAuth_test_remote_local_auth(t *testing.T) {

	cx := ctx.NewContext()
	secret := generateRandomSecret()

	remotephrase := generateChallengePhrase()
	remotemd5sum := generateMD5Sum(secret, remotephrase)

	as, err := NewAuthStruct(cx, secret)
	if err != nil {
		t.Errorf("NewAuthStruct() Err:%s", err)
	}

	UpCh := as.link.LinkUpCh()
	select {
	case <-UpCh:
		t.Error("Premature up")
	default:
	}

	//
	// Send Challenge
	//
	p, err := packet.NewAuth(packet.AuthChallenge, remotephrase)
	if err != nil {
		t.Errorf("NewAuth() Err:%s", err)
	}
	as.GetRecvCh() <- p

	select {
	case <-UpCh:
		t.Error("Premature up")
	default:
	}
	//
	// Get Response
	//
	response := <-as.GetSendCh()
	if response.Action() != packet.AuthResponse {
		t.Error("Expected Response Packet")
	}

	if string(response.Text()) != remotemd5sum {
		t.Error("Expected matching md5sums")
	}

	select {
	case <-UpCh:
		t.Error("Premature up")
	default:
	}
	//
	// Send Accept
	//
	p, err = packet.NewAuth(packet.AuthAccept)
	if err != nil {
		t.Errorf("NewAuth() Err:%s", err)
	}
	as.GetRecvCh() <- p

	select {
	case <-UpCh:
		t.Error("Premature up")
	default:
	}
	//
	// Get Challange
	//
	challenge := <-as.GetSendCh()
	if challenge.Action() != packet.AuthChallenge {
		t.Error("Expected Challenge Packet")
	}
	md5sum := generateMD5Sum(secret, string(challenge.Text()))

	select {
	case <-UpCh:
		t.Error("Premature up")
	default:
	}
	//
	// Send Response
	//
	p, err = packet.NewAuth(packet.AuthResponse, md5sum)
	if err != nil {
		t.Errorf("NewAuth() Err:%s", err)
	}
	as.GetRecvCh() <- p

	select {
	case <-UpCh:
		t.Error("Premature up")
	default:
	}

	//
	// Get Accept
	//
	reply := <-as.GetSendCh()
	if reply.Action() != packet.AuthAccept {
		t.Error("Expected Accept Packet")
	}

	select {
	case <-UpCh:
	case <-time.After(10*time.Millisecond):
		t.Error("Up TimeOut")
	}

	cx.Cancel()
}

// ----------------------------------------------------------------------------
// Local functions
//

func generateRandomSecret() string {

	size := 8

	random := make([]byte, size)
	_, err := rand.Read(random[:])
	if err != nil {
		log.Fatalf("ranf.Read() Random err:%s", err)
	}

	return base64.StdEncoding.EncodeToString(random)
}
