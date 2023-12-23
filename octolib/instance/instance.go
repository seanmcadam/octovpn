package instance

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/seanmcadam/counter"
	"github.com/seanmcadam/ctx"
	log "github.com/seanmcadam/loggy"
)

type InstanceName struct {
	name string
}

type Instance struct {
	module  string
	counter counter.Counter
}

func New() (i *Instance) {
	var module = ""

	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		log.FatalStack("runtime caller failed")
	}

	funcInfo := runtime.FuncForPC(pc)
	if funcInfo == nil {
		log.FatalStack("unknown module name")
	}

	fullName := funcInfo.Name()
	parts := strings.Split(fullName, "/")
	if len(parts) > 1 {
		module = parts[len(parts)-1]
	}

	i = &Instance{
		module:  module,
		counter: counter.NewCounter32(ctx.New()),
	}

	return i

}

func NewInstanceName(s string) *InstanceName {
	return &InstanceName{name: s}
}

// -
//
// -
func (i *Instance) Next() (in *InstanceName) {
	name := i.module
	count := i.counter.Next().Uint()
	in = &InstanceName{
		name: fmt.Sprintf("%s-%d", name, count),
	}
	return in
}

// -
//
// -
func (i *InstanceName) String() string {
	return i.name
}
