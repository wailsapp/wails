//go:build darwin


#ifndef WindowDelegate_h
#define WindowDelegate_h

#import <Cocoa/Cocoa.h>

@interface WindowDelegate : NSObject <NSWindowDelegate>

@property bool hideOnClose;

@end


#endif /* WindowDelegate_h */
