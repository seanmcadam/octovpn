package log


import (
	"io"
	systemlog "log"
)

//
// Mirror the log package for Logger
//
type Logger systemlog.Logger

//
//
// func Fatal(v ...any)
// Fatal is equivalent to Print() followed by a call to os.Exit(1).
//
func Fatal(v ...any) {
	systemlog.Fatal(v...)
}

//
// func Fatalf(format string, v ...any)
// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
//
func Fatalf(format string, v ...any) {
	systemlog.Fatalf(format, v...)
}

//
// func Fatalln(v ...any)
// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
//
func Fatalln(v ...any) {
	systemlog.Fatalln(v...)
}

//
// func Panic(v ...any)
// Panic is equivalent to Print() followed by a call to panic().
//
func Panic(v ...any) {
	systemlog.Panic(v...)
}

//
// func Panicf(format string, v ...any)
// Panicf is equivalent to Printf() followed by a call to panic().
//
func Panicf(format string, v ...any) {
	systemlog.Panicf(format, v...)
}

//
// func systemlog.Panicln(v ...any)
// Panicln is equivalent to Println() followed by a call to panic().
//
func Panicln(v ...any) {
	systemlog.Panicln(v...)
}

//
// func Print(v ...any)
// Print calls Output to print to the standard logger. Arguments are handled in the manner of fmt.Print.
//
func Print(v ...any) {
	systemlog.Print(v...)
}

//
// func Printf(format string, v ...any)
// Printf calls Output to print to the standard logger. Arguments are handled in the manner of fmt.Printf.
//
func Printf(format string, v ...any) {
	systemlog.Printf(format, v...)
}

//
// func Println(v ...any)
// Println calls Output to print to the standard logger. Arguments are handled in the manner of fmt.Println.
//
func Println(v ...any) {
	systemlog.Println(v...)
}

//
// func Flags() int
// Flags returns the output flags for the standard logger. The flag bits are Ldate, Ltime, and so on.
//
func Flags() int {
	return systemlog.Flags()
}

//
// func Output(calldepth int, s string) error
// Output writes the output for a logging event. The string s contains the text to print after the prefix
// specified by the flags of the Logger. A newline is appended if the last character of s is not already a newline.
// Calldepth is the count of the number of frames to skip when computing the file name and line number if Llongfile
// or Lshortfile is set; a value of 1 will print the details for the caller of Output.
//
func Output(calldepth int, s string) error {
	return systemlog.Output(calldepth, s)
}

//
// func Prefix() string
// Prefix returns the output prefix for the standard logger.
//
func Prefix() string {
	return systemlog.Prefix()
}

//
// func Writer() io.Writer
// Writer returns the output destination for the standard logger.
//
func Writer() io.Writer {
	return systemlog.Writer()
}

//
// func systemlog.Default() *systemlog.Logger
// Default returns the standard logger used by the package-level output functions.
//
func Default() *systemlog.Logger {
	return systemlog.Default()
}

//
// func systemlog.New(out io.Writer, prefix string, flag int) *systemlog.Logger
// New creates a new Logger. The out variable sets the destination to which log data will be written.
// The prefix appears at the beginning of each generated log line, or after the log header if the
// Lmsgprefix flag is provided. The flag argument defines the logging properties.
//
func New(out io.Writer, prefix string, flag int) *systemlog.Logger {
	return systemlog.New(out, prefix, flag)
}
