//go:build darwin

package application

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa

#include <stdlib.h>
#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>

static void SendDataToFirstInstance(char *singleInstanceUniqueId, char* message) {
	NSLog(@"SendDataToFirstInstance: singleInstanceUniqueId='%s' message='%@'\n", singleInstanceUniqueId, [NSString stringWithUTF8String:message]);
    [[NSDistributedNotificationCenter defaultCenter]
        postNotificationName:[NSString stringWithUTF8String:singleInstanceUniqueId]
        object:nil
        userInfo:@{@"message": [NSString stringWithUTF8String:message]}
        deliverImmediately:YES];
}

*/
import "C"
import (
	"encoding/json"
	"os"
	"syscall"
	"unsafe"
)

type unixLock struct {
	file     *os.File
	uniqueID string
	manager  *singleInstanceManager
}

func newPlatformLock(manager *singleInstanceManager) (platformLock, error) {
	return &unixLock{
		manager: manager,
	}, nil
}

func (l *unixLock) acquire(uniqueID string) error {
	l.uniqueID = uniqueID
	lockFilePath := os.TempDir()
	lockFileName := uniqueID + ".lock"
	var err error
	l.file, err = createLockFile(lockFilePath + "/" + lockFileName)
	if err != nil {
		return alreadyRunningError
	}
	return nil
}

func (l *unixLock) release() {
	if l.file != nil {
		syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
		l.file.Close()
		os.Remove(l.file.Name())
		l.file = nil
	}
}

func (l *unixLock) notify(data string) error {
	singleInstanceUniqueId := C.CString(l.uniqueID)
	defer C.free(unsafe.Pointer(singleInstanceUniqueId))
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))

	C.SendDataToFirstInstance(singleInstanceUniqueId, cData)

	os.Exit(l.manager.options.ExitCode)
	return nil
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

//export handleSecondInstanceData
func handleSecondInstanceData(secondInstanceMessage *C.char) {
	message := C.GoString(secondInstanceMessage)

	var secondInstanceData SecondInstanceData

	err := json.Unmarshal([]byte(message), &secondInstanceData)
	if err == nil {
		secondInstanceBuffer <- secondInstanceData
	}
}
