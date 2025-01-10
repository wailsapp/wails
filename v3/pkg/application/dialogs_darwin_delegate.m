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

    // If no extensions specified, allow all files
    if (self.allowedExtensions == nil || [self.allowedExtensions count] == 0) {
        return YES;
    }

    NSString *extension = [url.pathExtension lowercaseString];
    if (extension == nil || [extension isEqualToString:@""]) {
        return NO;
    }

    // Check if the extension is in our allowed list (case insensitive)
    for (NSString *allowedExt in self.allowedExtensions) {
        if ([[allowedExt lowercaseString] isEqualToString:extension]) {
            return YES;
        }
    }

    return NO;
}

@end
