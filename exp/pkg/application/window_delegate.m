//go:build darwin
//
//  WindowDelegate.m
//  test
//
//  Created by Lea Anthony on 10/10/21.
//

#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>
#import "window_delegate.h"

@implementation WindowDelegate
- (BOOL)windowShouldClose:(NSWindow *)sender {
    if( self.hideOnClose ) {
        [NSApp hide:nil];
        return false;
    }
    return true;
}


@end
