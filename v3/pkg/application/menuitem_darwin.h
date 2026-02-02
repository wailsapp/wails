#ifndef MenuItemDelegate_h
#define MenuItemDelegate_h

#import <Cocoa/Cocoa.h>

extern void processMenuItemClick(unsigned int);

@interface MenuItem : NSMenuItem

@property unsigned int menuItemID;

- (void) handleClick;

@end


#endif /* MenuItemDelegate_h */
