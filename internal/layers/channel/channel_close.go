package channel

import "github.com/seanmcadam/octovpn/octolib/log"

func (c *ChannelStruct) doneChan() <-chan struct{} {
	if c == nil {
		return nil
	}
	return c.cx.DoneChan()
}

func (c *ChannelStruct) Cancel() {
	if c == nil {
		return
	}
	log.GDebug("Cancel() %s", c.name)
	c.link.Down()
	c.link.Close()
	c.cx.Cancel()
}

func (t *ChannelStruct) closed() bool {
	if t == nil {
		return true
	}
	return t.cx.Done()
}
