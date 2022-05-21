package ctx

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

func FileLineFunc(depth int) string {
	return fileLineFunc(depth)
}

//
// [0060|file.go:153:function()]:
//  GID  file    Line# func()
//
func fileLineFunc(a ...int) string {
	var depth int = 1

	if len(a) != 0 {
		depth = a[0]
		if depth < 1 && depth > 10 {
			panic(FileLine() + " depth out of bounds")
		}
	}

	pc, file, line, ok := runtime.Caller(depth)
	if !ok {
		return "FileLine() failed"
	}

	lineno := strconv.Itoa(line)

	fileArr := strings.Split(file, "/")
	fileName := fileArr[len(fileArr)-1]

	funcPtr := runtime.FuncForPC(pc)

	var funcName string

	if funcPtr == nil {
		funcName = "TheUNKNOWNFunction()"
	} else {
		f := strings.Split(funcPtr.Name(), "/")
		funcName = f[len(f)-1]
		g := strings.Split(funcName, ".")
		funcName = g[len(g)-1]
	}

	gid := fmt.Sprintf("%04d", getGID())

	ret := "[" + gid + "|" + fileName + ":" + lineno + ":" + funcName + "()]:"
	return ret
}

//
//
//
func FileLine() string {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return "FileLine() failed"
	}

	f := strings.Split(file, "/")

	lineno := strconv.Itoa(line)

	return "[" + f[len(f)-1] + ":" + lineno + "]:"
}

//
//
//
func Funcname() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "TheUnknownFunction()"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "TheUnknownFunction()"
	}

	f := strings.Split(fn.Name(), "/")

	return f[len(f)-1]
}

//
//
//
func Errtrace() string {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return "file?[0]:func?"
	}

	lineno := strconv.Itoa(line)

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return file + "[" + lineno + "]:func?"
	}

	return file + "[" + lineno + "]:" + fn.Name()
}

//
//
//
func PanicHere(text ...string) string {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("Well, this is unexpected...")
	}

	f := strings.Split(file, "/")

	lineno := strconv.Itoa(line)

	panic(fmt.Sprintf("[%s:%s]:%s", f[len(f)-1], lineno, text[0]))
}

//
//
// borrowed from https://blog.sgmansfield.com/2015/12/goroutine-ids/
//
func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

//
// goCounter()
// Generates a UniqueID (int) and returns via supplied channel
//
func RunGoCounter(init int, ch chan int) {
	go func() {

		counter := init
		for {
			ch <- counter
			counter += 1
		}
	}()
}
