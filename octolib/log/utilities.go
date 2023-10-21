package log

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

//
// The parameter (a ...int) indicate the depth of the call stack that the file and line should taken from
//

// GidFileLine()
// Returns [GID|Filename:line#]:
// Useful for pointing out error message locations
func GidFileLine(a ...int) string {
	var depth int = 1

	if len(a) == 1 {
		depth = a[0]
	}

	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		return "GIDFileLine() failed"
	}

	lineno := strconv.Itoa(line)

	fileArr := strings.Split(file, "/")
	fileName := fileArr[len(fileArr)-1]

	gid := fmt.Sprintf("%04d", getGID())

	ret := "[" + gid + "|" + fileName + ":" + lineno + "]:"
	return ret
}

// FileLine()
// Returns [filename:Line#]
// Useful for pointing out error message locations
func FileLine(a ...int) string {
	var depth = 1

	// Get the first int as the depth
	if len(a) == 1 {
		depth = a[0]
	}

	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		return "FileLine() failed"
	}

	f := strings.Split(file, "/")

	lineno := strconv.Itoa(line)

	return "[" + f[len(f)-1] + ":" + lineno + "]:"
}

// FileLineFunc()
// Returns [Filename:line#:funcname()]:
// Useful for pointing out error message locations
func FileLineFunc(a ...int) string {
	var depth int = 1

	if len(a) == 1 {
		depth = a[0]
	}

	pc, file, line, ok := runtime.Caller(depth)
	if !ok {
		return "GidFileLine() failed"
	}

	lineno := strconv.Itoa(line)

	fileArr := strings.Split(file, "/")
	fileName := fileArr[len(fileArr)-1]

	funcPtr := runtime.FuncForPC(pc)

	var funcName string

	if funcPtr == nil {
		funcName = "UnknownFunction()"
	} else {
		f := strings.Split(funcPtr.Name(), "/")
		funcName = f[len(f)-1]
		g := strings.Split(funcName, ".")
		funcName = g[len(g)-1]
	}

	ret := "[" + fileName + ":" + lineno + ":" + funcName + "()]:"
	return ret
}

// GidFileLineFunc()
// Returns [GID|Filename:line#:funcname()]:
// Useful for pointing out error message locations
func GidFileLineFunc(a ...int) string {
	var depth int = 1

	if len(a) == 1 {
		depth = a[0]
	}

	pc, file, line, ok := runtime.Caller(depth)
	if !ok {
		return "GidFileLine() failed"
	}

	lineno := strconv.Itoa(line)

	fileArr := strings.Split(file, "/")
	fileName := fileArr[len(fileArr)-1]

	funcPtr := runtime.FuncForPC(pc)

	var funcName string

	if funcPtr == nil {
		funcName = "UnknownFunction()"
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

// func getGID()
// Returns the current threads GID
// Handy to see which GIDs are producing logs
func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// getStack()
func getStack(depth uint8) (ret []string) {

	pc := make([]uintptr, 20) // Adjust the size as needed
	n := runtime.Callers(0, pc)

	var j int = 0
	for i := depth; i < uint8(n); i++ {
		funcPtr := runtime.FuncForPC(pc[i])
		if funcPtr != nil {
			file, line := funcPtr.FileLine(pc[i])
			j++
			ret = append(ret, fmt.Sprintf("%d:[%s|%d]%s()\n", j, file, line, funcPtr.Name()))
		}
	}
	return ret
	//return reverseArray(ret)
}

