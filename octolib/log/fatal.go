package log

import (
	systemlog "log"
)

//
// func Debug(v ...any)
// Debug calls Output to print to the standard logger. Arguments are handled in the manner of fmt.Print.
//
func FFatal(v ...any) {
	var x []interface{}
	x = append(x,"ERR:")
	x = append(x, FileLine(2)+":")
	x = append(x, v)
	logLock.Lock()
	defer logLock.Unlock()
	systemlog.Fatal(x...)
}

//
// func Debugf(format string, v ...any)
// Debugf calls Output to print to the standard logger. Arguments are handled in the manner of fmt.Printf.
//
func FFatalf(format string, v ...any) {
	var x []any
	format = "ERR:%s:" + format
	x = append(x, FileLine(2))
	x = append(x, v...)
	logLock.Lock()
	defer logLock.Unlock()
	systemlog.Fatalf(format, x...)
}

//
// func Debugln(v ...any)
// Debugln calls Output to print to the standard logger. Arguments are handled in the manner of fmt.Println.
//
func FFatalln(v ...any) {
	var x []interface{}
	x = append(x,"ERR:")
	x = append(x, FileLine(2)+":")
	x = append(x, v)
	logLock.Lock()
	defer logLock.Unlock()
	systemlog.Fatalln(v...)
}
