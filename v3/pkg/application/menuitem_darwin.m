//go:build darwin && !ios && !server && !purego

#import <Foundation/Foundation.h>

#import "menuitem_darwin.h"

@implementation MenuItem

- (void) handleClick {
    processMenuItemClick(self.menuItemID);
}

@end
