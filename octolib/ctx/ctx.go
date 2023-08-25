package ctx

import "context"

type ContextValue string

type Ctx struct {
	mycontext context.Context
	cancel    func()
}

func NewContext() (c *Ctx) {
	ctx, cancel := context.WithCancel(context.Background())
	c = &Ctx{
		mycontext: ctx,
		cancel:    cancel,
	}
	return c
}

func (c *Ctx) NewWithCancel() (d *Ctx) {
	ctx, cancel := context.WithCancel(c.mycontext)
	d = &Ctx{
		mycontext: ctx,
		cancel:    cancel,
	}
	return d
}

func (c *Ctx) Context() (ctx context.Context) {
	return c.mycontext
}

func (c *Ctx) Cancel() {
	c.cancel()
}

func (c *Ctx) DoneChan() <-chan struct{} {
	return c.mycontext.Done()
}

func (c *Ctx) Done() bool {
	select {
	case <-c.DoneChan():
		return true
	default:
		return false
	}
}
