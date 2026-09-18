//go:build darwin && !ios && !server

#import "webview_window_accessory_darwin.h"

// Availability guards use __MAC_OS_X_VERSION_MAX_ALLOWED: the legacy
// MAC_OS_X_VERSION_MAX_ALLOWED macro is capped at macOS 14 by
// AvailabilityMacros.h, so a guard on it can never see macOS 26 API.
enum {
    WailsMacScrollEdgeEffectStyleAutomatic = 0,
    WailsMacScrollEdgeEffectStyleSoft = 1,
    WailsMacScrollEdgeEffectStyleHard = 2,
};

static bool isSupportedAccessoryViewController(id controller) {
    if (controller == nil) return false;
    if ([controller isKindOfClass:[NSTitlebarAccessoryViewController class]]) return true;
    Class splitAccessoryClass = NSClassFromString(@"NSSplitViewItemAccessoryViewController");
    return splitAccessoryClass != Nil && [controller isKindOfClass:splitAccessoryClass];
}

static int accessoryViewControllerKind(id controller) {
    if (controller == nil) return 0;
    if ([controller isKindOfClass:[NSTitlebarAccessoryViewController class]]) return 1;
    Class splitAccessoryClass = NSClassFromString(@"NSSplitViewItemAccessoryViewController");
    if (splitAccessoryClass != Nil && [controller isKindOfClass:splitAccessoryClass]) return 2;
    return 0;
}

int macAccessoryViewControllerKind(void* pointer) {
    if ([NSThread isMainThread]) return accessoryViewControllerKind((id)pointer);
    __block int result = 0;
    dispatch_sync(dispatch_get_main_queue(), ^{
        result = accessoryViewControllerKind((id)pointer);
    });
    return result;
}

static bool accessorySupportsScrollEdgeEffectStyle(id controller) {
    if (!isSupportedAccessoryViewController(controller)) return false;
#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 260100
    if (@available(macOS 26.1, *)) return true;
#endif
    return false;
}

bool macAccessoryViewControllerSupportsScrollEdgeEffectStyle(void* pointer) {
    if ([NSThread isMainThread]) return accessorySupportsScrollEdgeEffectStyle((id)pointer);
    __block bool result = false;
    dispatch_sync(dispatch_get_main_queue(), ^{
        result = accessorySupportsScrollEdgeEffectStyle((id)pointer);
    });
    return result;
}

static int accessorySetScrollEdgeEffectStyle(id controller, int style) {
    if (!isSupportedAccessoryViewController(controller)) {
        return WailsMacAccessoryStyleInvalidController;
    }
    if (style < WailsMacScrollEdgeEffectStyleAutomatic ||
        style > WailsMacScrollEdgeEffectStyleHard) {
        return WailsMacAccessoryStyleInvalidController;
    }
#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 260100
    if (@available(macOS 26.1, *)) {
        NSScrollEdgeEffectStyle* nativeStyle = NSScrollEdgeEffectStyle.automaticStyle;
        if (style == WailsMacScrollEdgeEffectStyleSoft) {
            nativeStyle = NSScrollEdgeEffectStyle.softStyle;
        } else if (style == WailsMacScrollEdgeEffectStyleHard) {
            nativeStyle = NSScrollEdgeEffectStyle.hardStyle;
        }
        [(NSTitlebarAccessoryViewController*)controller setPreferredScrollEdgeEffectStyle:nativeStyle];
        return WailsMacAccessoryStyleApplied;
    }
#endif
    return style == WailsMacScrollEdgeEffectStyleAutomatic
        ? WailsMacAccessoryStyleApplied
        : WailsMacAccessoryStyleUnavailable;
}

int macAccessoryViewControllerSetScrollEdgeEffectStyle(void* pointer, int style) {
    if ([NSThread isMainThread]) return accessorySetScrollEdgeEffectStyle((id)pointer, style);
    __block int result = WailsMacAccessoryStyleInvalidController;
    dispatch_sync(dispatch_get_main_queue(), ^{
        result = accessorySetScrollEdgeEffectStyle((id)pointer, style);
    });
    return result;
}

static int accessoryScrollEdgeEffectStyle(id controller) {
    if (!isSupportedAccessoryViewController(controller)) {
        return WailsMacAccessoryStyleInvalidController;
    }
#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 260100
    if (@available(macOS 26.1, *)) {
        NSScrollEdgeEffectStyle* style =
            [(NSTitlebarAccessoryViewController*)controller preferredScrollEdgeEffectStyle];
        if ([style isEqual:NSScrollEdgeEffectStyle.softStyle]) {
            return WailsMacScrollEdgeEffectStyleSoft;
        }
        if ([style isEqual:NSScrollEdgeEffectStyle.hardStyle]) {
            return WailsMacScrollEdgeEffectStyleHard;
        }
        return WailsMacScrollEdgeEffectStyleAutomatic;
    }
#endif
    return WailsMacAccessoryStyleUnavailable;
}

int macAccessoryViewControllerScrollEdgeEffectStyle(void* pointer) {
    if ([NSThread isMainThread]) return accessoryScrollEdgeEffectStyle((id)pointer);
    __block int result = WailsMacAccessoryStyleInvalidController;
    dispatch_sync(dispatch_get_main_queue(), ^{
        result = accessoryScrollEdgeEffectStyle((id)pointer);
    });
    return result;
}

// ---------------------------------------------------------------------------
// Native control strips (MacAccessory)
// ---------------------------------------------------------------------------

#import <objc/runtime.h>

static const void* WailsAccessoryHostAssociationKey = &WailsAccessoryHostAssociationKey;
static const CGFloat WailsAccessoryHorizontalInset = 8.0;
static const CGFloat WailsAccessoryDefaultSearchWidth = 180.0;

// WailsAccessoryHost owns the control strip inside one accessory view
// controller: the container view, the NSStackView, the control lookup by
// identifier, and the target of every control action. It is associated
// with the controller so the controller's lifetime governs it.
@interface WailsAccessoryHost : NSObject
@property int kind;
@property int layout;
@property double height;
@property (retain) NSView* container;
@property (retain) NSStackView* stack;
@property (retain) NSLayoutConstraint* heightConstraint;
@property (retain) NSMutableDictionary<NSNumber*, NSView*>* controls;
@property (retain) NSMutableDictionary<NSNumber*, NSLayoutConstraint*>* widths;
@property (retain) NSMutableDictionary<NSNumber*, NSMenu*>* menus;
@end

@implementation WailsAccessoryHost
- (void)dealloc {
    [_container release];
    [_stack release];
    [_heightConstraint release];
    [_controls release];
    [_widths release];
    [_menus release];
    [super dealloc];
}

- (NSControlSize)controlSize {
    return self.kind == WailsMacAccessoryKindTitlebar ? NSControlSizeSmall : NSControlSizeRegular;
}

- (NSFont*)controlFont {
    return [NSFont systemFontOfSize:[NSFont systemFontSizeForControlSize:[self controlSize]]];
}

// relayout keeps the container frame in step with the strip. Leading and
// trailing titlebar accessories are sized to fit their controls; every other
// placement receives its width from AppKit and only needs the height.
- (void)relayout {
    if (self.heightConstraint != nil) self.heightConstraint.constant = self.height;
    [self.container layoutSubtreeIfNeeded];
    NSRect frame = self.container.frame;
    frame.size.height = self.height;
    if (self.kind == WailsMacAccessoryKindTitlebar &&
        (self.layout == WailsMacAccessoryLayoutLeading || self.layout == WailsMacAccessoryLayoutTrailing)) {
        NSSize fit = self.stack.fittingSize;
        frame.size.width = MAX(fit.width, 0) + 2 * WailsAccessoryHorizontalInset;
    }
    self.container.frame = frame;
}

- (NSView*)controlForID:(unsigned long long)controlID {
    return self.controls[@(controlID)];
}

- (void)setWidth:(double)width forControlID:(unsigned long long)controlID {
    NSView* view = [self controlForID:controlID];
    if (view == nil) return;
    NSNumber* key = @(controlID);
    NSLayoutConstraint* existing = self.widths[key];
    if (existing != nil) {
        existing.active = NO;
        [self.widths removeObjectForKey:key];
    }
    if (width <= 0 && [view isKindOfClass:[NSSearchField class]]) {
        width = WailsAccessoryDefaultSearchWidth;
    }
    if (width > 0) {
        NSLayoutConstraint* constraint = [view.widthAnchor constraintEqualToConstant:width];
        constraint.priority = NSLayoutPriorityRequired - 1;
        constraint.active = YES;
        self.widths[key] = constraint;
    }
    [self relayout];
}

- (void)addControlView:(NSView*)view forControlID:(unsigned long long)controlID
    tooltip:(const char*)tooltip disabled:(bool)disabled hidden:(bool)hidden width:(double)width {
    view.translatesAutoresizingMaskIntoConstraints = NO;
    if ([view isKindOfClass:[NSControl class]]) {
        NSControl* control = (NSControl*)view;
        control.tag = (NSInteger)controlID;
        control.enabled = !disabled;
    }
    view.toolTip = tooltip != NULL && strlen(tooltip) > 0 ? [NSString stringWithUTF8String:tooltip] : nil;
    view.hidden = hidden;
    self.controls[@(controlID)] = view;
    [self.stack addArrangedSubview:view];
    [self setWidth:width forControlID:controlID];
}

- (void)handleButton:(id)sender {
    if (![sender isKindOfClass:[NSControl class]]) return;
    processMacAccessoryClicked((unsigned long long)((NSControl*)sender).tag);
}

- (void)handleMenuButton:(id)sender {
    if (![sender isKindOfClass:[NSButton class]]) return;
    NSButton* button = (NSButton*)sender;
    NSMenu* menu = self.menus[@(button.tag)];
    if (menu == nil) return;
    [menu popUpMenuPositioningItem:nil atLocation:NSMakePoint(0, -4) inView:button];
}

- (void)handleSearch:(id)sender {
    if (![sender isKindOfClass:[NSSearchField class]]) return;
    NSSearchField* field = (NSSearchField*)sender;
    processMacAccessorySearched((unsigned long long)field.tag, (char*)field.stringValue.UTF8String);
}

- (void)handleSegment:(id)sender {
    if (![sender isKindOfClass:[NSSegmentedControl class]]) return;
    NSSegmentedControl* control = (NSSegmentedControl*)sender;
    processMacAccessorySegmentChanged((unsigned long long)control.tag, (int)control.selectedSegment);
}
@end

static WailsAccessoryHost* accessoryHost(void* controller) {
    if (controller == NULL) return nil;
    return objc_getAssociatedObject((id)controller, WailsAccessoryHostAssociationKey);
}

static NSImage* accessorySymbolImage(const char* symbolName, NSString* description) {
    if (symbolName == NULL || strlen(symbolName) == 0) return nil;
    if (@available(macOS 11.0, *)) {
        return [NSImage imageWithSystemSymbolName:[NSString stringWithUTF8String:symbolName]
                         accessibilityDescription:description];
    }
    return nil;
}

static NSArray<NSString*>* accessoryStringsFromJSON(const char* json) {
    if (json == NULL || strlen(json) == 0) return @[];
    NSData* data = [NSData dataWithBytes:json length:strlen(json)];
    id parsed = [NSJSONSerialization JSONObjectWithData:data options:0 error:nil];
    if (![parsed isKindOfClass:[NSArray class]]) return @[];
    NSMutableArray<NSString*>* strings = [NSMutableArray array];
    for (id value in (NSArray*)parsed) {
        [strings addObject:[value isKindOfClass:[NSString class]] ? value : @""];
    }
    return strings;
}

static NSString* accessoryString(const char* value) {
    return value == NULL ? @"" : [NSString stringWithUTF8String:value];
}

bool macAccessorySplitItemAccessoriesSupported(void) {
#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 260000
    if (@available(macOS 26.0, *)) return true;
#endif
    return false;
}

static NSLayoutAttribute accessoryTitlebarAttribute(int layout) {
    switch (layout) {
    case WailsMacAccessoryLayoutLeading:
        if (@available(macOS 11.0, *)) return NSLayoutAttributeLeading;
        return NSLayoutAttributeLeft;
    case WailsMacAccessoryLayoutTrailing:
        if (@available(macOS 11.0, *)) return NSLayoutAttributeTrailing;
        return NSLayoutAttributeRight;
    default:
        return NSLayoutAttributeBottom;
    }
}

void* macAccessoryCreate(int kind, int layout, double height) {
    NSViewController* controller = nil;
    if (kind == WailsMacAccessoryKindTitlebar) {
        NSTitlebarAccessoryViewController* titlebar = [[NSTitlebarAccessoryViewController alloc] init];
        titlebar.layoutAttribute = accessoryTitlebarAttribute(layout);
        controller = titlebar;
    } else if (kind == WailsMacAccessoryKindSplitItem) {
#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 260000
        if (@available(macOS 26.0, *)) {
            controller = [[NSSplitViewItemAccessoryViewController alloc] init];
        }
#endif
    }
    if (controller == nil) return NULL;

    WailsAccessoryHost* host = [[WailsAccessoryHost alloc] init];
    host.kind = kind;
    host.layout = layout;
    host.height = height;
    host.controls = [NSMutableDictionary dictionary];
    host.widths = [NSMutableDictionary dictionary];
    host.menus = [NSMutableDictionary dictionary];

    NSView* container = [[NSView alloc] initWithFrame:NSMakeRect(0, 0, 320, height)];
    NSStackView* stack = [[NSStackView alloc] initWithFrame:NSMakeRect(0, 0, 320, height)];
    stack.orientation = NSUserInterfaceLayoutOrientationHorizontal;
    stack.alignment = NSLayoutAttributeCenterY;
    stack.distribution = NSStackViewDistributionFill;
    stack.spacing = 8.0;
    stack.edgeInsets = NSEdgeInsetsMake(0, WailsAccessoryHorizontalInset, 0, WailsAccessoryHorizontalInset);
    stack.translatesAutoresizingMaskIntoConstraints = NO;
    [container addSubview:stack];
    [NSLayoutConstraint activateConstraints:@[
        [stack.leadingAnchor constraintEqualToAnchor:container.leadingAnchor],
        [stack.trailingAnchor constraintEqualToAnchor:container.trailingAnchor],
        [stack.centerYAnchor constraintEqualToAnchor:container.centerYAnchor],
        [stack.heightAnchor constraintLessThanOrEqualToAnchor:container.heightAnchor],
    ]];
    NSLayoutConstraint* heightConstraint = [container.heightAnchor constraintEqualToConstant:height];
    heightConstraint.priority = NSLayoutPriorityRequired - 1;
    heightConstraint.active = YES;
    host.heightConstraint = heightConstraint;
    host.container = container;
    host.stack = stack;
    [stack release];
    [container release];

    controller.view = container;
    objc_setAssociatedObject(controller, WailsAccessoryHostAssociationKey, host, OBJC_ASSOCIATION_RETAIN_NONATOMIC);
    [host release];
    return controller;
}

void macAccessoryRelease(void* controller) {
    if (controller == NULL) return;
    [(id)controller release];
}

bool macAccessoryAttachToWindow(void* nsWindow, void* controller) {
    NSWindow* window = (NSWindow*)nsWindow;
    if (window == nil || controller == NULL) return false;
    if (![(id)controller isKindOfClass:[NSTitlebarAccessoryViewController class]]) return false;
    [accessoryHost(controller) relayout];
    [window addTitlebarAccessoryViewController:(NSTitlebarAccessoryViewController*)controller];
    return true;
}

bool macAccessoryAttachToSplitItem(void* splitItem, void* controller, bool top) {
#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 260000
    if (@available(macOS 26.0, *)) {
        NSSplitViewItem* item = (NSSplitViewItem*)splitItem;
        if (item == nil || controller == NULL) return false;
        if (![(id)controller isKindOfClass:[NSSplitViewItemAccessoryViewController class]]) return false;
        [accessoryHost(controller) relayout];
        if (top) {
            [item addTopAlignedAccessoryViewController:(NSSplitViewItemAccessoryViewController*)controller];
        } else {
            [item addBottomAlignedAccessoryViewController:(NSSplitViewItemAccessoryViewController*)controller];
        }
        return true;
    }
#endif
    return false;
}

void macAccessoryDetachFromWindow(void* controller) {
    if (controller == NULL) return;
    if (![(id)controller isKindOfClass:[NSTitlebarAccessoryViewController class]]) return;
    NSTitlebarAccessoryViewController* accessory = (NSTitlebarAccessoryViewController*)controller;
    NSWindow* window = accessory.view.window;
    if (window == nil) return;
    NSUInteger index = [window.titlebarAccessoryViewControllers indexOfObjectIdenticalTo:accessory];
    if (index != NSNotFound) [window removeTitlebarAccessoryViewControllerAtIndex:(NSInteger)index];
}

void macAccessoryDetachFromSplitItem(void* splitItem, void* controller) {
#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 260000
    if (@available(macOS 26.0, *)) {
        NSSplitViewItem* item = (NSSplitViewItem*)splitItem;
        if (item == nil || controller == NULL) return;
        NSUInteger index = [item.topAlignedAccessoryViewControllers indexOfObjectIdenticalTo:(id)controller];
        if (index != NSNotFound) {
            [item removeTopAlignedAccessoryViewControllerAtIndex:(NSInteger)index];
            return;
        }
        index = [item.bottomAlignedAccessoryViewControllers indexOfObjectIdenticalTo:(id)controller];
        if (index != NSNotFound) [item removeBottomAlignedAccessoryViewControllerAtIndex:(NSInteger)index];
    }
#endif
}

void macAccessorySetHidden(void* controller, bool hidden) {
    if (controller == NULL) return;
    if ([(id)controller isKindOfClass:[NSTitlebarAccessoryViewController class]]) {
        if (@available(macOS 10.12, *)) {
            ((NSTitlebarAccessoryViewController*)controller).hidden = hidden;
        }
        return;
    }
#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 260000
    if (@available(macOS 26.0, *)) {
        if ([(id)controller isKindOfClass:[NSSplitViewItemAccessoryViewController class]]) {
            ((NSSplitViewItemAccessoryViewController*)controller).hidden = hidden;
        }
    }
#endif
}

void macAccessorySetHeight(void* controller, double height) {
    WailsAccessoryHost* host = accessoryHost(controller);
    if (host == nil || height <= 0) return;
    host.height = height;
    [host relayout];
    if ([(id)controller isKindOfClass:[NSTitlebarAccessoryViewController class]]) {
        // The window observes the accessory view's frame, so set it again
        // explicitly after the layout pass to make the new height stick.
        NSTitlebarAccessoryViewController* accessory = (NSTitlebarAccessoryViewController*)controller;
        NSRect frame = accessory.view.frame;
        frame.size.height = height;
        accessory.view.frame = frame;
    }
}

void macAccessorySetFullScreenMinHeight(void* controller, double height) {
    if (controller == NULL || ![(id)controller isKindOfClass:[NSTitlebarAccessoryViewController class]]) return;
    ((NSTitlebarAccessoryViewController*)controller).fullScreenMinHeight = MAX(0, height);
}

void macAccessorySetAutomaticallyAdjustsSize(void* controller, bool adjusts) {
    if (controller == NULL || ![(id)controller isKindOfClass:[NSTitlebarAccessoryViewController class]]) return;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
    if (@available(macOS 11.0, *)) {
        ((NSTitlebarAccessoryViewController*)controller).automaticallyAdjustsSize = adjusts;
    }
#endif
}

void macAccessorySetAppliesContentInsets(void* controller, bool applies) {
    if (controller == NULL) return;
#if __MAC_OS_X_VERSION_MAX_ALLOWED >= 260000
    if (@available(macOS 26.0, *)) {
        if ([(id)controller isKindOfClass:[NSSplitViewItemAccessoryViewController class]]) {
            ((NSSplitViewItemAccessoryViewController*)controller).automaticallyAppliesContentInsets = applies;
        }
    }
#endif
}

void macAccessoryAddSearch(void* controller, unsigned long long controlID,
    const char* placeholder, const char* text, bool incremental,
    const char* tooltip, bool disabled, bool hidden, double width) {
    WailsAccessoryHost* host = accessoryHost(controller);
    if (host == nil) return;
    NSSearchField* field = [[NSSearchField alloc] initWithFrame:NSMakeRect(0, 0, WailsAccessoryDefaultSearchWidth, 22)];
    field.placeholderString = accessoryString(placeholder);
    field.stringValue = accessoryString(text);
    field.controlSize = [host controlSize];
    field.font = [host controlFont];
    field.target = host;
    field.action = @selector(handleSearch:);
    field.sendsSearchStringImmediately = incremental;
    field.sendsWholeSearchString = !incremental;
    [host addControlView:field forControlID:controlID tooltip:tooltip disabled:disabled hidden:hidden width:width];
    [field release];
}

static void accessoryConfigureSegments(NSSegmentedControl* control, NSArray<NSString*>* labels,
    NSArray<NSString*>* symbols, int selected) {
    control.segmentCount = (NSInteger)labels.count;
    for (NSUInteger index = 0; index < labels.count; index++) {
        NSString* label = labels[index];
        NSString* symbol = index < symbols.count ? symbols[index] : @"";
        NSImage* image = accessorySymbolImage(symbol.UTF8String, label);
        if (image != nil) {
            [control setImage:image forSegment:(NSInteger)index];
            [control setImageScaling:NSImageScaleProportionallyDown forSegment:(NSInteger)index];
            [control setLabel:@"" forSegment:(NSInteger)index];
            if (@available(macOS 10.13, *)) [control setToolTip:label forSegment:(NSInteger)index];
        } else {
            [control setImage:nil forSegment:(NSInteger)index];
            [control setLabel:label forSegment:(NSInteger)index];
        }
        [control setWidth:0 forSegment:(NSInteger)index];
    }
    control.selectedSegment = selected >= 0 && selected < (int)labels.count ? selected : -1;
}

void macAccessoryAddSegmented(void* controller, unsigned long long controlID,
    const char* labelsJSON, const char* symbolsJSON, int selected,
    const char* tooltip, bool disabled, bool hidden, double width) {
    WailsAccessoryHost* host = accessoryHost(controller);
    if (host == nil) return;
    NSSegmentedControl* control = [[NSSegmentedControl alloc] initWithFrame:NSZeroRect];
    control.segmentStyle = NSSegmentStyleAutomatic;
    control.trackingMode = NSSegmentSwitchTrackingSelectOne;
    control.controlSize = [host controlSize];
    control.font = [host controlFont];
    accessoryConfigureSegments(control, accessoryStringsFromJSON(labelsJSON),
        accessoryStringsFromJSON(symbolsJSON), selected);
    control.target = host;
    control.action = @selector(handleSegment:);
    [host addControlView:control forControlID:controlID tooltip:tooltip disabled:disabled hidden:hidden width:width];
    [control release];
}

static void accessoryApplyButtonContent(NSButton* button, NSString* title, const char* symbol, bool menuButton) {
    NSImage* image = accessorySymbolImage(symbol, title.length > 0 ? title : accessoryString(symbol));
    button.title = title;
    if (image != nil) {
        button.image = image;
        button.imagePosition = title.length > 0 ? NSImageLeading : NSImageOnly;
    } else if (menuButton) {
        button.image = accessorySymbolImage("chevron.down", @"Menu");
        button.imagePosition = button.image != nil ? NSImageTrailing : NSNoImage;
    } else {
        button.image = nil;
        button.imagePosition = NSNoImage;
    }
}

void macAccessoryAddButton(void* controller, unsigned long long controlID,
    const char* title, const char* symbol, void* nsMenu,
    const char* tooltip, bool disabled, bool hidden, double width) {
    WailsAccessoryHost* host = accessoryHost(controller);
    if (host == nil) return;
    bool menuButton = nsMenu != NULL;
    NSButton* button = [[NSButton alloc] initWithFrame:NSZeroRect];
    button.bezelStyle = NSBezelStyleToolbar;
    button.controlSize = [host controlSize];
    button.font = [host controlFont];
    accessoryApplyButtonContent(button, accessoryString(title), symbol, menuButton);
    button.target = host;
    button.action = menuButton ? @selector(handleMenuButton:) : @selector(handleButton:);
    if (menuButton) host.menus[@(controlID)] = (NSMenu*)nsMenu;
    [host addControlView:button forControlID:controlID tooltip:tooltip disabled:disabled hidden:hidden width:width];
    [button release];
}

// A label is an inner stack of an optional symbol image and a text field so
// SetSymbol can add or remove the image later.
static NSImageView* accessoryLabelImageView(NSView* view) {
    if (![view isKindOfClass:[NSStackView class]]) return nil;
    NSArray<NSView*>* views = ((NSStackView*)view).arrangedSubviews;
    return views.count > 0 && [views[0] isKindOfClass:[NSImageView class]] ? (NSImageView*)views[0] : nil;
}

static NSTextField* accessoryLabelField(NSView* view) {
    if (![view isKindOfClass:[NSStackView class]]) return nil;
    for (NSView* candidate in ((NSStackView*)view).arrangedSubviews) {
        if ([candidate isKindOfClass:[NSTextField class]]) return (NSTextField*)candidate;
    }
    return nil;
}

void macAccessoryAddLabel(void* controller, unsigned long long controlID,
    const char* text, const char* symbol, const char* tooltip, bool hidden, double width) {
    WailsAccessoryHost* host = accessoryHost(controller);
    if (host == nil) return;
    NSStackView* group = [[NSStackView alloc] initWithFrame:NSZeroRect];
    group.orientation = NSUserInterfaceLayoutOrientationHorizontal;
    group.alignment = NSLayoutAttributeCenterY;
    group.spacing = 4.0;

    NSImageView* imageView = [[NSImageView alloc] initWithFrame:NSZeroRect];
    imageView.image = accessorySymbolImage(symbol, accessoryString(text));
    imageView.imageScaling = NSImageScaleProportionallyDown;
    if (@available(macOS 10.14, *)) imageView.contentTintColor = [NSColor secondaryLabelColor];
    imageView.hidden = imageView.image == nil;
    [group addArrangedSubview:imageView];
    [imageView release];

    NSTextField* field = [NSTextField labelWithString:accessoryString(text)];
    field.textColor = [NSColor secondaryLabelColor];
    field.font = [host controlFont];
    field.lineBreakMode = NSLineBreakByTruncatingTail;
    [field setContentCompressionResistancePriority:NSLayoutPriorityDefaultLow
                                    forOrientation:NSLayoutConstraintOrientationHorizontal];
    [group addArrangedSubview:field];

    [host addControlView:group forControlID:controlID tooltip:tooltip disabled:false hidden:hidden width:width];
    [group release];
}

void macAccessoryAddFlexibleSpace(void* controller, unsigned long long controlID) {
    WailsAccessoryHost* host = accessoryHost(controller);
    if (host == nil) return;
    NSView* spacer = [[NSView alloc] initWithFrame:NSZeroRect];
    [spacer setContentHuggingPriority:1 forOrientation:NSLayoutConstraintOrientationHorizontal];
    [spacer setContentCompressionResistancePriority:1 forOrientation:NSLayoutConstraintOrientationHorizontal];
    [host addControlView:spacer forControlID:controlID tooltip:NULL disabled:false hidden:false width:0];
    [[spacer.widthAnchor constraintGreaterThanOrEqualToConstant:0] setActive:YES];
    [spacer release];
}

void macAccessoryAddNativeView(void* controller, unsigned long long controlID,
    void* nsView, bool hidden, double width) {
    WailsAccessoryHost* host = accessoryHost(controller);
    NSView* view = (NSView*)nsView;
    if (host == nil || view == nil || ![view isKindOfClass:[NSView class]]) return;
    if (view.superview != nil) return;
    [host addControlView:view forControlID:controlID tooltip:NULL disabled:false hidden:hidden width:width];
}

void macAccessoryControlSetEnabled(void* controller, unsigned long long controlID, bool enabled) {
    NSView* view = [accessoryHost(controller) controlForID:controlID];
    if ([view isKindOfClass:[NSControl class]]) ((NSControl*)view).enabled = enabled;
}

void macAccessoryControlSetHidden(void* controller, unsigned long long controlID, bool hidden) {
    WailsAccessoryHost* host = accessoryHost(controller);
    NSView* view = [host controlForID:controlID];
    if (view == nil) return;
    view.hidden = hidden;
    [host relayout];
}

void macAccessoryControlSetTooltip(void* controller, unsigned long long controlID, const char* tooltip) {
    NSView* view = [accessoryHost(controller) controlForID:controlID];
    if (view == nil) return;
    view.toolTip = tooltip != NULL && strlen(tooltip) > 0 ? [NSString stringWithUTF8String:tooltip] : nil;
}

void macAccessoryControlSetWidth(void* controller, unsigned long long controlID, double width) {
    [accessoryHost(controller) setWidth:width forControlID:controlID];
}

void macAccessoryControlSetText(void* controller, unsigned long long controlID, const char* text) {
    WailsAccessoryHost* host = accessoryHost(controller);
    NSView* view = [host controlForID:controlID];
    if (view == nil) return;
    NSString* value = accessoryString(text);
    if ([view isKindOfClass:[NSSearchField class]]) {
        ((NSSearchField*)view).stringValue = value;
    } else if ([view isKindOfClass:[NSButton class]]) {
        NSButton* button = (NSButton*)view;
        button.title = value;
        if (button.image != nil) {
            button.imagePosition = value.length > 0
                ? (host.menus[@(controlID)] != nil && button.imagePosition == NSImageTrailing ? NSImageTrailing : NSImageLeading)
                : NSImageOnly;
        }
    } else {
        NSTextField* field = accessoryLabelField(view);
        if (field != nil) field.stringValue = value;
    }
    [host relayout];
}

void macAccessoryControlSetSymbol(void* controller, unsigned long long controlID, const char* symbol) {
    WailsAccessoryHost* host = accessoryHost(controller);
    NSView* view = [host controlForID:controlID];
    if (view == nil) return;
    if ([view isKindOfClass:[NSButton class]]) {
        NSButton* button = (NSButton*)view;
        accessoryApplyButtonContent(button, button.title, symbol, host.menus[@(controlID)] != nil);
    } else {
        NSImageView* imageView = accessoryLabelImageView(view);
        if (imageView != nil) {
            NSTextField* field = accessoryLabelField(view);
            imageView.image = accessorySymbolImage(symbol, field != nil ? field.stringValue : @"");
            imageView.hidden = imageView.image == nil;
        }
    }
    [host relayout];
}

void macAccessoryControlSetPlaceholder(void* controller, unsigned long long controlID, const char* placeholder) {
    NSView* view = [accessoryHost(controller) controlForID:controlID];
    if ([view isKindOfClass:[NSSearchField class]]) {
        ((NSSearchField*)view).placeholderString = accessoryString(placeholder);
    }
}

void macAccessoryControlSetIncremental(void* controller, unsigned long long controlID, bool incremental) {
    NSView* view = [accessoryHost(controller) controlForID:controlID];
    if (![view isKindOfClass:[NSSearchField class]]) return;
    NSSearchField* field = (NSSearchField*)view;
    field.sendsSearchStringImmediately = incremental;
    field.sendsWholeSearchString = !incremental;
}

void macAccessoryControlSetSegments(void* controller, unsigned long long controlID,
    const char* labelsJSON, const char* symbolsJSON, int selected) {
    WailsAccessoryHost* host = accessoryHost(controller);
    NSView* view = [host controlForID:controlID];
    if (![view isKindOfClass:[NSSegmentedControl class]]) return;
    accessoryConfigureSegments((NSSegmentedControl*)view, accessoryStringsFromJSON(labelsJSON),
        accessoryStringsFromJSON(symbolsJSON), selected);
    [host relayout];
}

void macAccessoryControlSetSelectedSegment(void* controller, unsigned long long controlID, int selected) {
    NSView* view = [accessoryHost(controller) controlForID:controlID];
    if (![view isKindOfClass:[NSSegmentedControl class]]) return;
    NSSegmentedControl* control = (NSSegmentedControl*)view;
    control.selectedSegment = selected >= 0 && selected < control.segmentCount ? selected : -1;
}

void macAccessoryControlSetMenu(void* controller, unsigned long long controlID, void* nsMenu) {
    WailsAccessoryHost* host = accessoryHost(controller);
    if (host == nil || [host controlForID:controlID] == nil) return;
    if (nsMenu == NULL) {
        [host.menus removeObjectForKey:@(controlID)];
        return;
    }
    host.menus[@(controlID)] = (NSMenu*)nsMenu;
}
