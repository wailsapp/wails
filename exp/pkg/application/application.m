//go:build darwin

#import "app_delegate.h"

AppDelegate *appDelegate = nil;

void Init(void) {
    [NSApplication sharedApplication];
    appDelegate = [[AppDelegate alloc] init];
    [NSApp setDelegate:appDelegate];
}

void SetActivationPolicy(int policy) {
    [appDelegate setApplicationActivationPolicy:policy];
}

void Run(void) {
    @autoreleasepool {
        [NSApp run];
        [appDelegate release];
    }
}
