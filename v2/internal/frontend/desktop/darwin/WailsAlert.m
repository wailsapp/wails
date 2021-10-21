//
//  WailsAlert.m
//  test
//
//  Created by Lea Anthony on 20/10/21.
//

#import <Foundation/Foundation.h>

#import "WailsAlert.h"

@implementation WailsAlert

- (void)addButton:(const char*)text :(const char*)defaultButton :(const char*)cancelButton {
    if( text == nil ) {
        return;
    }
    NSButton *button = [self addButtonWithTitle:[NSString stringWithUTF8String:text]];
    if( defaultButton != nil && strcmp(text, defaultButton) == 0) {
        [button setKeyEquivalent:@"\r"];
    } else if( cancelButton != nil && strcmp(text, cancelButton) == 0) {
        [button setKeyEquivalent:@"\033"];
    } else {
        [button setKeyEquivalent:@""];
    }
}

@end


