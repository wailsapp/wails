//go:build darwin && !ios

#import <Foundation/Foundation.h>

#import "menuitem_darwin.h"

@implementation MenuItem

- (void) handleClick {
    processMenuItemClick(self.menuItemID);
}

@end
