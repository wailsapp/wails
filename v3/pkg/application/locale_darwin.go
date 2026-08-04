//go:build darwin && !ios && !server

package application

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation
#import <Foundation/Foundation.h>
#include <stdlib.h>

static const char* getSystemLocale() {
    // Use languageIdentifier (macOS 13+ / iOS 16+) for clean BCP-47 output.
    // localeIdentifier can include POSIX modifiers like @collation=pinyin.
    if (@available(macOS 13, *)) {
        NSString *tag = [[NSLocale currentLocale] languageIdentifier];
        return strdup([tag UTF8String]);
    }
    // Fallback for macOS < 13: build from components to preserve script subtags
    // (e.g. zh-Hant-TW, not just zh-TW).
    NSLocale *locale = [NSLocale currentLocale];
    NSString *lang = [locale languageCode];
    NSString *script = [locale scriptCode];
    NSString *country = [locale countryCode];
    NSMutableString *tag = [NSMutableString stringWithString:lang ?: @"en"];
    if (script.length > 0) [tag appendFormat:@"-%@", script];
    if (country.length > 0) [tag appendFormat:@"-%@", country];
    return strdup([tag UTF8String]);
}
*/
import "C"
import "unsafe"

// SystemLocale returns the system's configured locale as a BCP-47 language tag
// (e.g. "nb-NO", "en-US").
func SystemLocale() string {
	cStr := C.getSystemLocale()
	if cStr == nil {
		return "en"
	}
	defer C.free(unsafe.Pointer(cStr))
	return C.GoString(cStr)
}
