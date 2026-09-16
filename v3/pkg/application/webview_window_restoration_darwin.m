//go:build darwin && !ios && !server

#import "webview_window_restoration_darwin.h"
#import <WebKit/WebKit.h>
#import <objc/runtime.h>
#include <stdlib.h>
#include <string.h>

extern void processMacWindowRestore(unsigned long long requestID, const char* identifier, const char* json);
extern void processMacWindowRestorationVisible(void* nsWindow);

static char WailsRestorationDataAssociationKey;
static char WailsRestorationCheckedAssociationKey;

static NSMutableDictionary<NSNumber*, id>* wailsRestorationPending(void) {
    static NSMutableDictionary<NSNumber*, id>* pending = nil;
    if (pending == nil) pending = [[NSMutableDictionary alloc] init];
    return pending;
}

static unsigned long long wailsRestorationNextRequestID = 0;

// The WebviewWindow class declares a webView property; the NativeWindow
// class has none. Resolving it dynamically keeps this file free of the
// WebView-only headers so it also compiles in a wails_native build.
static WKWebView* restorationWebView(void* nsWindow) {
    if (nsWindow == NULL) return nil;
    id window = (id)nsWindow;
    if (![window respondsToSelector:@selector(webView)]) return nil;
    id view = [window valueForKey:@"webView"];
    if (![view isKindOfClass:[WKWebView class]]) return nil;
    return (WKWebView*)view;
}

@implementation WailsWindowRestoration

+ (void)restoreWindowWithIdentifier:(NSUserInterfaceItemIdentifier)identifier
                              state:(NSCoder*)state
                  completionHandler:(void (^)(NSWindow*, NSError*))completionHandler {
    NSString* json = nil;
    if (state.allowsKeyedCoding) {
        NSString* key = [NSString stringWithUTF8String:WailsRestorationDataKey];
        json = [state decodeObjectOfClass:[NSString class] forKey:key];
    }
    unsigned long long requestID = ++wailsRestorationNextRequestID;
    wailsRestorationPending()[@(requestID)] = [[completionHandler copy] autorelease];
    processMacWindowRestore(requestID, identifier.UTF8String, json == nil ? "" : json.UTF8String);
}

@end

// window:willEncodeRestorableState: is installed on the WebviewWindowDelegate
// class at load time (see wailsRestorationInstall) so the delegate's own
// file stays untouched. It writes the JSON stored with
// windowRestorationSetData under WailsRestorationDataKey.
static void wailsRestorationWillEncode(id self, SEL _cmd, NSWindow* window, NSCoder* coder) {
    NSString* json = objc_getAssociatedObject(window, &WailsRestorationDataAssociationKey);
    if (json == nil || !coder.allowsKeyedCoding) return;
    [coder encodeObject:json forKey:[NSString stringWithUTF8String:WailsRestorationDataKey]];
}

// Windows created with MacWindow.RestorationID, or given one before their
// native window existed, are configured the first time AppKit reports them
// on screen: the creation path does not call into this file, and AppKit
// only saves windows that are ordered in. The first of several cheap
// notifications (occlusion change, becoming key or main, update) triggers
// the check; each window is asked once and marked.
@interface WailsRestorationObserver : NSObject
- (void)windowDidAppear:(NSNotification*)notification;
@end

@implementation WailsRestorationObserver
- (void)windowDidAppear:(NSNotification*)notification {
    NSWindow* window = notification.object;
    if (![window isKindOfClass:[NSWindow class]]) return;
    if (objc_getAssociatedObject(window, &WailsRestorationCheckedAssociationKey) != nil) return;
    if (!window.isVisible) return;
    objc_setAssociatedObject(window, &WailsRestorationCheckedAssociationKey, @YES, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    processMacWindowRestorationVisible((void*)window);
}
@end

__attribute__((constructor))
static void wailsRestorationInstall(void) {
    @autoreleasepool {
        Class delegateClass = NSClassFromString(@"WebviewWindowDelegate");
        if (delegateClass != Nil) {
            SEL selector = @selector(window:willEncodeRestorableState:);
            if (class_getInstanceMethod(delegateClass, selector) == NULL) {
                class_addMethod(delegateClass, selector, (IMP)wailsRestorationWillEncode, "v@:@@");
            }
        }
        static WailsRestorationObserver* observer = nil;
        if (observer == nil) {
            observer = [[WailsRestorationObserver alloc] init];
            NSNotificationCenter* center = [NSNotificationCenter defaultCenter];
            for (NSNotificationName name in @[NSWindowDidChangeOcclusionStateNotification,
                                              NSWindowDidBecomeKeyNotification,
                                              NSWindowDidBecomeMainNotification,
                                              NSWindowDidUpdateNotification]) {
                [center addObserver:observer selector:@selector(windowDidAppear:) name:name object:nil];
            }
        }
    }
}

void windowRestorationConfigure(void* nsWindow, const char* identifier) {
    if (nsWindow == NULL) return;
    NSWindow* window = (NSWindow*)nsWindow;
    // Configured explicitly: the visibility observer need not ask Go again.
    objc_setAssociatedObject(window, &WailsRestorationCheckedAssociationKey, @YES, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    if (identifier == NULL || strlen(identifier) == 0) {
        window.restorationClass = Nil;
        window.identifier = nil;
        window.restorable = NO;
        return;
    }
    window.identifier = [NSString stringWithUTF8String:identifier];
    window.restorationClass = [WailsWindowRestoration class];
    window.restorable = YES;
    [window invalidateRestorableState];
}

void windowRestorationSetData(void* nsWindow, const char* json) {
    if (nsWindow == NULL) return;
    NSWindow* window = (NSWindow*)nsWindow;
    NSString* value = nil;
    if (json != NULL && strlen(json) > 0) value = [NSString stringWithUTF8String:json];
    objc_setAssociatedObject(window, &WailsRestorationDataAssociationKey, value, OBJC_ASSOCIATION_COPY_NONATOMIC);
    [window invalidateRestorableState];
}

char* windowRestorationData(void* nsWindow) {
    if (nsWindow == NULL) return NULL;
    NSString* value = objc_getAssociatedObject((NSWindow*)nsWindow, &WailsRestorationDataAssociationKey);
    if (value == nil) return NULL;
    return strdup(value.UTF8String);
}

void windowRestorationComplete(unsigned long long requestID, void* nsWindow, const char* error) {
    NSNumber* key = @(requestID);
    void (^completion)(NSWindow*, NSError*) = wailsRestorationPending()[key];
    if (completion == nil) return;
    [[completion retain] autorelease];
    [wailsRestorationPending() removeObjectForKey:key];
    if (nsWindow != NULL) {
        completion((NSWindow*)nsWindow, nil);
        return;
    }
    NSString* message = error != NULL ? [NSString stringWithUTF8String:error] : @"window restoration declined";
    NSError* nsError = [NSError errorWithDomain:@"io.wails.restoration" code:1
                                       userInfo:@{NSLocalizedDescriptionKey: message}];
    completion(nil, nsError);
}

int windowRestorationInteractionState(void* nsWindow, void** bytes, int* length) {
    *bytes = NULL;
    *length = 0;
    WKWebView* webView = restorationWebView(nsWindow);
    if (webView == nil) return WailsRestorationStateNoWebView;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 120000
    if (@available(macOS 12.0, *)) {
        id state = webView.interactionState;
        if (![state isKindOfClass:[NSData class]]) return WailsRestorationStateOK;
        NSData* data = (NSData*)state;
        if (data.length == 0) return WailsRestorationStateOK;
        void* copy = malloc(data.length);
        if (copy == NULL) return WailsRestorationStateInvalid;
        memcpy(copy, data.bytes, data.length);
        *bytes = copy;
        *length = (int)data.length;
        return WailsRestorationStateOK;
    }
#endif
    return WailsRestorationStateUnsupported;
}

int windowRestorationSetInteractionState(void* nsWindow, const void* bytes, int length) {
    WKWebView* webView = restorationWebView(nsWindow);
    if (webView == nil) return WailsRestorationStateNoWebView;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 120000
    if (@available(macOS 12.0, *)) {
        if (bytes == NULL || length <= 0) return WailsRestorationStateInvalid;
        NSData* data = [NSData dataWithBytes:bytes length:(NSUInteger)length];
        @try {
            webView.interactionState = data;
        } @catch (NSException* exception) {
            return WailsRestorationStateInvalid;
        }
        return WailsRestorationStateOK;
    }
#endif
    return WailsRestorationStateUnsupported;
}
