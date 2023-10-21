package log

import (
	systemlog "log"
)

// func Debug(v ...any)
// Debug calls Output to print to the standard logger. Arguments are handled in the manner of fmt.Print.
func stack() {
	var start uint8 = 4

	systemlog.Println("STACK TRACE:")
	stack := getStack(start)
	for _, s := range stack {
		systemlog.Printf("\t%s", s)
	}
}
