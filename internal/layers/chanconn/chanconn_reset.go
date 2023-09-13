package chanconn

import (
	"github.com/seanmcadam/octovpn/octolib/errors"
	"github.com/seanmcadam/octovpn/octolib/log"
)

func (cs *ChanconnStruct) Reset() error {
	if cs == nil {
		return errors.ErrNetNilMethodPointer(log.Errf(""))
	}

	return cs.conn.Reset()
}
