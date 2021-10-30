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
- (void)addButton:(NSString*)text :(NSString*)defaultButton :(NSString*)cancelButton;
@end


#endif /* WailsAlert_h */
