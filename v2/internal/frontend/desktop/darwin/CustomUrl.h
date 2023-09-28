#ifndef CustomUrl_h
#define CustomUrl_h

#import <Cocoa/Cocoa.h>

extern void HandleCustomURL(char*);

@interface CustomUrlSchemeHandler : NSObject
+ (void)handleGetURLEvent:(NSAppleEventDescriptor *)event withReplyEvent:(NSAppleEventDescriptor *)replyEvent;
@end

void StartCustomURLHandler(void);

#endif /* CustomUrl_h */
