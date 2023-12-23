package instance

import (
	"testing"

	log "github.com/seanmcadam/loggy"
)

func TestNew(t *testing.T) {

	i := New()
	log.Infof("Name:'%s'", i.Next())
	log.Infof("Name:'%s'", i.Next())
}
