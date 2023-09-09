package log

import (
	"fmt"
	systemlog "log"
)

// func Err(v ...any)
func Errf(format string, v ...any) error {
	format = "[%s]:" + format
	var x []interface{}
	x = append(x, "Error:")
	x = append(x, FileLine(2)+":")
	x = append(x, v)
	y := append(x, v...)
	return fmt.Errorf(format, y...)
}

// func Debug(v ...any)
// Debug calls Output to print to the standard logger. Arguments are handled in the manner of fmt.Print.
func Error(v ...any) {
	var x []interface{}
	x = append(x, "ERR:")
	x = append(x, FileLine(2)+":")
	x = append(x, v)
	logLock.Lock()
	defer logLock.Unlock()
	systemlog.Print(x...)
}

// func Errorf(format string, v ...any)
// Debugf calls Output to print to the standard logger. Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, v ...any) {
	var x []any
	xformat := "ERR:%s:" + format
	x = append(x, FileLine(2))
	x = append(x, v...)
	logLock.Lock()
	defer logLock.Unlock()
	systemlog.Printf(xformat, x...)
}

// func Errorln(v ...any)
// Debugln calls Output to print to the standard logger. Arguments are handled in the manner of fmt.Println.
func Errorln(v ...any) {
	var x []interface{}
	x = append(x, "ERR:")
	x = append(x, FileLine(2)+":")
	x = append(x, v)
	logLock.Lock()
	defer logLock.Unlock()
	systemlog.Println(v...)
}

// func Debug(v ...any)
// Debug calls Output to print to the standard logger. Arguments are handled in the manner of fmt.Print.
func ErrorStack(v ...any) {
	var x []interface{}
	x = append(x, "ERR:")
	x = append(x, FileLine(2)+":")
	x = append(x, v)
	logLock.Lock()
	defer logLock.Unlock()
	stack()
	systemlog.Print(x...)
}

// func Errorf(format string, v ...any)
// Debugf calls Output to print to the standard logger. Arguments are handled in the manner of fmt.Printf.
func ErrorfStack(format string, v ...any) {
	var x []any
	xformat := "ERR:%s:" + format
	x = append(x, FileLine(2))
	x = append(x, v...)
	logLock.Lock()
	defer logLock.Unlock()
	stack()
	systemlog.Printf(xformat, x...)
}

// func Errorln(v ...any)
// Debugln calls Output to print to the standard logger. Arguments are handled in the manner of fmt.Println.
func ErrorlnStack(v ...any) {
	var x []interface{}
	x = append(x, "ERR:")
	x = append(x, FileLine(2)+":")
	x = append(x, v)
	logLock.Lock()
	defer logLock.Unlock()
	stack()
	systemlog.Println(v...)
}
