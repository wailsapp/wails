//go:build darwin && !ios && !server

#ifndef WailsPresentationOptions_h
#define WailsPresentationOptions_h

#import <Cocoa/Cocoa.h>

// NSApplication.presentationOptions. Both functions expect the main thread.
// presentationOptionsSet returns NULL on success or a malloc'd message when
// AppKit raised for the combination; the caller frees it. Combinations are
// validated in Go first, so a message here means the Go rules and AppKit's
// disagree.
char* presentationOptionsSet(unsigned long options);
unsigned long presentationOptionsGet(void);

#endif
