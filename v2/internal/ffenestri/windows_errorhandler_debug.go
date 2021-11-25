//go:build windows && debug
// +build windows,debug

package ffenestri

import (
	"fmt"
	"github.com/ztrue/tracerr"
	"runtime"
	"strings"
)

func wall(err error, inputs ...interface{}) error {
	if err == nil {
		return nil
	}
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	splitName := strings.Split(funcName, ".")
	message := "[" + splitName[len(splitName)-1] + "]"
	if len(inputs) > 0 {
		params := []string{}
		for _, param := range inputs {
			params = append(params, fmt.Sprintf("%v", param))
		}
		message += "(" + strings.Join(params, " ") + ")"
	}
	return tracerr.Errorf(message)
}
