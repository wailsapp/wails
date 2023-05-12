//go:build darwin

#import "dialogs_darwin_delegate.h"

// Override shouldEnableURL
@implementation OpenPanelDelegate
- (BOOL)panel:(id)sender shouldEnableURL:(NSURL *)url {
    if (url == nil) {
        return NO;
    }
    NSFileManager *fileManager = [NSFileManager defaultManager];
    BOOL isDirectory = NO;
    if ([fileManager fileExistsAtPath:url.path isDirectory:&isDirectory] && isDirectory) {
        return YES;
    }
    if (self.allowedExtensions == nil) {
        return YES;
    }
    NSString *extension = url.pathExtension;
    if (extension == nil) {
        return NO;
    }
    if ([extension isEqualToString:@""]) {
        return NO;
    }
    if ([self.allowedExtensions containsObject:extension]) {
        return YES;
    }
    return NO;
}

@end




