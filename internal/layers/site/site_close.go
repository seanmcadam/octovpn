package site

import "github.com/seanmcadam/octovpn/octolib/log"

func (c *SiteStruct) doneChan() <-chan struct{} {
	if c == nil {
		return nil
	}
	return c.cx.DoneChan()
}

func (c *SiteStruct) Cancel() {
	if c == nil {
		return
	}
	log.GDebug("Cancel() %s", c.name)
	c.link.Down()
	c.link.Close()
	c.cx.Cancel()
}

func (t *SiteStruct) closed() bool {
	if t == nil {
		return true
	}
	return t.cx.Done()
}
