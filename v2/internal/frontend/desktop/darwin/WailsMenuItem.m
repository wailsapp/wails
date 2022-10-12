//go:build darwin
//
//  WailsMenuItem.m
//  test
//
//  Created by Lea Anthony on 27/10/21.
//

#import <Foundation/Foundation.h>

#import "WailsMenuItem.h"
#include "message.h"


@implementation WailsMenuItem

- (void) handleClick {
    processCallback(self.menuItemID);
}

@end
