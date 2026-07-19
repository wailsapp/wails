#import "webcontentsview_darwin.h"
#import <objc/runtime.h>

extern void browserViewSnapshotCallback(uintptr_t callbackID, const char* base64);
extern void browserViewAutomationCommandCallback(uintptr_t callbackID, const char* resultJSON, const char* errorMessage);
extern void browserViewAutomationEventCallback(uintptr_t viewID, const char* method, const char* payloadJSON);

static const void* WailsAutomationBridgeAssociationKey = &WailsAutomationBridgeAssociationKey;

static int64_t WailsAutomationTimestamp(void) {
    return (int64_t)([[NSDate date] timeIntervalSince1970] * 1000.0);
}

static NSString* WailsAutomationURLString(WKWebView* webView) {
    if (webView.URL == nil) {
        return @"";
    }
    return webView.URL.absoluteString ?: @"";
}

static NSString* WailsAutomationString(id value) {
    if (value == nil || value == [NSNull null]) {
        return @"";
    }
    if ([value isKindOfClass:[NSString class]]) {
        return value;
    }
    if ([value respondsToSelector:@selector(stringValue)]) {
        return [value stringValue];
    }
    return [value description] ?: @"";
}

static NSString* WailsAutomationCookieSameSite(NSHTTPCookie* cookie) {
    if (@available(macOS 10.15, *)) {
        NSString* sameSite = cookie.properties[NSHTTPCookieSameSitePolicy];
        return sameSite ?: @"";
    }
    return @"";
}

static NSDictionary* WailsAutomationCookieDictionary(NSHTTPCookie* cookie) {
    NSMutableDictionary* result = [NSMutableDictionary dictionaryWithDictionary:@{
        @"name": cookie.name ?: @"",
        @"value": cookie.value ?: @"",
        @"domain": cookie.domain ?: @"",
        @"path": cookie.path ?: @"/",
        @"secure": @(cookie.secure),
        @"httpOnly": @(cookie.HTTPOnly),
        @"sessionOnly": @(cookie.sessionOnly),
    }];
    if (cookie.expiresDate != nil) {
        result[@"expires"] = @((int64_t)([cookie.expiresDate timeIntervalSince1970] * 1000.0));
    }
    NSString* sameSite = WailsAutomationCookieSameSite(cookie);
    if (sameSite.length > 0) {
        result[@"sameSite"] = sameSite;
    }
    return result;
}

static WKHTTPCookieStore* WailsAutomationCookieStore(WKWebView* webView) {
    return webView.configuration.websiteDataStore.httpCookieStore;
}

static id WailsAutomationSanitizeJSONValue(id value) {
    if (value == nil) {
        return [NSNull null];
    }
    if ([value isKindOfClass:[NSString class]] || [value isKindOfClass:[NSNumber class]] || [value isKindOfClass:[NSNull class]]) {
        return value;
    }
    if ([value isKindOfClass:[NSDate class]]) {
        return @((int64_t)([(NSDate*)value timeIntervalSince1970] * 1000.0));
    }
    if ([value isKindOfClass:[NSURL class]]) {
        return [(NSURL*)value absoluteString] ?: @"";
    }
    if ([value isKindOfClass:[NSArray class]]) {
        NSMutableArray* result = [NSMutableArray arrayWithCapacity:[(NSArray*)value count]];
        for (id item in (NSArray*)value) {
            [result addObject:WailsAutomationSanitizeJSONValue(item)];
        }
        return result;
    }
    if ([value isKindOfClass:[NSDictionary class]]) {
        NSMutableDictionary* result = [NSMutableDictionary dictionaryWithCapacity:[(NSDictionary*)value count]];
        for (id key in (NSDictionary*)value) {
            NSString* stringKey = [key isKindOfClass:[NSString class]] ? key : [key description];
            result[stringKey] = WailsAutomationSanitizeJSONValue([(NSDictionary*)value objectForKey:key]);
        }
        return result;
    }
    return [value description] ?: @"";
}

static NSString* WailsAutomationJSONString(id value) {
    id sanitized = WailsAutomationSanitizeJSONValue(value);
    NSError* error = nil;
    NSData* jsonData = nil;

    if ([sanitized isKindOfClass:[NSDictionary class]] || [sanitized isKindOfClass:[NSArray class]]) {
        jsonData = [NSJSONSerialization dataWithJSONObject:sanitized options:0 error:&error];
    } else {
        jsonData = [NSJSONSerialization dataWithJSONObject:@[sanitized ?: [NSNull null]] options:0 error:&error];
        if (jsonData != nil) {
            NSString* wrappedJSON = [[NSString alloc] initWithData:jsonData encoding:NSUTF8StringEncoding];
            if (wrappedJSON.length >= 2) {
                return [wrappedJSON substringWithRange:NSMakeRange(1, wrappedJSON.length - 2)];
            }
        }
    }

    if (jsonData == nil || error != nil) {
        return @"null";
    }

    return [[NSString alloc] initWithData:jsonData encoding:NSUTF8StringEncoding];
}

static void WailsAutomationDispatchEvent(uintptr_t viewID, NSString* method, id payload) {
    NSString* payloadJSON = WailsAutomationJSONString(payload);
    browserViewAutomationEventCallback(viewID, method.UTF8String, payloadJSON.UTF8String);
}

static void WailsAutomationCompleteCommand(uintptr_t callbackID, id result, NSError* error) {
    if (error != nil) {
        NSString* message = error.localizedDescription ?: @"JavaScript execution failed";
        browserViewAutomationCommandCallback(callbackID, NULL, message.UTF8String);
        return;
    }

    NSString* resultJSON = WailsAutomationJSONString(result);
    browserViewAutomationCommandCallback(callbackID, resultJSON.UTF8String, NULL);
}

static void WailsAutomationCompleteCommandMessage(uintptr_t callbackID, NSString* message) {
    browserViewAutomationCommandCallback(callbackID, NULL, message.UTF8String);
}

static void WailsRunOnMainThread(void (^block)(void)) {
    if ([NSThread isMainThread]) {
        block();
        return;
    }
    dispatch_sync(dispatch_get_main_queue(), block);
}

static NSString* WailsPageEvaluationFunctionBody(void) {
    return @"const serializeValue = (value) => {"
           "if (value === undefined) { return { type: 'undefined' }; }"
           "if (value === null) { return { type: 'object', subtype: 'null', value: null }; }"
           "const type = typeof value;"
           "if (type === 'string' || type === 'number' || type === 'boolean') { return { type, value }; }"
           "if (type === 'bigint') { return { type: 'bigint', description: value.toString() }; }"
           "if (type === 'function') { return { type: 'function', description: value.name || 'anonymous' }; }"
           "if (Array.isArray(value)) { try { return { type: 'object', subtype: 'array', value: JSON.parse(JSON.stringify(value)) }; } catch (_) { return { type: 'object', subtype: 'array', description: 'Array(' + value.length + ')' }; } }"
           "if (value instanceof Element) { return { type: 'object', subtype: 'node', description: value.tagName ? value.tagName.toLowerCase() : 'element', value: { tagName: value.tagName ? value.tagName.toLowerCase() : '', id: value.id || '', text: value.innerText || value.textContent || '' } }; }"
           "try { return { type: 'object', value: JSON.parse(JSON.stringify(value)) }; } catch (_) { return { type: 'object', description: Object.prototype.toString.call(value) }; }"
           "};"
           "let value = (0, eval)(expression);"
           "if (awaitPromise) { value = await Promise.resolve(value); }"
           "return serializeValue(value);";
}

@interface WailsAutomationBridge : NSObject <WKNavigationDelegate, WKUIDelegate, WKScriptMessageHandler>

@property (nonatomic, assign) uintptr_t viewID;
@property (nonatomic, retain) WKContentWorld* automationWorld;
@property (nonatomic, copy) NSString* pageScript;
@property (nonatomic, copy) NSString* automationScript;
@property (nonatomic, assign) BOOL installed;

- (instancetype)initWithViewID:(uintptr_t)viewID pageScript:(NSString*)pageScript automationScript:(NSString*)automationScript;
- (void)installIntoWebView:(WKWebView*)webView;

@end

@implementation WailsAutomationBridge

- (instancetype)initWithViewID:(uintptr_t)viewID pageScript:(NSString*)pageScript automationScript:(NSString*)automationScript {
    self = [super init];
    if (self != nil) {
        self.viewID = viewID;
        self.pageScript = pageScript ?: @"";
        self.automationScript = automationScript ?: @"";
        self.automationWorld = [WKContentWorld worldWithName:@"WailsAutomation"];
    }
    return self;
}

- (void)installIntoWebView:(WKWebView*)webView {
    if (self.installed) {
        return;
    }

    WKUserContentController* controller = webView.configuration.userContentController;
    [controller addScriptMessageHandler:self contentWorld:WKContentWorld.pageWorld name:@"wailsAutomationPage"];

    WKUserScript* pageScript = [[WKUserScript alloc] initWithSource:self.pageScript
                                                      injectionTime:WKUserScriptInjectionTimeAtDocumentStart
                                                   forMainFrameOnly:NO
                                                     inContentWorld:WKContentWorld.pageWorld];
    [controller addUserScript:pageScript];

    WKUserScript* automationScript = [[WKUserScript alloc] initWithSource:self.automationScript
                                                            injectionTime:WKUserScriptInjectionTimeAtDocumentStart
                                                         forMainFrameOnly:NO
                                                           inContentWorld:self.automationWorld];
    [controller addUserScript:automationScript];

    webView.navigationDelegate = self;
    webView.UIDelegate = self;
    self.installed = YES;
}

- (void)emitMethod:(NSString*)method payload:(NSDictionary*)payload {
    WailsAutomationDispatchEvent(self.viewID, method, payload);
}

- (void)userContentController:(WKUserContentController*)userContentController didReceiveScriptMessage:(WKScriptMessage*)message {
    if (![message.body isKindOfClass:[NSDictionary class]]) {
        return;
    }

    NSDictionary* body = (NSDictionary*)message.body;
    NSString* type = body[@"type"];
    NSDictionary* payload = [body[@"payload"] isKindOfClass:[NSDictionary class]] ? body[@"payload"] : @{};

    if ([type isEqualToString:@"console"]) {
        [self emitMethod:@"Console.messageAdded" payload:payload];
        return;
    }
    if ([type isEqualToString:@"exception"]) {
        [self emitMethod:@"Runtime.exceptionThrown" payload:payload];
        return;
    }
    if ([type isEqualToString:@"networkRequest"]) {
        [self emitMethod:@"Network.requestWillBeSent" payload:payload];
        return;
    }
    if ([type isEqualToString:@"networkResponse"]) {
        [self emitMethod:@"Network.responseReceived" payload:payload];
        return;
    }
    if ([type isEqualToString:@"networkFinished"]) {
        [self emitMethod:@"Network.loadingFinished" payload:payload];
        return;
    }
    if ([type isEqualToString:@"networkFailed"]) {
        [self emitMethod:@"Network.loadingFailed" payload:payload];
        return;
    }
    if ([type isEqualToString:@"domcontentloaded"]) {
        [self emitMethod:@"Page.domContentEventFired" payload:payload];
    }
}

- (void)webView:(WKWebView*)webView didStartProvisionalNavigation:(WKNavigation*)navigation {
    [self emitMethod:@"Page.frameStartedLoading" payload:@{
        @"url": WailsAutomationURLString(webView),
        @"timestamp": @(WailsAutomationTimestamp()),
    }];
}

- (void)webView:(WKWebView*)webView didCommitNavigation:(WKNavigation*)navigation {
    [self emitMethod:@"Page.frameNavigated" payload:@{
        @"url": WailsAutomationURLString(webView),
        @"title": webView.title ?: @"",
        @"timestamp": @(WailsAutomationTimestamp()),
    }];
}

- (void)webView:(WKWebView*)webView didFinishNavigation:(WKNavigation*)navigation {
    [self emitMethod:@"Page.loadEventFired" payload:@{
        @"url": WailsAutomationURLString(webView),
        @"title": webView.title ?: @"",
        @"timestamp": @(WailsAutomationTimestamp()),
    }];
}

- (nullable WKWebView*)webView:(WKWebView*)webView createWebViewWithConfiguration:(WKWebViewConfiguration*)configuration forNavigationAction:(WKNavigationAction*)navigationAction windowFeatures:(WKWindowFeatures*)windowFeatures {
    NSString* url = navigationAction.request.URL.absoluteString ?: @"";
    [self emitMethod:@"Page.windowOpenRequested" payload:@{
        @"url": url,
        @"timestamp": @(WailsAutomationTimestamp()),
    }];
    return nil;
}

- (void)webViewWebContentProcessDidTerminate:(WKWebView*)webView {
    [self emitMethod:@"Page.webContentProcessTerminated" payload:@{
        @"timestamp": @(WailsAutomationTimestamp()),
    }];
}

@end

static WailsAutomationBridge* WailsAutomationBridgeForView(WKWebView* webView) {
    return (WailsAutomationBridge*)objc_getAssociatedObject(webView, WailsAutomationBridgeAssociationKey);
}

void* createWebContentsView(int x, int y, int w, int h, WebContentsViewPreferences prefs) {
    NSRect frame = NSMakeRect(x, y, w, h);
    WKWebViewConfiguration* config = [[WKWebViewConfiguration alloc] init];
    config.userContentController = [[WKUserContentController alloc] init];

    WKPreferences* preferences = [[WKPreferences alloc] init];

    @try {
        if (@available(macOS 10.11, *)) {
            [preferences setValue:@(prefs.devTools) forKey:@"developerExtrasEnabled"];
        }

        if (@available(macOS 11.0, *)) {
            WKWebpagePreferences* webpagePreferences = [[WKWebpagePreferences alloc] init];
            webpagePreferences.allowsContentJavaScript = prefs.javascript;
            config.defaultWebpagePreferences = webpagePreferences;
        } else {
            preferences.javaScriptEnabled = prefs.javascript;
        }

        if (!prefs.webSecurity) {
            @try {
                [config.preferences setValue:@YES forKey:@"allowFileAccessFromFileURLs"];
                [config.preferences setValue:@YES forKey:@"allowUniversalAccessFromFileURLs"];
            } @catch (NSException* e) {
            }
        }
        config.preferences = preferences;
    } @catch (NSException* e) {
    }

    WKWebView* webView = [[WKWebView alloc] initWithFrame:frame configuration:config];
    [webView setValue:@(NO) forKey:@"drawsTransparentBackground"];

    if (prefs.userAgent != NULL) {
        NSString* customUA = [NSString stringWithUTF8String:prefs.userAgent];
        [webView setCustomUserAgent:customUA];
    }

    return webView;
}

void webContentsViewConfigureAutomation(void* view, uintptr_t viewID, const char* pageScript, const char* automationScript) {
    WKWebView* webView = (WKWebView*)view;
    NSString* pageScriptString = pageScript != NULL ? [NSString stringWithUTF8String:pageScript] : @"";
    NSString* automationScriptString = automationScript != NULL ? [NSString stringWithUTF8String:automationScript] : @"";

    WailsRunOnMainThread(^{
        WailsAutomationBridge* bridge = WailsAutomationBridgeForView(webView);
        if (bridge == nil) {
            bridge = [[WailsAutomationBridge alloc] initWithViewID:viewID pageScript:pageScriptString automationScript:automationScriptString];
            objc_setAssociatedObject(webView, WailsAutomationBridgeAssociationKey, bridge, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
        }
        [bridge installIntoWebView:webView];
    });
}

void webContentsViewSetBounds(void* view, int x, int y, int w, int h) {
    WKWebView* webView = (WKWebView*)view;
    NSView* superview = [webView superview];

    if (superview != nil) {
        CGFloat superHeight = superview.bounds.size.height;
        CGFloat cocoaY = superHeight - y - h;
        [webView setFrame:NSMakeRect(x, cocoaY, w, h)];
    } else {
        [webView setFrame:NSMakeRect(x, y, w, h)];
    }
}

void webContentsViewSetURL(void* view, const char* url) {
    WKWebView* webView = (WKWebView*)view;
    NSString* nsURL = [NSString stringWithUTF8String:url];
    NSURLRequest* request = [NSURLRequest requestWithURL:[NSURL URLWithString:nsURL]];
    dispatch_async(dispatch_get_main_queue(), ^{
        [webView loadRequest:request];
    });
}

void webContentsViewExecJS(void* view, const char* js) {
    WKWebView* webView = (WKWebView*)view;
    NSString* script = [NSString stringWithUTF8String:js];
    dispatch_async(dispatch_get_main_queue(), ^{
        [webView evaluateJavaScript:script completionHandler:nil];
    });
}

void webContentsViewAutomationEvaluate(void* view, const char* expression, int world, bool awaitPromise, uintptr_t callbackID) {
    WKWebView* webView = (WKWebView*)view;
    NSString* expressionString = expression != NULL ? [NSString stringWithUTF8String:expression] : @"";

    WailsRunOnMainThread(^{
        WailsAutomationBridge* bridge = WailsAutomationBridgeForView(webView);
        if (bridge == nil) {
            WailsAutomationCompleteCommandMessage(callbackID, @"automation bridge is not configured");
            return;
        }

        WKContentWorld* contentWorld = world == 1 ? bridge.automationWorld : WKContentWorld.pageWorld;
        NSString* functionBody = world == 1
            ? @"return await globalThis.__wailsAutomation.evaluate(expression, awaitPromise);"
            : WailsPageEvaluationFunctionBody();

        [webView callAsyncJavaScript:functionBody
                            arguments:@{
                                @"expression": expressionString,
                                @"awaitPromise": @(awaitPromise),
                            }
                              inFrame:nil
                       inContentWorld:contentWorld
                    completionHandler:^(id result, NSError* error) {
            WailsAutomationCompleteCommand(callbackID, result, error);
        }];
    });
}

void webContentsViewAutomationInvoke(void* view, const char* method, const char* paramsJSON, uintptr_t callbackID) {
    WKWebView* webView = (WKWebView*)view;
    NSString* methodString = method != NULL ? [NSString stringWithUTF8String:method] : @"";

    NSError* jsonError = nil;
    id paramsObject = [NSNull null];
    if (paramsJSON != NULL) {
        NSData* data = [[NSString stringWithUTF8String:paramsJSON] dataUsingEncoding:NSUTF8StringEncoding];
        if (data != nil && data.length > 0) {
            paramsObject = [NSJSONSerialization JSONObjectWithData:data options:0 error:&jsonError];
            if (paramsObject == nil) {
                WailsAutomationCompleteCommandMessage(callbackID, jsonError.localizedDescription ?: @"invalid automation params");
                return;
            }
        }
    }

    WailsRunOnMainThread(^{
        WailsAutomationBridge* bridge = WailsAutomationBridgeForView(webView);
        if (bridge == nil) {
            WailsAutomationCompleteCommandMessage(callbackID, @"automation bridge is not configured");
            return;
        }

        [webView callAsyncJavaScript:@"return await globalThis.__wailsAutomation.dispatch(method, params || {});"
                            arguments:@{
                                @"method": methodString,
                                @"params": paramsObject ?: [NSNull null],
                            }
                              inFrame:nil
                       inContentWorld:bridge.automationWorld
                    completionHandler:^(id result, NSError* error) {
            WailsAutomationCompleteCommand(callbackID, result, error);
        }];
    });
}

void webContentsViewCreatePDF(void* view, uintptr_t callbackID) {
    WKWebView* webView = (WKWebView*)view;
    dispatch_async(dispatch_get_main_queue(), ^{
        if (@available(macOS 11.0, *)) {
            WKPDFConfiguration* config = [[WKPDFConfiguration alloc] init];
            [webView createPDFWithConfiguration:config completionHandler:^(NSData* pdfDocumentData, NSError* error) {
                if (error != nil || pdfDocumentData == nil) {
                    WailsAutomationCompleteCommand(callbackID, nil, error);
                    return;
                }

                NSString* base64String = [pdfDocumentData base64EncodedStringWithOptions:0];
                NSString* payload = [NSString stringWithFormat:@"data:application/pdf;base64,%@", base64String ?: @""];
                WailsAutomationCompleteCommand(callbackID, payload, nil);
            }];
            return;
        }

        WailsAutomationCompleteCommandMessage(callbackID, @"PDF export requires macOS 11.0 or newer");
    });
}

void webContentsViewGetCookies(void* view, uintptr_t callbackID) {
    WKWebView* webView = (WKWebView*)view;
    dispatch_async(dispatch_get_main_queue(), ^{
        [WailsAutomationCookieStore(webView) getAllCookies:^(NSArray<NSHTTPCookie*>* cookies) {
            NSMutableArray* result = [NSMutableArray arrayWithCapacity:cookies.count];
            for (NSHTTPCookie* cookie in cookies) {
                [result addObject:WailsAutomationCookieDictionary(cookie)];
            }
            WailsAutomationCompleteCommand(callbackID, result, nil);
        }];
    });
}

void webContentsViewSetCookie(void* view, const char* cookieJSON, uintptr_t callbackID) {
    WKWebView* webView = (WKWebView*)view;
    NSData* data = cookieJSON != NULL ? [[NSString stringWithUTF8String:cookieJSON] dataUsingEncoding:NSUTF8StringEncoding] : nil;
    NSError* jsonError = nil;
    NSDictionary* cookiePayload = data != nil ? [NSJSONSerialization JSONObjectWithData:data options:0 error:&jsonError] : nil;
    if (![cookiePayload isKindOfClass:[NSDictionary class]]) {
        WailsAutomationCompleteCommandMessage(callbackID, jsonError.localizedDescription ?: @"invalid cookie payload");
        return;
    }

    NSMutableDictionary* properties = [NSMutableDictionary dictionary];
    NSString* name = WailsAutomationString(cookiePayload[@"name"]);
    NSString* value = WailsAutomationString(cookiePayload[@"value"]);
    NSString* domain = WailsAutomationString(cookiePayload[@"domain"]);
    NSString* path = WailsAutomationString(cookiePayload[@"path"]);
    if (path.length == 0) {
        path = @"/";
    }

    if (name.length == 0 || domain.length == 0) {
        WailsAutomationCompleteCommandMessage(callbackID, @"cookie name and domain are required");
        return;
    }

    properties[NSHTTPCookieName] = name;
    properties[NSHTTPCookieValue] = value;
    properties[NSHTTPCookieDomain] = domain;
    properties[NSHTTPCookiePath] = path;

    id secure = cookiePayload[@"secure"];
    if (secure != nil) {
        properties[NSHTTPCookieSecure] = [secure boolValue] ? @"TRUE" : @"FALSE";
    }

    id expires = cookiePayload[@"expires"];
    if (expires != nil && expires != [NSNull null]) {
        NSTimeInterval interval = [expires doubleValue] / 1000.0;
        properties[NSHTTPCookieExpires] = [NSDate dateWithTimeIntervalSince1970:interval];
    }

    NSString* sameSite = WailsAutomationString(cookiePayload[@"sameSite"]);
    if (@available(macOS 10.15, *)) {
        if (sameSite.length > 0) {
            properties[NSHTTPCookieSameSitePolicy] = sameSite;
        }
    }

    NSHTTPCookie* cookie = [NSHTTPCookie cookieWithProperties:properties];
    if (cookie == nil) {
        WailsAutomationCompleteCommandMessage(callbackID, @"unable to create cookie from supplied properties");
        return;
    }

    dispatch_async(dispatch_get_main_queue(), ^{
        [WailsAutomationCookieStore(webView) setCookie:cookie completionHandler:^{
            WailsAutomationCompleteCommand(callbackID, @{ @"stored": @YES }, nil);
        }];
    });
}

void webContentsViewDeleteCookie(void* view, const char* name, const char* domain, const char* path, uintptr_t callbackID) {
    WKWebView* webView = (WKWebView*)view;
    NSString* cookieName = name != NULL ? [NSString stringWithUTF8String:name] : @"";
    NSString* cookieDomain = domain != NULL ? [NSString stringWithUTF8String:domain] : @"";
    NSString* cookiePath = path != NULL ? [NSString stringWithUTF8String:path] : @"";

    dispatch_async(dispatch_get_main_queue(), ^{
        [WailsAutomationCookieStore(webView) getAllCookies:^(NSArray<NSHTTPCookie*>* cookies) {
            NSMutableArray* matches = [NSMutableArray array];
            for (NSHTTPCookie* cookie in cookies) {
                if (cookieName.length > 0 && ![cookie.name isEqualToString:cookieName]) {
                    continue;
                }
                if (cookieDomain.length > 0 && ![cookie.domain isEqualToString:cookieDomain]) {
                    continue;
                }
                if (cookiePath.length > 0 && ![cookie.path isEqualToString:cookiePath]) {
                    continue;
                }
                [matches addObject:cookie];
            }

            if (matches.count == 0) {
                WailsAutomationCompleteCommand(callbackID, @{ @"deleted": @NO }, nil);
                return;
            }

            __block NSUInteger remaining = matches.count;
            for (NSHTTPCookie* cookie in matches) {
                [WailsAutomationCookieStore(webView) deleteCookie:cookie completionHandler:^{
                    remaining -= 1;
                    if (remaining == 0) {
                        WailsAutomationCompleteCommand(callbackID, @{ @"deleted": @YES }, nil);
                    }
                }];
            }
        }];
    });
}

void webContentsViewClearCookies(void* view, uintptr_t callbackID) {
    WKWebView* webView = (WKWebView*)view;
    dispatch_async(dispatch_get_main_queue(), ^{
        [WailsAutomationCookieStore(webView) getAllCookies:^(NSArray<NSHTTPCookie*>* cookies) {
            if (cookies.count == 0) {
                WailsAutomationCompleteCommand(callbackID, @{ @"cleared": @YES }, nil);
                return;
            }

            __block NSUInteger remaining = cookies.count;
            for (NSHTTPCookie* cookie in cookies) {
                [WailsAutomationCookieStore(webView) deleteCookie:cookie completionHandler:^{
                    remaining -= 1;
                    if (remaining == 0) {
                        WailsAutomationCompleteCommand(callbackID, @{ @"cleared": @YES }, nil);
                    }
                }];
            }
        }];
    });
}

bool webContentsViewAutomationSetInspectable(void* view, bool enabled) {
    __block bool success = false;
    WKWebView* webView = (WKWebView*)view;

    WailsRunOnMainThread(^{
        if (@available(macOS 13.3, *)) {
            webView.inspectable = enabled;
            success = true;
        }
    });

    return success;
}

bool webContentsViewAutomationSupportsInspection(void) {
    if (@available(macOS 13.3, *)) {
        return true;
    }
    return false;
}

bool webContentsViewAutomationSupportsProxyCapture(void) {
    if (@available(macOS 14.0, *)) {
        return true;
    }
    return false;
}

void windowAddWebContentsView(void* nsWindow, void* view) {
    NSWindow* window = (NSWindow*)nsWindow;
    WKWebView* webView = (WKWebView*)view;
    [window.contentView addSubview:webView];
    [webView setWantsLayer:YES];
    webView.layer.zPosition = 9999.0;
}

void windowRemoveWebContentsView(void* nsWindow, void* view) {
    WKWebView* webView = (WKWebView*)view;
    [webView removeFromSuperview];
}

void webContentsViewGoBack(void* view) {
    WKWebView* webView = (WKWebView*)view;
    dispatch_async(dispatch_get_main_queue(), ^{
        if ([webView canGoBack]) {
            [webView goBack];
        }
    });
}

const char* webContentsViewGetURL(void* view) {
    __block const char* result = NULL;
    WKWebView* webView = (WKWebView*)view;

    if ([NSThread isMainThread]) {
        if (webView.URL != nil) {
            result = strdup(webView.URL.absoluteString.UTF8String);
        }
    } else {
        dispatch_sync(dispatch_get_main_queue(), ^{
            if (webView.URL != nil) {
                result = strdup(webView.URL.absoluteString.UTF8String);
            }
        });
    }
    return result;
}

void webContentsViewTakeSnapshot(void* view, uintptr_t callbackID) {
    WKWebView* webView = (WKWebView*)view;
    dispatch_async(dispatch_get_main_queue(), ^{
        @try {
            if (@available(macOS 10.13, *)) {
                WKSnapshotConfiguration* config = [[WKSnapshotConfiguration alloc] init];
                [webView takeSnapshotWithConfiguration:config completionHandler:^(NSImage* image, NSError* error) {
                    if (error != nil || image == nil) {
                        browserViewSnapshotCallback(callbackID, NULL);
                        return;
                    }

                    @try {
                        CGImageRef cgRef = [image CGImageForProposedRect:NULL context:nil hints:nil];
                        if (cgRef == NULL) {
                            browserViewSnapshotCallback(callbackID, NULL);
                            return;
                        }

                        NSBitmapImageRep* newRep = [[NSBitmapImageRep alloc] initWithCGImage:cgRef];
                        if (newRep == nil) {
                            browserViewSnapshotCallback(callbackID, NULL);
                            return;
                        }

                        [newRep setSize:[image size]];
                        NSData* pngData = [newRep representationUsingType:NSBitmapImageFileTypePNG properties:@{}];

                        if (pngData == nil) {
                            browserViewSnapshotCallback(callbackID, NULL);
                            return;
                        }

                        NSString* base64String = [pngData base64EncodedStringWithOptions:0];
                        if (base64String == nil) {
                            browserViewSnapshotCallback(callbackID, NULL);
                            return;
                        }

                        NSString* fullBase64 = [NSString stringWithFormat:@"data:image/png;base64,%@", base64String];
                        browserViewSnapshotCallback(callbackID, [fullBase64 UTF8String]);
                    } @catch (NSException* innerException) {
                        NSLog(@"Error processing snapshot image: %@", innerException.reason);
                        browserViewSnapshotCallback(callbackID, NULL);
                    }
                }];
            } else {
                browserViewSnapshotCallback(callbackID, NULL);
            }
        } @catch (NSException* e) {
            NSLog(@"Exception in snapshot configuration: %@", e.reason);
            browserViewSnapshotCallback(callbackID, NULL);
        }
    });
}
