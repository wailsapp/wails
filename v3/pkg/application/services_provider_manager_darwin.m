//go:build darwin && !ios && !server

#import <Cocoa/Cocoa.h>
#import <objc/runtime.h>

#include "services_provider_manager_darwin.h"

// servicesProviderHandle is the Go side. It receives the service name and
// the request JSON ({"text","files","data":{type:base64},"types"}) and
// returns a strdup'd response JSON ({"text","files","data","error"}) that
// the caller frees.
extern char* servicesProviderHandle(char *name, char *requestJSON);

static NSMutableDictionary<NSString *, NSArray<NSString *> *> *wailsRegisteredServices(void) {
    static NSMutableDictionary *registered = nil;
    static dispatch_once_t once;
    dispatch_once(&once, ^{
        registered = [[NSMutableDictionary alloc] init];
    });
    return registered;
}

static NSString *wailsServiceNameForSelector(SEL selector) {
    NSArray<NSString *> *parts = [NSStringFromSelector(selector) componentsSeparatedByString:@":"];
    // "<name>:userData:error:" splits into name, userData, error and "".
    if (parts.count != 4 || ![parts[1] isEqualToString:@"userData"] || ![parts[2] isEqualToString:@"error"]) {
        return nil;
    }
    return parts[0];
}

static NSData *wailsJSONData(id object) {
    if (object == nil || ![NSJSONSerialization isValidJSONObject:object]) {
        return nil;
    }
    return [NSJSONSerialization dataWithJSONObject:object options:0 error:nil];
}

// wailsServiceInvoke is the IMP shared by every registered service selector.
// It runs on the main thread while the requesting application waits.
static void wailsServiceInvoke(id self, SEL _cmd, NSPasteboard *pboard, NSString *userData, NSString **error) {
    @autoreleasepool {
        NSString *name = wailsServiceNameForSelector(_cmd);
        if (name == nil) {
            return;
        }
        NSArray<NSString *> *sendTypes = wailsRegisteredServices()[name] ?: @[];

        NSMutableDictionary *request = [NSMutableDictionary dictionary];
        request[@"text"] = [pboard stringForType:NSPasteboardTypeString] ?: @"";
        request[@"types"] = pboard.types ?: @[];

        NSMutableArray *files = [NSMutableArray array];
        NSArray *urls = [pboard readObjectsForClasses:@[[NSURL class]]
                                              options:@{NSPasteboardURLReadingFileURLsOnlyKey: @YES}];
        for (NSURL *url in urls) {
            if (url.path != nil) {
                [files addObject:url.path];
            }
        }
        request[@"files"] = files;

        NSMutableDictionary *data = [NSMutableDictionary dictionary];
        for (NSString *type in sendTypes) {
            NSData *bytes = [pboard dataForType:type];
            if (bytes != nil) {
                data[type] = [bytes base64EncodedStringWithOptions:0];
            }
        }
        request[@"data"] = data;

        NSData *requestJSON = wailsJSONData(request);
        if (requestJSON == nil) {
            if (error) {
                *error = @"Wails: cannot encode the pasteboard for the service handler";
            }
            return;
        }
        NSString *requestString = [[[NSString alloc] initWithData:requestJSON encoding:NSUTF8StringEncoding] autorelease];

        char *rawResponse = servicesProviderHandle((char *)[name UTF8String], (char *)[requestString UTF8String]);
        if (rawResponse == NULL) {
            return;
        }
        NSData *responseJSON = [NSData dataWithBytes:rawResponse length:strlen(rawResponse)];
        free(rawResponse);
        NSDictionary *response = [NSJSONSerialization JSONObjectWithData:responseJSON options:0 error:nil];
        if (![response isKindOfClass:[NSDictionary class]]) {
            return;
        }

        NSString *message = response[@"error"];
        if ([message isKindOfClass:[NSString class]] && message.length > 0) {
            if (error) {
                *error = message;
            }
            return;
        }

        NSString *text = [response[@"text"] isKindOfClass:[NSString class]] ? response[@"text"] : @"";
        NSArray *paths = [response[@"files"] isKindOfClass:[NSArray class]] ? response[@"files"] : @[];
        NSDictionary *payload = [response[@"data"] isKindOfClass:[NSDictionary class]] ? response[@"data"] : @{};
        if (text.length == 0 && paths.count == 0 && payload.count == 0) {
            return;
        }

        [pboard clearContents];
        if (paths.count > 0) {
            NSMutableArray *fileURLs = [NSMutableArray array];
            for (NSString *path in paths) {
                if ([path isKindOfClass:[NSString class]]) {
                    [fileURLs addObject:[NSURL fileURLWithPath:path]];
                }
            }
            [pboard writeObjects:fileURLs];
        }
        if (text.length > 0) {
            [pboard setString:text forType:NSPasteboardTypeString];
        }
        for (NSString *type in payload) {
            NSString *encoded = payload[type];
            if (![encoded isKindOfClass:[NSString class]]) {
                continue;
            }
            NSData *bytes = [[[NSData alloc] initWithBase64EncodedString:encoded options:0] autorelease];
            if (bytes != nil) {
                [pboard setData:bytes forType:type];
            }
        }
    }
}

// WailsServicesProvider is installed as NSApplication.servicesProvider. It
// has no methods of its own: the Objective-C runtime asks
// resolveInstanceMethod: for `<name>:userData:error:` the first time AppKit
// sends it, and the method is added for every registered service name.
@interface WailsServicesProvider : NSObject
+ (instancetype)shared;
@end

@implementation WailsServicesProvider

+ (instancetype)shared {
    static WailsServicesProvider *provider = nil;
    static dispatch_once_t once;
    dispatch_once(&once, ^{
        provider = [[WailsServicesProvider alloc] init];
    });
    return provider;
}

+ (BOOL)resolveInstanceMethod:(SEL)selector {
    NSString *name = wailsServiceNameForSelector(selector);
    if (name != nil && wailsRegisteredServices()[name] != nil) {
        class_addMethod(self, selector, (IMP)wailsServiceInvoke, "v@:@@^@");
        return YES;
    }
    return [super resolveInstanceMethod:selector];
}

@end

void wailsServicesProviderRegister(const char *name, const char *sendTypesJSON) {
    @autoreleasepool {
        NSString *serviceName = [NSString stringWithUTF8String:name];
        NSArray *sendTypes = @[];
        if (sendTypesJSON != NULL) {
            NSData *data = [NSData dataWithBytes:sendTypesJSON length:strlen(sendTypesJSON)];
            id parsed = [NSJSONSerialization JSONObjectWithData:data options:0 error:nil];
            if ([parsed isKindOfClass:[NSArray class]]) {
                sendTypes = parsed;
            }
        }
        wailsRegisteredServices()[serviceName] = sendTypes;
        NSApplication *app = [NSApplication sharedApplication];
        if (app.servicesProvider != [WailsServicesProvider shared]) {
            app.servicesProvider = [WailsServicesProvider shared];
        }
        NSUpdateDynamicServices();
    }
}
