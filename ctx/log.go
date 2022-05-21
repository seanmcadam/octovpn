package ctx

import "log"

func (*Ctx) LogPanic(v ...any)            { log.Panic(string(LogLevelPanic)+":", v) }
func (*Ctx) LogPanicf(f string, v ...any) { log.Panicf(string(LogLevelPanic)+":"+f, v) }
func (*Ctx) LogPanicln(v ...any)          { log.Panicln(string(LogLevelPanic)+":", v) }
func (*Ctx) LogFatal(v ...any)            { log.Fatal(string(LogLevelFatal)+":", v) }
func (*Ctx) LogFatalf(f string, v ...any) { log.Fatalf(string(LogLevelFatal)+":"+f, v) }
func (*Ctx) LogFatalln(v ...any)          { log.Fatalln(string(LogLevelFatal)+":", v) }
func (*Ctx) LogError(v ...any)            { log.Print(string(LogLevelError)+":", v) }
func (*Ctx) LogErrorf(f string, v ...any) { log.Printf(string(LogLevelError)+":"+f, v) }
func (*Ctx) LogErrorln(v ...any)          { log.Println(string(LogLevelError)+":", v) }
func (*Ctx) LogWarn(v ...any)             { log.Print(string(LogLevelWarn)+":", v) }
func (*Ctx) LogWarnf(f string, v ...any)  { log.Printf(string(LogLevelWarn)+":"+f, v) }
func (*Ctx) LogWarnln(v ...any)           { log.Println(string(LogLevelWarn)+":", v) }
func (*Ctx) LogInfo(v ...any)             { log.Print(string(LogLevelInfo)+":", v) }
func (*Ctx) LogInfof(f string, v ...any)  { log.Printf(string(LogLevelInfo)+":"+f, v) }
func (*Ctx) LogInfoln(v ...any)           { log.Println(string(LogLevelInfo)+":", v) }
func (*Ctx) LogDebug(v ...any)            { log.Print(string(LogLevelDebug)+":", v) }
func (*Ctx) LogDebugf(f string, v ...any) { log.Printf(string(LogLevelDebug)+":"+f, v) }
func (*Ctx) LogDebugln(v ...any)          { log.Println(string(LogLevelDebug)+":", v) }
func (*Ctx) LogTrace(v ...any)            { log.Print(string(LogLevelTrace)+":", v) }
func (*Ctx) LogTracef(f string, v ...any) { log.Printf(string(LogLevelTrace)+":"+f, v) }
func (*Ctx) LogTraceln(v ...any)          { log.Println(string(LogLevelTrace)+":", v) }

func (c *Ctx) Log(l LogLevel, v ...any) {
	flf := fileLineFunc(2)
	switch l {
	case LogLevelPanic:
		c.LogPanic(flf, v)
	case LogLevelFatal:
		c.LogFatal(flf, v)
	case LogLevelError:
		c.LogError(flf, v)
	case LogLevelWarn:
		c.LogWarn(flf, v)
	case LogLevelInfo:
		c.LogInfo(flf, v)
	case LogLevelDebug:
		c.LogDebug(flf, v)
	case LogLevelTrace:
		c.LogTrace(flf, v)
	default:
		c.LogPanicf("Reached default: %s", l)
	}
}
func (c *Ctx) Logf(l LogLevel, f string, v ...any) {
	flf := fileLineFunc(2)
	switch l {
	case LogLevelPanic:
		if len(v) > 0 {
			c.LogPanicf(flf+f, v...)
		} else {
			c.LogPanic(flf + f)
		}
	case LogLevelFatal:
		if len(v) > 0 {
			c.LogFatalf(flf+f, v...)
		} else {
			c.LogFatal(flf + f)
		}
	case LogLevelError:
		if len(v) > 0 {
			c.LogErrorf(flf+f, v...)
		} else {
			c.LogError(flf + f)
		}
	case LogLevelWarn:
		if len(v) > 0 {
			c.LogWarnf(flf+f, v...)
		} else {
			c.LogWarn(flf + f)
		}
	case LogLevelInfo:
		if len(v) > 0 {
			c.LogInfof(flf+f, v...)
		} else {
			c.LogInfo(flf + f)
		}
	case LogLevelDebug:
		if len(v) > 0 {
			c.LogDebugf(flf+f, v)
		} else {
			c.LogDebug(flf + f)
		}
	case LogLevelTrace:
		if len(v) > 0 {
			c.LogTracef(flf+f, v)
		} else {
			c.LogTrace(flf + f)
		}
	default:
		c.LogPanicf("Reached default: %s", l)
	}
}
func (c *Ctx) Logln(l LogLevel, v ...any) {
	flf := fileLineFunc(2)
	switch l {
	case LogLevelPanic:
		c.LogPanicln(flf, v)
	case LogLevelFatal:
		c.LogFatalln(flf, v)
	case LogLevelError:
		c.LogErrorln(flf, v)
	case LogLevelWarn:
		c.LogWarnln(flf, v)
	case LogLevelInfo:
		c.LogInfoln(flf, v)
	case LogLevelDebug:
		c.LogDebugln(flf, v)
	case LogLevelTrace:
		c.LogTraceln(flf, v)
	default:
		c.LogPanicf("Reached default: %s", l)
	}
}
