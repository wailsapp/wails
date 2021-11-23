//
//  WindowDelegate.h
//  test
//
//  Created by Lea Anthony on 10/10/21.
//

#ifndef WindowDelegate_h
#define WindowDelegate_h

#import "WailsContext.h"

@interface WindowDelegate : NSObject <NSWindowDelegate>

@property bool hideOnClose;

@property (assign) WailsContext* ctx;

- (void)windowDidExitFullScreen:(NSNotification *)notification;


@end


#endif /* WindowDelegate_h */
