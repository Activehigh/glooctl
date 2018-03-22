package log

import (
	"fmt"
	"os"
	"runtime"

	"regexp"

	"io"

	"github.com/k0kubun/pp"
)

var debugMode = os.Getenv("DEBUG") == "1"
var DefaultOut io.Writer = os.Stdout

var rxp = regexp.MustCompile(".*/src/")

func Sprintf(format string, a ...interface{}) string {
	return pp.Sprintf("%v\t"+format, append([]interface{}{line()}, a...)...)
}

func GreyPrintf(format string, a ...interface{}) {
	fmt.Fprintf(DefaultOut, "%v\t"+format+"\n", append([]interface{}{line()}, a...)...)
}

func Printf(format string, a ...interface{}) {
	pp.Fprintf(DefaultOut, "%v\t"+format+"\n", append([]interface{}{line()}, a...)...)
}

func Warnf(format string, a ...interface{}) {
	pp.Fprintf(DefaultOut, "WARNING: %v\t"+format+"\n", append([]interface{}{line()}, a...)...)
}

func Debugf(format string, a ...interface{}) {
	if debugMode {
		pp.Fprintf(DefaultOut, "%v\t"+format+"\n", append([]interface{}{line()}, a...)...)
	}
}

func Fatalf(format string, a ...interface{}) {
	pp.Fprintf(DefaultOut, "%v\t"+format+"\n", append([]interface{}{line()}, a...)...)
}

func line() string {
	_, file, line, _ := runtime.Caller(2)
	file = rxp.ReplaceAllString(file, "")
	return fmt.Sprintf("%v:%v", file, line)
}
