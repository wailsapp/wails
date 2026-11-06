//go:build darwin && !server && (ios || !purego)

package application

/*
#cgo LDFLAGS: -framework CoreFoundation
#include <CoreFoundation/CoreFoundation.h>
#include <stdlib.h>

// wailsPreferredLanguages returns the user's preferred languages, most
// preferred first, joined by newlines, as a malloc'd UTF-8 string the caller
// frees, or NULL when the list is unavailable.
//
// CFLocaleCopyPreferredLanguages is the list in Language & Region on macOS and
// iOS, the same as +[NSLocale preferredLanguages]. The current locale
// (CFLocaleCopyCurrent, -[NSLocale currentLocale]) is the regional format
// instead, which can differ from the display language.
static char* wailsPreferredLanguages(void) {
	CFArrayRef langs = CFLocaleCopyPreferredLanguages();
	if (langs == NULL) {
		return NULL;
	}
	CFStringRef joined = CFStringCreateByCombiningStrings(kCFAllocatorDefault, langs, CFSTR("\n"));
	CFRelease(langs);
	if (joined == NULL) {
		return NULL;
	}
	CFIndex size = CFStringGetMaximumSizeForEncoding(CFStringGetLength(joined), kCFStringEncodingUTF8) + 1;
	char *buf = malloc(size);
	if (buf != NULL && !CFStringGetCString(joined, buf, size, kCFStringEncodingUTF8)) {
		free(buf);
		buf = NULL;
	}
	CFRelease(joined);
	return buf;
}
*/
import "C"

import (
	"strings"
	"unsafe"
)

// platformPreferredLanguages returns the list in Language & Region on macOS
// (System Settings › General) and iOS (Settings › General).
func platformPreferredLanguages() []string {
	joined := C.wailsPreferredLanguages()
	if joined == nil {
		return nil
	}
	defer C.free(unsafe.Pointer(joined))
	return strings.Split(C.GoString(joined), "\n")
}
