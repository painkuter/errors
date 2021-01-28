package errors

import (
	"fmt"
	"regexp"
	"runtime"
)

const (
	maxCallers  = 20
	skipCallers = 3
)

var scopeRe = regexp.MustCompile(`\.\([^)]+\)`)

func getStackTrace() []string {
	result := make([]string, 0, maxCallers)
	pc := make([]uintptr, maxCallers)
	n := runtime.Callers(skipCallers, pc)
	if n == 0 {
		return result
	}

	pc = pc[:n]
	frames := runtime.CallersFrames(pc)
	frames.Next() // skip errors.ErrWithStack call

	more := true
	var frame runtime.Frame
	for i := 0; more && i < maxCallers; i++ {
		frame, more = frames.Next()
		funcName := scopeRe.ReplaceAllString(frame.Function, "")
		result = append(result, fmt.Sprintf("%s at %s:%d", funcName, frame.File, frame.Line))
	}

	return result
}
