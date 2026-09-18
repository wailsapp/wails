//go:build darwin && !ios && !server

#import "presentation_options_darwin.h"
#include <stdlib.h>
#include <string.h>

char* presentationOptionsSet(unsigned long options) {
    @try {
        NSApp.presentationOptions = (NSApplicationPresentationOptions)options;
    } @catch (NSException* exception) {
        const char* reason = exception.reason.UTF8String;
        if (reason == NULL) reason = exception.name.UTF8String;
        if (reason == NULL) reason = "AppKit rejected the presentation options";
        return strdup(reason);
    }
    return NULL;
}

unsigned long presentationOptionsGet(void) {
    return (unsigned long)NSApp.presentationOptions;
}
