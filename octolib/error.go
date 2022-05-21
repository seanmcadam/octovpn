package octolib

import (
	"fmt"

	"github.com/seanmcadam/octovpn/ctx"
)

type any interface{}

func ErrorLocationf(f string, v ...any) (e error) {
	flf := ctx.FileLineFunc(3)
	e = fmt.Errorf(flf+f, v)
	return e
}
