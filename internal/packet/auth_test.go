package packet

import (
	"testing"

	"github.com/seanmcadam/octovpn/octolib/log"
)

func TestNewAuth_compile(t *testing.T) {
	_, _ = NewAuth(AuthAccept)
}

func TestNewAuth_min_raw(t *testing.T) {
	var b []byte
	b = append(b, uint8(0xF0))
	_, err := MakeAuth(b)
	if err == nil {
		t.Errorf("MakeAuth() did not return error with only size")
	}
	b = append(b, uint8(AuthAccept))
	_, err = MakeAuth(b)
	if err == nil {
		t.Errorf("MakeAuth() did not return error")
	}

}

func TestMakeAuth_packets(t *testing.T) {

	tests := []struct {
		testno       int
		action       AuthPacketActionType
		text         string
		err_newauth  bool
		err_makeauth bool
	}{
		{1, AuthChallenge, "1234", false, false},
		{2, AuthResponse, "5678", false, false},
		{3, AuthAccept, "", false, false},
		{4, AuthReject, "", false, false},
		{5, AuthError, "this is an error", false, false},
		{6, AuthChallenge, "", true, false},
		{7, AuthResponse, "", true, false},
		{7, AuthError, "", true, false},
		{8, AuthAccept, "1234", false, true},
	}

	for _, test := range tests {
		log.Debugf("Run Test:%d", test.testno)

		AS, err := NewAuth(test.action, test.text)
		if err != nil {
			if !test.err_newauth {
				t.Errorf("NewAuth() did return error test:%d, err:%s", test.testno, err)
			}
		} else {
			if test.err_newauth {
				t.Errorf("NewAuth() did not return error test:%d", test.testno)
			} else {

				raw := AS.ToByte()
				if len(raw) != int(AS.Size()) {
					t.Errorf("ToByte bad len:%d", len(raw))
				}

				NewAS, err := MakeAuth(raw)
				if err != nil {
					if !test.err_newauth {
						t.Errorf("MakeAuth() did return error test:%d, err:%s", test.testno, err)
					}
				} else {
					if test.err_newauth {
						t.Errorf("NewAuth() did not return error test:%d", test.testno)
					} else {

						if NewAS.Action() != test.action {
							t.Errorf("NewAuth() actions dont match test:%d", test.testno)
						} else if NewAS.Size() != PacketSizeType(len(raw)) {
							t.Errorf("NewAuth() lengths dont match test:%d", test.testno)
						} else if string(NewAS.Text()) != string((AS.Text())) {
							t.Errorf("NewAuth() texts dont match test:%d", test.testno)
						}
					}
				}
			}
		}
	}
}
