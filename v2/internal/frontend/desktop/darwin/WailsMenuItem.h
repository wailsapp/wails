//
//  WailsMenuItem.h
//  test
//
//  Created by Lea Anthony on 27/10/21.
//

#ifndef WailsMenuItem_h
#define WailsMenuItem_h

#import <Cocoa/Cocoa.h>

@interface WailsMenuItem : NSMenuItem

@property int menuItemID;

- (void) handleClick;

@end


#endif /* WailsMenuItem_h */
