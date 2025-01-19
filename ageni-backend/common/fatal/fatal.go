package fatal

import (
	"fmt"

	"github.com/readonme/open-studio/common/log"

	"runtime"
	"runtime/debug"
	"strings"
)

// RecoverFromPanic recovers the panic for the calling function and logs panic data to logs and metrics
// Example: defer fatal.RecoverFromPanic("content-service", "panic")
// This DOES NOT work: defer func() {fatal.RecoverFromPanic("content-service", "panic")} ()
func RecoverFromPanic() {
	if r := recover(); r != nil {
		LogPanicRecover(r, string(debug.Stack()), 0)
	}
}

// ConvertPanicToErrorAndLog converts panic value to error and logs panic data to logs and metrics
// Example:
//
//	defer func() {
//	  if err:= fatal.ConvertPanicToErrorAndLog(recover()); err != nil {
//	    // error handling
//	  }
//	}()
func ConvertPanicToErrorAndLog(r interface{}) error {
	if r != nil {
		LogPanicRecover(r, string(debug.Stack()), 1)
		return fmt.Errorf("panic: %v", r)
	}
	return nil
}

func getPanicCodeLocation(skip int) string {
	var (
		pc   [16]uintptr
		file string
		line int
	)
	n := runtime.Callers(skip, pc[:])
	for _, pc := range pc[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		if !strings.HasPrefix(fn.Name(), "runtime.") {
			file, line = fn.FileLine(pc)
			break
		}
	}
	fileSplit := strings.Split(file, "/")
	shortFile := fileSplit[len(fileSplit)-1]
	return fmt.Sprintf("%s-%d", shortFile, line)
}

// LogPanicRecover prints unified error logs
// The argument skip is the number of stack frames to ascend, with 0 identifying the caller of LogPanicRecover.
func LogPanicRecover(r interface{}, stack string, skip int) {
	codeLocation := getPanicCodeLocation(skip + 4)
	log.Errorf("[Panic] codeLocation:%s panic:%v stack:%v", codeLocation, r, stack)
}
