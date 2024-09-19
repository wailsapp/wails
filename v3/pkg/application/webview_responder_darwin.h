#ifndef WEBVIEW_RESPONDER_DARWIN
#define WEBVIEW_RESPONDER_DARWIN

#import <Cocoa/Cocoa.h>

@interface WebviewResponder : NSResponder
@property(assign) NSWindow *w;

- (WebviewResponder *) initAttachToWindow:(NSWindow *)window;
- (void)keyDown:(NSEvent *)event;
- (NSString *)keyStringFromEvent:(NSEvent *)event;
- (NSString *)specialKeyStringFromEvent:(NSEvent *)event;
@end

#endif /* WEBVIEW_RESPONDER_DARWIN */
