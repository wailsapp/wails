//go:build darwin && !ios

package dock

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

void hideDockIcon() {
    dispatch_sync(dispatch_get_main_queue(), ^{
        [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
    });
}

void showDockIcon() {
    dispatch_sync(dispatch_get_main_queue(), ^{
        [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
    });
}

bool setBadge(const char *label) {
    __block bool success = false;
    dispatch_sync(dispatch_get_main_queue(), ^{
        // Ensure the app is in Regular activation policy (dock icon visible)
        NSApplicationActivationPolicy currentPolicy = [NSApp activationPolicy];
        if (currentPolicy != NSApplicationActivationPolicyRegular) {
            success = false;
            return;
        }

        NSString *nsLabel = nil;
		if (label != NULL) {
			nsLabel = [NSString stringWithUTF8String:label];
		}
		[[NSApp dockTile] setBadgeLabel:nsLabel];
		[[NSApp dockTile] display];
		success = true;
    });
    return success;
}

// WailsDockProgressView draws a rounded progress bar. It sits above an
// NSImageView showing the app icon inside the dock tile's contentView, so the
// icon keeps drawing and the Dock still overlays the badge label.
@interface WailsDockProgressView : NSView
@property double fraction;
@end

@implementation WailsDockProgressView
- (void)drawRect:(NSRect)dirtyRect {
    NSRect bounds = self.bounds;
    CGFloat radius = bounds.size.height / 2.0;
    NSBezierPath *track = [NSBezierPath bezierPathWithRoundedRect:bounds xRadius:radius yRadius:radius];
    [[NSColor colorWithCalibratedWhite:1.0 alpha:0.85] setFill];
    [track fill];
    [[NSColor colorWithCalibratedWhite:0.0 alpha:0.25] setStroke];
    track.lineWidth = 1.0;
    [track stroke];

    NSRect fillRect = NSInsetRect(bounds, 2.0, 2.0);
    fillRect.size.width = MAX(0.0, fillRect.size.width * MIN(1.0, MAX(0.0, self.fraction)));
    if (fillRect.size.width > 0.0) {
        CGFloat fillRadius = fillRect.size.height / 2.0;
        NSBezierPath *fill = [NSBezierPath bezierPathWithRoundedRect:fillRect xRadius:fillRadius yRadius:fillRadius];
        NSColor *tint = nil;
        if (@available(macOS 10.14, *)) {
            tint = [NSColor controlAccentColor];
        } else {
            tint = [NSColor systemBlueColor];
        }
        [tint setFill];
        [fill fill];
    }
}
@end

static NSView *dockProgressContainer = nil;
static NSImageView *dockProgressIcon = nil;
static WailsDockProgressView *dockProgressBar = nil;

// setDockProgress installs (once) an icon + bar content view on the dock
// tile and updates the fraction. The Dock draws contentView instead of the
// app icon, so the NSImageView keeps the icon visible underneath.
//
// dockOnMain runs block synchronously on the main thread. dispatch_sync from
// the main thread deadlocks, so short-circuit when already there (the
// service may be driven from a menu callback or an InvokeSync block).
static void dockOnMain(dispatch_block_t block) {
    if ([NSThread isMainThread]) {
        block();
    } else {
        dispatch_sync(dispatch_get_main_queue(), block);
    }
}

void setDockProgress(double fraction) {
    dockOnMain(^{
        NSDockTile *tile = [NSApp dockTile];
        NSSize size = [tile size];
        if (size.width <= 0 || size.height <= 0) {
            size = NSMakeSize(128, 128);
        }
        NSRect frame = NSMakeRect(0, 0, size.width, size.height);
        if (dockProgressContainer == nil) {
            dockProgressContainer = [[NSView alloc] initWithFrame:frame];
            dockProgressIcon = [[NSImageView alloc] initWithFrame:frame];
            dockProgressIcon.imageScaling = NSImageScaleProportionallyUpOrDown;
            dockProgressIcon.autoresizingMask = NSViewWidthSizable | NSViewHeightSizable;
            [dockProgressContainer addSubview:dockProgressIcon];
            dockProgressBar = [[WailsDockProgressView alloc] initWithFrame:NSZeroRect];
            dockProgressBar.autoresizingMask = NSViewWidthSizable | NSViewMinYMargin;
            [dockProgressContainer addSubview:dockProgressBar];
        }
        dockProgressContainer.frame = frame;
        dockProgressIcon.frame = frame;
        dockProgressIcon.image = [NSApp applicationIconImage];
        CGFloat barHeight = MAX(8.0, size.height * 0.11);
        dockProgressBar.frame = NSMakeRect(size.width * 0.1, size.height * 0.08, size.width * 0.8, barHeight);
        dockProgressBar.fraction = fraction;
        [dockProgressBar setNeedsDisplay:YES];
        if (tile.contentView != dockProgressContainer) {
            tile.contentView = dockProgressContainer;
        }
        [tile display];
    });
}

// clearDockProgress removes the content view so the Dock draws the plain
// app icon again. The badge label is untouched.
void clearDockProgress(void) {
    dockOnMain(^{
        NSDockTile *tile = [NSApp dockTile];
        if (tile.contentView == dockProgressContainer) {
            tile.contentView = nil;
        }
        [tile display];
    });
}
*/
import "C"
import (
	"context"
	"fmt"
	"sync"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type darwinDock struct {
	mu       sync.RWMutex
	Badge    *string
	progress *float64
}

// Creates a new Dock Service.
func New() *DockService {
	return &DockService{
		impl: &darwinDock{
			Badge: nil,
		},
	}
}

// NewWithOptions creates a new dock service with badge options.
// Currently, options are not available on macOS and are ignored.
func NewWithOptions(options BadgeOptions) *DockService {
	return New()
}

func (d *darwinDock) Startup(ctx context.Context, options application.ServiceOptions) error {
	return nil
}

func (d *darwinDock) Shutdown() error {
	return nil
}

// HideAppIcon hides the app icon in the macOS Dock.
func (d *darwinDock) HideAppIcon() {
	C.hideDockIcon()
}

// ShowAppIcon shows the app icon in the macOS Dock.
// Note: After showing the dock icon, you may need to call SetBadge again
// to reapply any previously set badge, as changing activation policies clears the badge.
func (d *darwinDock) ShowAppIcon() {
	C.showDockIcon()
}

// setBadge handles the C call and updates the internal badge state with locking.
func (d *darwinDock) setBadge(label *string) error {
	var cLabel *C.char
	if label != nil {
		cLabel = C.CString(*label)
		defer C.free(unsafe.Pointer(cLabel))
	}

	success := C.setBadge(cLabel)
	if !success {
		return fmt.Errorf("failed to set badge")
	}

	d.mu.Lock()
	d.Badge = label
	d.mu.Unlock()

	return nil
}

// SetBadge sets the badge label on the application icon.
// Available default badge labels:
// Single space " " empty badge
// Empty string "" dot "●" indeterminate badge
func (d *darwinDock) SetBadge(label string) error {
	// Always pick a label (use "●" if empty), then allocate + free exactly once.
	if label == "" {
		label = "●" // Default badge character
	}
	return d.setBadge(&label)
}

// SetCustomBadge is not supported on macOS, SetBadge is called instead.
func (d *darwinDock) SetCustomBadge(label string, options BadgeOptions) error {
	return d.SetBadge(label)
}

// RemoveBadge removes the badge label from the application icon.
func (d *darwinDock) RemoveBadge() error {
	return d.setBadge(nil)
}

// GetBadge returns the badge label on the application icon.
func (d *darwinDock) GetBadge() *string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.Badge
}

// SetProgress draws a progress bar over the Dock icon. The tile content
// view keeps the app icon underneath and the badge label still overlays it.
func (d *darwinDock) SetProgress(fraction float64) error {
	C.setDockProgress(C.double(fraction))
	d.mu.Lock()
	d.progress = &fraction
	d.mu.Unlock()
	return nil
}

// ClearProgress removes the progress bar from the Dock icon.
func (d *darwinDock) ClearProgress() error {
	C.clearDockProgress()
	d.mu.Lock()
	d.progress = nil
	d.mu.Unlock()
	return nil
}

// GetProgress returns the current progress fraction, or nil when cleared.
func (d *darwinDock) GetProgress() *float64 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.progress
}
