//
//  WailsMenu.h
//  test
//
//  Created by Lea Anthony on 25/10/21.
//

#ifndef WailsMenu_h
#define WailsMenu_h

#import <Cocoa/Cocoa.h>
#import "Role.h"
#import "WailsMenu.h"
#import "WailsContext.h"

@interface WailsMenu : NSMenu

//- (void) AddMenuByRole :(Role)role;
- (WailsMenu*) initWithNSTitle :(NSString*)title;
- (void) appendSubmenu :(WailsMenu*)child;
- (void) appendRole :(WailsContext*)ctx :(Role)role;

- (NSMenuItem*) newMenuItemWithContext :(WailsContext*)ctx :(NSString*)title :(SEL)selector :(NSString*)key :(NSEventModifierFlags)flags;
- (void*) AppendMenuItem :(WailsContext*)ctx :(NSString*)label :(NSString *)shortcutKey :(int)modifiers :(bool)disabled :(bool)checked :(int)menuItemID;
- (void) AppendSeparator;

@end


#endif /* WailsMenu_h */
