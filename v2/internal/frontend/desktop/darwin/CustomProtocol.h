#ifndef CustomProtocol_h
#define CustomProtocol_h

#import <Cocoa/Cocoa.h>

extern void HandleCustomProtocol(char*);

@interface CustomProtocolSchemeHandler : NSObject
+ (void)handleGetURLEvent:(NSAppleEventDescriptor *)event withReplyEvent:(NSAppleEventDescriptor *)replyEvent;
@end

void StartCustomProtocolHandler(void);

#endif /* CustomProtocol_h */
