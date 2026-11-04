//go:build darwin && !ios && !server && !purego

package application

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa

#include <stdlib.h>
#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>

static void SendDataToFirstInstance(char *singleInstanceUniqueId, char* message) {
    // Pass message via object, not userInfo: sandboxed apps cannot post distributed
    // notifications with userInfo (the entire post call is blocked by the sandbox).
    NSString *messageStr = [NSString stringWithUTF8String:message];
    [[NSDistributedNotificationCenter defaultCenter]
        postNotificationName:[NSString stringWithUTF8String:singleInstanceUniqueId]
        object:messageStr
        userInfo:nil
        deliverImmediately:YES];
}

static char* getMacOSTempDir() {
    NSString *tempDir = NSTemporaryDirectory();
    return strdup([tempDir UTF8String]);
}

*/
import "C"
import (
	"os"
	"unsafe"
)

type darwinLock struct {
	file     *os.File
	uniqueID string
	manager  *singleInstanceManager
}

func (l *darwinLock) notify(data string) error {
	singleInstanceUniqueId := C.CString(l.uniqueID)
	defer C.free(unsafe.Pointer(singleInstanceUniqueId))
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))

	C.SendDataToFirstInstance(singleInstanceUniqueId, cData)

	os.Exit(l.manager.options.ExitCode)
	return nil
}

//export handleSecondInstanceData
func handleSecondInstanceData(secondInstanceMessage *C.char) {
	message := C.GoString(secondInstanceMessage)
	secondInstanceBuffer <- message
}

// macOSTempDir returns NSTemporaryDirectory() as a Go string.
func macOSTempDir() string {
	cTmpDir := C.getMacOSTempDir()
	defer C.free(unsafe.Pointer(cTmpDir))
	return C.GoString(cTmpDir)
}
