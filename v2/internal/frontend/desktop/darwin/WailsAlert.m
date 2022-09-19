//go:build darwin
//
//  WailsAlert.m
//  test
//
//  Created by Lea Anthony on 20/10/21.
//

#import <Foundation/Foundation.h>

#import "WailsAlert.h"

@implementation WailsAlert

- (void)addButton:(NSString*)text :(NSString*)defaultButton :(NSString*)cancelButton {
    if( text == nil ) {
        return;
    }
    NSButton *button = [self addButtonWithTitle:text];
    if( defaultButton != nil && [text isEqualToString:defaultButton]) {
        [button setKeyEquivalent:@"\r"];
    } else if( cancelButton != nil && [text isEqualToString:cancelButton]) {
        [button setKeyEquivalent:@"\033"];
    } else {
        [button setKeyEquivalent:@""];
    }
}

@end


