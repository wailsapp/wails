//go:build darwin && !ios && !server && (!production || devtools)

package application

/*
#cgo CFLAGS: -x objective-c
#import <WebKit/WebKit.h>
static void embeddedPanelDevTools(void* panel, bool enabled) {
 WKWebView* view = (WKWebView*)panel;
 [view.configuration.preferences setValue:@(enabled) forKey:@"developerExtrasEnabled"];
}
static void embeddedPanelInspector(void* panel) {
 WKWebView* view = (WKWebView*)panel;
 @try {
  id inspector = [view valueForKey:@"_inspector"];
  if ([inspector respondsToSelector:@selector(show)]) [inspector performSelector:@selector(show)];
 } @catch (NSException* exception) {
  NSLog(@"Opening panel inspector failed: %@", exception.reason);
 }
}
*/
import "C"
import "unsafe"

func configureEmbeddedPanelDevTools(view unsafe.Pointer, enabled bool) {
	C.embeddedPanelDevTools(view, C.bool(enabled))
}
func openEmbeddedPanelDevTools(view unsafe.Pointer) { C.embeddedPanelInspector(view) }
