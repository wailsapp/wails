//go:build darwin && !ios && !server

#import "dialogs_darwin_delegate.h"

// Override shouldEnableURL
@implementation OpenPanelDelegate
- (void)dealloc {
    [_allowedExtensions release];
    [_allowedTypeIdentifiers release];
    [super dealloc];
}

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

#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
    // Content type identifiers (AddContentType) are honoured alongside the
    // extension list: a file conforming to any of them is enabled.
    if (@available(macOS 11, *)) {
        if ([self.allowedTypeIdentifiers count] > 0) {
            UTType *fileType = nil;
            [url getResourceValue:&fileType forKey:NSURLContentTypeKey error:nil];
            if (fileType == nil && [[url pathExtension] length] > 0) {
                fileType = [UTType typeWithFilenameExtension:[url pathExtension]];
            }
            if (fileType != nil) {
                for (NSString *identifier in self.allowedTypeIdentifiers) {
                    UTType *allowedType = [UTType typeWithIdentifier:identifier];
                    if (allowedType != nil && [fileType conformsToType:allowedType]) {
                        return YES;
                    }
                }
            }
        }
    }
#endif

    NSString *filename = [[url lastPathComponent] lowercaseString];
    if (filename == nil || [filename isEqualToString:@""]) {
        return NO;
    }

    // Check if the extension is in our allowed list (case insensitive)
    for (NSString *allowedExt in self.allowedExtensions) {
        NSString *allowed = [allowedExt lowercaseString];
        if ([allowed length] == 0) {
            continue;
        }
        if ([filename hasSuffix:[@"." stringByAppendingString:allowed]]) {
            return YES;
        }
    }

    return NO;
}

@end
