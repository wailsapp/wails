package xattr

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation
#import <sys/xattr.h>
#include <stdlib.h>
#include <string.h>

void removeXattrFrom(const char *path) {
    ssize_t xattrNamesSize = listxattr(path, NULL, 0, XATTR_NOFOLLOW);
    if (xattrNamesSize <= 0) return;

    char *xattrNames = (char *)malloc(xattrNamesSize);
    xattrNamesSize = listxattr(path, xattrNames, xattrNamesSize, XATTR_NOFOLLOW);

    ssize_t pos = 0;
    while (pos < xattrNamesSize) {
        char *name = xattrNames + pos;
        removexattr(path, name, XATTR_NOFOLLOW);
        pos += strlen(name) + 1;
    }

    free(xattrNames);
}
*/
import "C"
import "unsafe"

func RemoveXAttr(filepath string) {
	cpath := C.CString(filepath)
	defer C.free(unsafe.Pointer(cpath))
	C.removeXattrFrom(cpath)
}
