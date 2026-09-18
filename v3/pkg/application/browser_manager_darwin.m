//go:build darwin && !ios && !server

#import <Cocoa/Cocoa.h>

#include "browser_manager_darwin.h"

// wailsCompletionCallback is the Go side of asynchronous NSWorkspace calls:
// message is NULL on success, otherwise an error description.
extern void wailsCompletionCallback(unsigned long long id, char *message);

static void wailsCompleteWithError(unsigned long long id, NSError *error) {
    if (error == nil) {
        wailsCompletionCallback(id, NULL);
        return;
    }
    wailsCompletionCallback(id, (char *)[error.localizedDescription UTF8String]);
}

static void wailsCompleteWithMessage(unsigned long long id, NSString *message) {
    wailsCompletionCallback(id, (char *)[message UTF8String]);
}

// wailsApplicationInfo describes an application bundle as
// {"name","bundleID","path"}. It is shared with environment_manager_darwin.m.
NSDictionary *wailsApplicationInfo(NSURL *appURL) {
    if (appURL == nil) {
        return nil;
    }
    NSBundle *bundle = [NSBundle bundleWithURL:appURL];
    NSString *name = [bundle objectForInfoDictionaryKey:@"CFBundleDisplayName"];
    if (![name isKindOfClass:[NSString class]] || name.length == 0) {
        name = [bundle objectForInfoDictionaryKey:@"CFBundleName"];
    }
    if (![name isKindOfClass:[NSString class]] || name.length == 0) {
        name = [[NSFileManager defaultManager] displayNameAtPath:appURL.path];
    }
    if (![name isKindOfClass:[NSString class]]) {
        name = [appURL.lastPathComponent stringByDeletingPathExtension] ?: @"";
    }
    return @{
        @"name": name ?: @"",
        @"bundleID": bundle.bundleIdentifier ?: @"",
        @"path": appURL.path ?: @"",
    };
}

// wailsApplicationURL resolves a bundle identifier or a .app path to the
// bundle URL, or nil when nothing matches.
NSURL *wailsApplicationURL(NSString *app) {
    if (app.length == 0) {
        return nil;
    }
    if ([app hasPrefix:@"/"] || [app hasPrefix:@"~"] || [app.pathExtension isEqualToString:@"app"]) {
        NSString *path = [app stringByExpandingTildeInPath];
        BOOL isDirectory = NO;
        if ([[NSFileManager defaultManager] fileExistsAtPath:path isDirectory:&isDirectory] && isDirectory) {
            return [NSURL fileURLWithPath:path];
        }
        return nil;
    }
    NSWorkspace *workspace = [NSWorkspace sharedWorkspace];
    if (@available(macOS 10.15, *)) {
        return [workspace URLForApplicationWithBundleIdentifier:app];
    }
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
    NSString *path = [workspace absolutePathForAppBundleWithIdentifier:app];
#pragma clang diagnostic pop
    return path != nil ? [NSURL fileURLWithPath:path] : nil;
}

static char *wailsJSONCStringFromObject(id object) {
    NSData *data = [NSJSONSerialization dataWithJSONObject:object options:0 error:nil];
    if (data == nil) {
        return strdup("[]");
    }
    NSString *json = [[[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding] autorelease];
    return strdup([json UTF8String]);
}

void wailsWorkspaceOpenWith(unsigned long long id, const char *path, const char *app) {
    @autoreleasepool {
        NSString *appName = [NSString stringWithUTF8String:app];
        NSURL *appURL = wailsApplicationURL(appName);
        if (appURL == nil) {
            wailsCompleteWithMessage(id, [NSString stringWithFormat:@"application not found: %@", appName]);
            return;
        }
        NSURL *fileURL = [NSURL fileURLWithPath:[NSString stringWithUTF8String:path]];
        NSWorkspace *workspace = [NSWorkspace sharedWorkspace];
        if (@available(macOS 10.15, *)) {
            NSWorkspaceOpenConfiguration *configuration = [NSWorkspaceOpenConfiguration configuration];
            [workspace openURLs:@[fileURL]
             withApplicationAtURL:appURL
                    configuration:configuration
                completionHandler:^(NSRunningApplication *running, NSError *error) {
                    wailsCompleteWithError(id, error);
                }];
            return;
        }
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
        BOOL ok = [workspace openFile:fileURL.path withApplication:appURL.path];
#pragma clang diagnostic pop
        if (ok) {
            wailsCompletionCallback(id, NULL);
        } else {
            wailsCompleteWithMessage(id, @"the application declined to open the file");
        }
    }
}

char *wailsWorkspaceApplicationsForFile(const char *path) {
    @autoreleasepool {
        NSURL *fileURL = [NSURL fileURLWithPath:[NSString stringWithUTF8String:path]];
        NSArray<NSURL *> *urls = nil;
        if (@available(macOS 12.0, *)) {
            urls = [[NSWorkspace sharedWorkspace] URLsForApplicationsToOpenURL:fileURL];
        } else {
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
            CFArrayRef list = LSCopyApplicationURLsForURL((CFURLRef)fileURL, kLSRolesAll);
#pragma clang diagnostic pop
            if (list != NULL) {
                urls = [(NSArray *)list autorelease];
            }
        }
        NSMutableArray *result = [NSMutableArray array];
        NSMutableSet *seen = [NSMutableSet set];
        for (NSURL *url in urls) {
            NSDictionary *info = wailsApplicationInfo(url);
            if (info == nil || [seen containsObject:info[@"path"]]) {
                continue;
            }
            [seen addObject:info[@"path"]];
            [result addObject:info];
        }
        return wailsJSONCStringFromObject(result);
    }
}

char *wailsWorkspaceActivateApplication(const char *bundleID) {
    @autoreleasepool {
        NSString *identifier = [NSString stringWithUTF8String:bundleID];
        NSArray<NSRunningApplication *> *running = [NSRunningApplication runningApplicationsWithBundleIdentifier:identifier];
        if (running.count == 0) {
            return strdup("not running");
        }
        BOOL activated = NO;
        for (NSRunningApplication *app in running) {
            // macOS 14 makes activation cooperative: the active application
            // hands focus over with activateFromApplication:options:, and
            // the older activateWithOptions: is honoured only in limited
            // cases. Use the cooperative form where available.
            if (@available(macOS 14.0, *)) {
                if ([app activateFromApplication:[NSRunningApplication currentApplication]
                                          options:NSApplicationActivateAllWindows]) {
                    activated = YES;
                    continue;
                }
            }
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
            if ([app activateWithOptions:NSApplicationActivateAllWindows]) {
                activated = YES;
            }
#pragma clang diagnostic pop
        }
        if (!activated) {
            return strdup("the application refused activation");
        }
        return NULL;
    }
}
