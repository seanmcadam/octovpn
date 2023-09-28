package instance

import (
	"testing"

	"github.com/seanmcadam/octovpn/octolib/log"
)

func TestNew(t *testing.T) {

	i := New()
	log.Infof("Name:%s", i.Next())
	log.Infof("Name:%s", i.Next())
}
