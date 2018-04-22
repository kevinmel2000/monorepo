package print

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/lab46/monorepo/gopkg/errors"
)

// this package is used to produce pretty print output
// mostly used by cli app
// all print should go to stdout instead of stderr

// var for printing prefix, this is naive and need a better implementation
var (
	prefixPrint = func(prePrefix, prefix string) string {
		return color.HiMagentaString(fmt.Sprintf("%s", prePrefix) + prefix)
	}
	prefixDebug = func(prePrefix, prefix string) string {
		return color.YellowString(fmt.Sprintf("%s", prePrefix) + prefix)
	}
	prefixInfo = func(prePrefix, prefix string) string {
		return color.GreenString(fmt.Sprintf("%s", prePrefix) + prefix)
	}
	prefixWarn = func(prePrefix, prefix string) string {
		return color.HiCyanString(fmt.Sprintf("%s", prePrefix) + prefix)
	}
	prefixError = func(prePrefix, prefix string) string {
		return color.RedString(fmt.Sprintf("%s", prePrefix) + prefix)
	}
	// debug var to identify if debug print is allowed or not
	isDebug bool
)

// list of pre-prefix
const (
	InfoPrePrefix  = "[INFO]"
	DebugPrePrefix = "[DEBUG]"
	WarnPrePrepix  = "[WARN]"
	ErrorPrePrefix = "[ERROR]"
)

func SetDebug(debug bool) {
	isDebug = debug
}

func Debug(v ...interface{}) {
	// idiomatic debug
	if !isDebug {
		return
	}
	print(prefixDebug(DebugPrePrefix, ""), v...)
}

func Info(v ...interface{}) {
	print(prefixInfo(InfoPrePrefix, ""), v...)
}

func Warn(v ...interface{}) {
	print(prefixWarn(WarnPrePrepix, ""), v...)
}

func Error(v ...interface{}) {
	print(prefixError(ErrorPrePrefix, ""), v...)
}

func Fatal(err error) {
	if err == nil {
		return
	}
	Error(err)
	os.Exit(2)
}

func print(prefix string, v ...interface{}) {
	// naively reject if only tag
	if len(v) == 0 {
		return
	}
	// return if parased argument is not valid
	parsedIntf := parseArgs(v...)
	if len(parsedIntf) == 0 {
		return
	}
	newIntf := []interface{}{prefix}
	newIntf = append(newIntf, parsedIntf...)
	fmt.Println(newIntf...)
}

// TODO: count on interface{} length and discard append to reduce memory use
func parseArgs(v ...interface{}) []interface{} {
	var newIntf []interface{}
	for key, val := range v {
		switch val.(type) {
		// dispatch if array of string
		case []string:
			arrOfString := val.([]string)
			for _, stringval := range arrOfString {
				newIntf = append(newIntf, stringval)
			}
			continue
		case nil:
			continue
		// pretty print errors if available
		case *errors.Errs:
			err := val.(*errors.Errs)
			newIntf = append(newIntf, err.Error())
			fields := err.GetFields()
			for key, val := range fields {
				s := fmt.Sprintf("%v=%v", key, val)
				newIntf = append(newIntf, s)
			}
			file, line := err.GetFileAndLine()
			if line != 0 {
				newIntf = append(newIntf, fmt.Sprintf("err_file=%s", file))
				newIntf = append(newIntf, fmt.Sprintf("err_line=%d", line))
			}
			continue
		}
		newIntf = append(newIntf, v[key])
	}
	return newIntf
}
