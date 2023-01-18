//go:build darwin

#import <Foundation/Foundation.h>

#import "menuitem.h"

@implementation MenuItem

- (void) handleClick {
    processMenuItemClick(self.menuItemID);
}

@end
