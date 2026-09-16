//go:build darwin && !ios && !server

#import "quicklook_manager_darwin.h"
#import <Quartz/Quartz.h>
#import <QuickLookThumbnailing/QuickLookThumbnailing.h>

// Delivers a finished thumbnail (PNG bytes) or an error message to Go. The
// callee copies both before returning.
extern void quickLookThumbnailResult(unsigned long long requestID, const void *png, int length, const char *errorMessage);

// WailsQuickLookController is the standalone data source and delegate for
// the shared QLPreviewPanel. Nothing in the responder chain claims the panel,
// so it is assigned directly.
@interface WailsQuickLookController : NSObject <QLPreviewPanelDataSource, QLPreviewPanelDelegate>
@property (nonatomic, strong) NSArray<NSURL *> *items;
+ (instancetype)shared;
@end

@implementation WailsQuickLookController

+ (instancetype)shared {
    static WailsQuickLookController *shared = nil;
    static dispatch_once_t once;
    dispatch_once(&once, ^{
        shared = [[WailsQuickLookController alloc] init];
        shared.items = @[];
    });
    return shared;
}

- (NSInteger)numberOfPreviewItemsInPreviewPanel:(QLPreviewPanel *)panel {
    return (NSInteger)self.items.count;
}

- (id<QLPreviewItem>)previewPanel:(QLPreviewPanel *)panel previewItemAtIndex:(NSInteger)index {
    if (index < 0 || (NSUInteger)index >= self.items.count) {
        return nil;
    }
    return self.items[(NSUInteger)index];
}

@end

bool quickLookPreview(const char *pathsJSON) {
  @autoreleasepool {
    NSData *data = [NSData dataWithBytes:pathsJSON length:strlen(pathsJSON)];
    NSArray *paths = [NSJSONSerialization JSONObjectWithData:data options:0 error:nil];
    if (![paths isKindOfClass:[NSArray class]] || paths.count == 0) {
        return false;
    }
    NSMutableArray<NSURL *> *urls = [NSMutableArray arrayWithCapacity:paths.count];
    for (id path in paths) {
        if ([path isKindOfClass:[NSString class]]) {
            [urls addObject:[NSURL fileURLWithPath:path]];
        }
    }
    if (urls.count == 0) {
        return false;
    }
    WailsQuickLookController *controller = [WailsQuickLookController shared];
    controller.items = urls;

    QLPreviewPanel *panel = [QLPreviewPanel sharedPreviewPanel];
    if (panel == nil) {
        return false;
    }
    panel.dataSource = controller;
    panel.delegate = controller;
    [panel reloadData];
    panel.currentPreviewItemIndex = 0;
    if (![panel isVisible]) {
        [panel makeKeyAndOrderFront:nil];
    }
    return true;
  }
}

void quickLookClosePreview(void) {
    if (![QLPreviewPanel sharedPreviewPanelExists]) {
        return;
    }
    QLPreviewPanel *panel = [QLPreviewPanel sharedPreviewPanel];
    if ([panel isVisible]) {
        // QLPreviewPanel ignores orderOut: and setIsVisible: while no object
        // in the responder chain controls it (the panel is driven through
        // its data source here), but the underlying NSWindow ordering call
        // hides it reliably.
        [panel orderWindow:NSWindowOut relativeTo:0];
    }
}

bool quickLookIsPreviewOpen(void) {
    return [QLPreviewPanel sharedPreviewPanelExists] && [[QLPreviewPanel sharedPreviewPanel] isVisible];
}

bool quickLookThumbnailAvailable(void) {
    if (@available(macOS 10.15, *)) {
        return true;
    }
    return false;
}

static NSData *wailsQuickLookPNG(CGImageRef image) {
    if (image == NULL) {
        return nil;
    }
    NSBitmapImageRep *rep = [[[NSBitmapImageRep alloc] initWithCGImage:image] autorelease];
    return [rep representationUsingType:NSBitmapImageFileTypePNG properties:@{}];
}

void quickLookThumbnail(unsigned long long requestID, const char *path, int width, int height, double scale, bool iconMode) {
  @autoreleasepool {
    if (@available(macOS 10.15, *)) {
        NSURL *url = [NSURL fileURLWithPath:[NSString stringWithUTF8String:path]];
        QLThumbnailGenerationRequest *request =
            [[[QLThumbnailGenerationRequest alloc] initWithFileAtURL:url
                                                                 size:CGSizeMake(width, height)
                                                                scale:(CGFloat)scale
                                                  representationTypes:QLThumbnailGenerationRequestRepresentationTypeAll] autorelease];
        request.iconMode = iconMode;
        [[QLThumbnailGenerator sharedGenerator]
            generateBestRepresentationForRequest:request
                               completionHandler:^(QLThumbnailRepresentation *thumbnail, NSError *error) {
            // Runs on a Quick Look queue, never the main thread.
            @autoreleasepool {
                NSData *png = wailsQuickLookPNG(thumbnail.CGImage);
                if (png == nil || png.length == 0) {
                    NSString *message = error.localizedDescription ?: @"quick look: no thumbnail available";
                    quickLookThumbnailResult(requestID, NULL, 0, [message UTF8String]);
                    return;
                }
                quickLookThumbnailResult(requestID, png.bytes, (int)png.length, NULL);
            }
        }];
        return;
    }
    quickLookThumbnailResult(requestID, NULL, 0, "quick look: thumbnails need macOS 10.15 or later");
  }
}
