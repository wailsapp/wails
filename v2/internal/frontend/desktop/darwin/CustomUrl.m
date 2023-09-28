#include "CustomUrl.h"

@implementation CustomUrlSchemeHandler
+ (void)handleGetURLEvent:(NSAppleEventDescriptor *)event withReplyEvent:(NSAppleEventDescriptor *)replyEvent {
  [event paramDescriptorForKeyword:keyDirectObject];

   NSString *urlStr = [[event paramDescriptorForKeyword:keyDirectObject] stringValue];

   HandleCustomURL((char*)[[[event paramDescriptorForKeyword:keyDirectObject] stringValue] UTF8String]);
}
@end

void StartCustomURLHandler(void) {
	NSAppleEventManager *appleEventManager = [NSAppleEventManager sharedAppleEventManager];

	[appleEventManager setEventHandler:[CustomUrlSchemeHandler class]
	    andSelector:@selector(handleGetURLEvent:withReplyEvent:)
	    forEventClass:kInternetEventClass
	    andEventID: kAEGetURL];
}
