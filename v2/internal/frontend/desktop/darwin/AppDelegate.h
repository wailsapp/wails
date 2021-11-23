//
//  AppDelegate.h
//  test
//
//  Created by Lea Anthony on 10/10/21.
//

#ifndef AppDelegate_h
#define AppDelegate_h

#import <Cocoa/Cocoa.h>
#import "WailsContext.h"

@interface AppDelegate : NSResponder <NSTouchBarProvider>

@property bool alwaysOnTop;
@property bool startHidden;
@property bool startFullscreen;
@property (retain) WailsWindow* mainWindow;

@end

#endif /* AppDelegate_h */
