// +build windows,!debug

package ffenestri

import "C"
import (
	"fmt"
	"golang.org/x/sys/windows"
	"log"
	"os"
	"runtime"
	"strings"
	"syscall"
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

	title, err := syscall.UTF16PtrFromString("Fatal Error")
	if err != nil {
		log.Fatal(err)
	}

	text, err := syscall.UTF16PtrFromString("There has been a fatal error. Details:\n" + message)
	if err != nil {
		log.Fatal(err)
	}

	var flags uint32 = windows.MB_ICONERROR | windows.MB_OK

	_, err = windows.MessageBox(0, text, title, flags|windows.MB_SYSTEMMODAL)
	os.Exit(1)
	return err
}
