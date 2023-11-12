//go:build darwin
// +build darwin

package darwin

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa
#import "AppDelegate.h"

*/
import "C"

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"github.com/wailsapp/wails/v2/pkg/options"
)

func SetupSingleInstance(uniqueID string) {
	lockFilePath := getTempDir()
	lockFileName := uniqueID + ".lock"
	_, err := createLockFile(lockFilePath + "/" + lockFileName)
	// if lockFile exist – send notification to second instance
	if err != nil {
		c := NewCalloc()
		defer c.Free()
		singleInstanceUniqueId := c.String(uniqueID)

		data, err := options.NewSecondInstanceData()
		if err != nil {
			return
		}

		serialized, err := json.Marshal(data)
		if err != nil {
			return
		}

		C.SendDataToFirstInstance(singleInstanceUniqueId, c.String(string(serialized)))

		os.Exit(0)
	}
}

//export HandleSecondInstanceData
func HandleSecondInstanceData(secondInstanceMessage *C.char) {
	message := C.GoString(secondInstanceMessage)

	var secondInstanceData options.SecondInstanceData

	err := json.Unmarshal([]byte(message), &secondInstanceData)
	if err == nil {
		secondInstanceBuffer <- secondInstanceData
	}
}

// CreateLockFile tries to create a file with given name and acquire an
// exclusive lock on it. If the file already exists AND is still locked, it will
// fail.
func createLockFile(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		fmt.Printf("Failed to open lockfile %s: %s", filename, err)
		return nil, err
	}

	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		// Flock failed for some other reason than other instance already lock it. Print it in logs for possible debugging.
		if !strings.Contains(err.Error(), "resource temporarily unavailable") {
			fmt.Printf("Failed to lock lockfile %s: %s", filename, err)
		}
		file.Close()
		return nil, err
	}

	return file, nil
}

// If app is sandboxed, golang os.TempDir() will return path that will not be accessible. So use native macOS temp dir function.
func getTempDir() string {
	cstring := C.GetMacOsNativeTempDir()
	path := C.GoString(cstring)
	C.free(unsafe.Pointer(cstring))

	return path
}
