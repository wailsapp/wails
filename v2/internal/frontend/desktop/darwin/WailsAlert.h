//
//  WailsAlert.h
//  test
//
//  Created by Lea Anthony on 20/10/21.
//

#ifndef WailsAlert_h
#define WailsAlert_h

#import <Cocoa/Cocoa.h>

@interface WailsAlert : NSAlert 
- (void)addButton:(const char*)text :(const char*)defaultButton :(const char*)cancelButton;
@end


#endif /* WailsAlert_h */
