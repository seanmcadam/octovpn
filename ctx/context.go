package ctx

import "context"

type ContextValue string
type LogLevel string

//type any interface{}

const ContextKey ContextValue = "ContextKey"

const (
	LogLevelPanic LogLevel = "PANIC"
	LogLevelFatal LogLevel = "FATAL"
	LogLevelError LogLevel = "ERROR"
	LogLevelWarn  LogLevel = " WARN"
	LogLevelInfo  LogLevel = " INFO"
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelTrace LogLevel = "TRACE"
)

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

func (c *Ctx) Done() <-chan struct{} {
	return c.mycontext.Done()
}
