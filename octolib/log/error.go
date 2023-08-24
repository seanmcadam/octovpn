package log

import (
	systemlog "log"
)

//
// func Debug(v ...any)
// Debug calls Output to print to the standard logger. Arguments are handled in the manner of fmt.Print.
//
func Error(v ...any) {
	var x []interface{}
	x = append(x,"ERR:")
	x = append(x, FileLine(2)+":")
	x = append(x, v)
	systemlog.Print(x...)
}

//
// func Debugf(format string, v ...any)
// Debugf calls Output to print to the standard logger. Arguments are handled in the manner of fmt.Printf.
//
func Errorf(format string, v ...any) {
	var x []any
	format = "ERR:%s:" + format
	x = append(x, FileLine(2))
	x = append(x, v...)
	systemlog.Printf(format, x...)
}

//
// func Debugln(v ...any)
// Debugln calls Output to print to the standard logger. Arguments are handled in the manner of fmt.Println.
//
func Errorln(v ...any) {
	var x []interface{}
	x = append(x,"ERR:")
	x = append(x, FileLine(2)+":")
	x = append(x, v)
	systemlog.Println(v...)
}
