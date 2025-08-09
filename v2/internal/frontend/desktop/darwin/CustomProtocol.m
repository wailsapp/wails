#include "CustomProtocol.h"

@implementation CustomProtocolSchemeHandler
+ (void)handleGetURLEvent:(NSAppleEventDescriptor *)event withReplyEvent:(NSAppleEventDescriptor *)replyEvent {
  [event paramDescriptorForKeyword:keyDirectObject];

   NSString *urlStr = [[event paramDescriptorForKeyword:keyDirectObject] stringValue];

   HandleCustomProtocol((char*)[[[event paramDescriptorForKeyword:keyDirectObject] stringValue] UTF8String]);
}
@end

void StartCustomProtocolHandler(void) {
	NSAppleEventManager *appleEventManager = [NSAppleEventManager sharedAppleEventManager];

	[appleEventManager setEventHandler:[CustomProtocolSchemeHandler class]
	    andSelector:@selector(handleGetURLEvent:withReplyEvent:)
	    forEventClass:kInternetEventClass
	    andEventID: kAEGetURL];
}
