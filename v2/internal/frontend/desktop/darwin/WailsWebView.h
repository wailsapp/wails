#ifndef WailsWebView_h
#define WailsWebView_h

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

@interface WKWebView (DeviceScaleFactorSPI)
- (void)_setOverrideDeviceScaleFactor:(CGFloat)scaleFactor;
@end

@interface WailsWebView : WKWebView
@property bool disableWebViewDragAndDrop;
@property bool enableDragAndDrop;
@end

#endif /* WailsWebView_h */
