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
	"github.com/wailsapp/wails/v2/pkg/options"
	"os"
	"syscall"
)

func SetupSingleInstance(uniqueID string) {
	lockFilePath := os.TempDir()
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
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		file.Close()
		return nil, err
	}

	return file, nil
}
