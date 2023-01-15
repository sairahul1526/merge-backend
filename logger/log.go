package logger

import (
	"fmt"
	CONFIG "merge-backend/config"
	"path/filepath"
	"runtime"
	"strings"
)

// Log - logs based on config
func Log(str ...interface{}) {
	if CONFIG.Log {
		fmt.Println(`{"type": "log", "location": "`+functionLoc(2)+`", "function": "`+functionName(2)+`", "log": "`, str, `"}`)
	}
}

// Log - logs all
func Warn(str ...interface{}) {
	fmt.Println(`{"type": "warn", "location": "`+functionLoc(2)+`", "function": "`+functionName(2)+`", "log": "`, str, `"}`)
}

// functionLoc - returns last two path tokens of caller
func functionLoc(optFuncLevel ...int) string {
	frameLevel := 1 // default to the caller's frame
	if len(optFuncLevel) > 0 {
		frameLevel = optFuncLevel[0]
	}

	_, fPath, line, ok := runtime.Caller(frameLevel)
	if !ok {
		return ""
	}

	return fmt.Sprintf("%s:%d", last2Tokens(fPath, "/"), line)
}

// functionName - returns the function name of the caller
func functionName(optFuncLevel ...int) (funcName string) {
	frameLevel := 1 // default to the caller's frame
	if len(optFuncLevel) > 0 {
		frameLevel = optFuncLevel[0]
	}

	if pc, _, _, ok := runtime.Caller(frameLevel); ok {
		fPtr := runtime.FuncForPC(pc)
		if fPtr == nil {
			return
		}

		return last2Tokens(fPtr.Name(), "/")
	}
	return
}

func last2Tokens(str, separator string) (lastTokens string) {
	tokens := strings.Split(str, "/")
	if len(tokens) >= 2 {
		lastTokens = strings.Join(tokens[len(tokens)-2:], separator)
	} else {
		lastTokens = filepath.Base(str)
	}
	return
}
