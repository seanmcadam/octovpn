package log

import (
	systemlog "log"
)

func init() {
	systemlog.Print("DEBUG Logging enabled")
}

//
// func Debug(v ...any)
// Debug calls Output to print to the standard logger. Arguments are handled in the manner of fmt.Print.
//
func Debug(v ...any) {
	var x []interface{}
	x = append(x,"DBG:")
	x = append(x, FileLine(2)+":")
	x = append(x, v)
	systemlog.Print(x...)
}

//
// func Debugf(format string, v ...any)
// Debugf calls Output to print to the standard logger. Arguments are handled in the manner of fmt.Printf.
//
func Debugf(format string, v ...any) {
	var x []any
	format = "DBG:%s:" + format
	x = append(x, FileLine(2))
	x = append(x, v...)
	systemlog.Printf(format, x...)
}

//
// func Debugln(v ...any)
// Debugln calls Output to print to the standard logger. Arguments are handled in the manner of fmt.Println.
//
func Debugln(v ...any) {
	var x []interface{}
	x = append(x,"DBG:")
	x = append(x, FileLine(2)+":")
	x = append(x, v)
	systemlog.Println(v...)
}
