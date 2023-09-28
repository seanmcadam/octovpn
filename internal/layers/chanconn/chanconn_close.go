package chanconn

import "github.com/seanmcadam/octovpn/octolib/log"

func (c *ChanconnStruct) doneChan() <-chan struct{} {
	if c == nil {
		return nil
	}
	return c.cx.DoneChan()
}

func (c *ChanconnStruct) Cancel() {
	if c == nil {
		return
	}
	log.GDebugf("Cancel() %s", c.name)
	c.link.Close()
	c.cx.Cancel()
}

func (t *ChanconnStruct) closed() bool {
	if t == nil {
		return true
	}
	return t.cx.Done()
}
