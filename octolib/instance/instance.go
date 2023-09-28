package instance

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/seanmcadam/octovpn/internal/counter"
	"github.com/seanmcadam/octovpn/octolib/ctx"
	"github.com/seanmcadam/octovpn/octolib/log"
)

type Instance struct {
	module  string
	counter counter.CounterStruct
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
	parts := strings.Split(fullName, ".")
	if len(parts) > 1 {
		module = parts[0]
	}

	i = &Instance{
		module:  module,
		counter: counter.NewCounter32(ctx.NewContext()),
	}

	return i

}

// -
//
// -
func (i *Instance) Next() string {
	return fmt.Sprintf("%s-%d", i.module, i.counter.Next().Uint().(uint32))
}
