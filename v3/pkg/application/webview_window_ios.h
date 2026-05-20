//go:build ios

#ifndef WEBVIEW_WINDOW_IOS_H
#define WEBVIEW_WINDOW_IOS_H

#import <UIKit/UIKit.h>
#import <WebKit/WebKit.h>

#import "application_ios.h"

// Globals provided by application_ios.m
extern WailsAppDelegate *appDelegate;
extern unsigned int nextWindowID;

// Scheme handler bridging to Go asset server
@interface WailsSchemeHandler : NSObject <WKURLSchemeHandler>
@property (nonatomic, assign) unsigned int windowID;
- (instancetype)initWithWindowID:(unsigned int)windowID;
@end

// Script message handler bridging JS -> Go
@interface WailsMessageHandler : NSObject <WKScriptMessageHandler>
@property (nonatomic, assign) unsigned int windowID;
- (instancetype)initWithWindowID:(unsigned int)windowID;
@end

// Main view controller owning the WKWebView
@interface WailsViewController : UIViewController <WKNavigationDelegate, UITabBarDelegate>
@property (nonatomic, strong) WKWebView *webView;
@property (nonatomic, strong) WailsSchemeHandler *schemeHandler;
@property (nonatomic, strong) WailsMessageHandler *messageHandler;
@property (nonatomic, assign) unsigned int windowID;
- (void)enableNativeTabs:(BOOL)enabled;
- (void)selectNativeTabIndex:(NSInteger)index;
@property (nonatomic, strong) UITabBar *tabBar;
- (instancetype)initWithWindowID:(unsigned int)windowID;
- (void)executeJavaScript:(NSString *)js;
@end

// C-callable bridge used by Go
unsigned int ios_create_webview(void);
void* ios_create_webview_with_id(unsigned int wailsID);
void ios_execute_javascript(unsigned int windowID, const char* js);
void ios_window_exec_js(void* viewController, const char* js);
void ios_window_load_url(void* viewController, const char* url);
void ios_window_set_html(void* viewController, const char* html);
unsigned int ios_window_get_id(void* viewController);
void ios_window_release_handle(void* viewController);

// Console logging bridge (broadcast to all WKWebViews)
void ios_console_log(const char* level, const char* message);

// Set background color (RGBA 0-255) for a specific window's root view and webview
void ios_window_set_background_color(void* viewController, unsigned char r, unsigned char g, unsigned char b, unsigned char a);

#endif /* WEBVIEW_WINDOW_IOS_H */