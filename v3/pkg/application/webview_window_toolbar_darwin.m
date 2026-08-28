//go:build darwin && !ios && !server

#import "webview_window_toolbar_darwin.h"
#import <dispatch/dispatch.h>
#import <objc/runtime.h>
#import <stdlib.h>
#import <string.h>

typedef struct {
    NSToolbar* toolbar;
    WailsToolbarDelegate* delegate;
    NSWindow* window;
} WailsToolbarHandle;

static const void* WailsToolbarDelegateAssociationKey = &WailsToolbarDelegateAssociationKey;
static const void* WailsToolbarSearchTargetAssociationKey = &WailsToolbarSearchTargetAssociationKey;
static const void* WailsToolbarGroupTargetAssociationKey = &WailsToolbarGroupTargetAssociationKey;
static const void* WailsToolbarGroupAssociationKey = &WailsToolbarGroupAssociationKey;
static const void* WailsToolbarGroupIndexAssociationKey = &WailsToolbarGroupIndexAssociationKey;
static const void* WailsToolbarShareTargetAssociationKey = &WailsToolbarShareTargetAssociationKey;
static const void* WailsToolbarShareProviderLifetimeAssociationKey = &WailsToolbarShareProviderLifetimeAssociationKey;

@implementation WailsToolbarItem

- (void)handleClick {
    processToolbarItemClick(self.itemID);
}

@end

@interface WailsToolbarTitleField : NSTextField
@property BOOL windowDraggable;
@end

@implementation WailsToolbarTitleField
- (BOOL)mouseDownCanMoveWindow { return self.windowDraggable; }
@end

@implementation WailsToolbarShareProviderLifetime

- (void)dealloc {
    processToolbarShareProviderRelease(self.providerID);
    [super dealloc];
}

@end

@implementation WailsToolbarGroupTarget

- (void)dealloc {
    [_itemIDs release];
    [super dealloc];
}

- (void)handleClick:(id)sender {
    NSInteger selectedIndex = -1;
    if ([sender isKindOfClass:[NSToolbarItemGroup class]]) {
        selectedIndex = ((NSToolbarItemGroup*)sender).selectedIndex;
    } else if ([sender isKindOfClass:[NSSegmentedControl class]]) {
        selectedIndex = ((NSSegmentedControl*)sender).selectedSegment;
    }
    if (selectedIndex >= 0 && selectedIndex < (NSInteger)self.itemIDs.count) {
        processToolbarItemClick(self.itemIDs[selectedIndex].unsignedIntValue);
    }
}

@end

@implementation WailsToolbarSearchTarget

- (void)handleSearch:(id)sender {
    NSSearchField* field = (NSSearchField*)sender;
    processToolbarSearch(self.itemID, (char*)field.stringValue.UTF8String);
}

@end

@implementation WailsToolbarShareTarget

- (void)dealloc {
    [_items release];
    [_subject release];
    [_activePicker release];
    [super dealloc];
}

- (NSArray*)itemsForSharingServicePickerToolbarItem:(NSSharingServicePickerToolbarItem*)pickerToolbarItem {
    return self.items ?: @[];
}

- (id<NSSharingServiceDelegate>)sharingServicePicker:(NSSharingServicePicker*)sharingServicePicker
    delegateForSharingService:(NSSharingService*)sharingService {
    if (self.subject.length > 0) sharingService.subject = self.subject;
    return self;
}

- (void)sharingServicePicker:(NSSharingServicePicker*)sharingServicePicker
    didChooseSharingService:(NSSharingService*)sharingService {
    if (sharingService != nil && self.subject.length > 0) {
        sharingService.subject = self.subject;
    }
    if (self.activePicker != nil) {
        // Keep the picker alive until the current AppKit callback unwinds.
        NSSharingServicePicker* closingPicker = [[self.activePicker retain] autorelease];
        self.activePicker = nil;
        (void)closingPicker;
    }
}

- (void)sharingService:(NSSharingService*)sharingService didShareItems:(NSArray*)items {
    NSString* service = sharingService.title ?: @"";
    processToolbarShareResult(self.itemID, (char*)service.UTF8String, (char*)"");
}

- (void)sharingService:(NSSharingService*)sharingService
    didFailToShareItems:(NSArray*)items error:(NSError*)error {
    NSString* service = sharingService.title ?: @"";
    NSString* message = error.localizedDescription ?: @"Sharing failed";
    processToolbarShareResult(self.itemID, (char*)service.UTF8String, (char*)message.UTF8String);
}

- (NSWindow*)sharingService:(NSSharingService*)sharingService
    sourceWindowForShareItems:(NSArray*)items sharingContentScope:(NSSharingContentScope*)sharingContentScope {
    if (sharingContentScope != NULL) *sharingContentScope = NSSharingContentScopeFull;
    return self.window;
}

- (void)showSharePicker:(id)sender {
    if (self.items.count == 0) return;
    NSToolbarItem* toolbarItem = [sender isKindOfClass:[NSToolbarItem class]] ? sender : nil;
    NSView* anchor = toolbarItem.view;
    if (anchor == nil) anchor = self.window.contentView;
    if (anchor == nil) return;

    NSSharingServicePicker* picker = [[NSSharingServicePicker alloc] initWithItems:self.items];
    picker.delegate = self;
    self.activePicker = picker;
    [picker release];

    NSRect rect = anchor.bounds;
    if (toolbarItem.view == nil) {
        rect = NSMakeRect(NSMaxX(anchor.bounds) - 1, NSMaxY(anchor.bounds) - 1, 1, 1);
    }
    [self.activePicker showRelativeToRect:rect ofView:anchor preferredEdge:NSMinYEdge];
}

@end

@implementation WailsToolbarDelegate

- (instancetype)init {
    self = [super init];
    if (self) {
        _orderedIdentifiers = [[NSMutableArray alloc] init];
        _itemsByIdentifier = [[NSMutableDictionary alloc] init];
    }
    return self;
}

- (void)dealloc {
    [_orderedIdentifiers release];
    [_itemsByIdentifier release];
    [super dealloc];
}

- (NSArray<NSToolbarItemIdentifier>*)toolbarDefaultItemIdentifiers:(NSToolbar*)toolbar {
    return self.orderedIdentifiers;
}

- (NSArray<NSToolbarItemIdentifier>*)toolbarAllowedItemIdentifiers:(NSToolbar*)toolbar {
    return self.orderedIdentifiers;
}

- (NSToolbarItem*)toolbar:(NSToolbar*)toolbar
        itemForItemIdentifier:(NSToolbarItemIdentifier)itemIdentifier
    willBeInsertedIntoToolbar:(BOOL)flag {
    return self.itemsByIdentifier[itemIdentifier];
}

@end

static WailsToolbarHandle* toolbarHandle(void* handlePtr) {
    return (WailsToolbarHandle*)handlePtr;
}

static WailsToolbarDelegate* toolbarDelegate(void* handlePtr) {
    WailsToolbarHandle* handle = toolbarHandle(handlePtr);
    return handle == NULL ? nil : handle->delegate;
}

static NSImage* toolbarSymbolImage(const char* symbolName, NSString* accessibilityLabel) {
    if (symbolName == NULL || strlen(symbolName) == 0) return nil;
    if (@available(macOS 11.0, *)) {
        return [NSImage imageWithSystemSymbolName:[NSString stringWithUTF8String:symbolName]
                         accessibilityDescription:accessibilityLabel];
    }
    return nil;
}

static void applyCommonItemStyle(NSToolbarItem* item, const char* tooltip,
    bool bordered, bool prominent, bool disabled, bool hidden,
    bool hasTint, double tintR, double tintG, double tintB, double tintA, int badgeCount) {
    item.toolTip = tooltip != NULL && strlen(tooltip) > 0
        ? [NSString stringWithUTF8String:tooltip]
        : nil;
    // Wails owns enabled state through MacToolbarItem.SetEnabled. AppKit's
    // target/action auto-validation otherwise re-enables a disabled item as
    // soon as it discovers that the target implements the action selector.
    item.autovalidates = NO;
    item.enabled = !disabled;
    item.hidden = hidden;
    if (@available(macOS 10.15, *)) {
        item.bordered = bordered;
    }
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 260000
    if (@available(macOS 26.0, *)) {
        item.style = prominent ? NSToolbarItemStyleProminent : NSToolbarItemStylePlain;
        item.backgroundTintColor = hasTint
            ? [NSColor colorWithRed:tintR green:tintG blue:tintB alpha:tintA]
            : nil;
        item.badge = badgeCount > 0 ? [NSItemBadge badgeWithCount:badgeCount] : nil;
    }
#endif
}

void* toolbarCreate(const char* identifier) {
    if (identifier == NULL || strlen(identifier) == 0) return NULL;

    WailsToolbarHandle* handle = calloc(1, sizeof(WailsToolbarHandle));
    if (handle == NULL) return NULL;

    WailsToolbarDelegate* delegate = [[WailsToolbarDelegate alloc] init];
    NSToolbar* toolbar = [[NSToolbar alloc]
        initWithIdentifier:[NSString stringWithUTF8String:identifier]];
    if (delegate == nil || toolbar == nil) {
        [delegate release];
        [toolbar release];
        free(handle);
        return NULL;
    }

    toolbar.displayMode = NSToolbarDisplayModeIconAndLabel;
    toolbar.allowsUserCustomization = NO;
    toolbar.autosavesConfiguration = NO;
    toolbar.visible = YES;

    // NSToolbar.delegate is weak/assign. The association makes the delegate's
    // lifetime exactly match the toolbar while the handle owns the toolbar.
    objc_setAssociatedObject(toolbar, WailsToolbarDelegateAssociationKey, delegate, OBJC_ASSOCIATION_RETAIN);
    [delegate release];

    handle->toolbar = toolbar; // +1 from alloc
    handle->delegate = delegate; // retained by the toolbar association
    return handle;
}

void toolbarAttach(void* nsWindow, void* handlePtr, int style) {
    NSWindow* window = (NSWindow*)nsWindow;
    WailsToolbarHandle* handle = toolbarHandle(handlePtr);
    if (window == nil || handle == NULL || handle->toolbar == nil || handle->delegate == nil) return;

    NSToolbar* toolbar = handle->toolbar;
    WailsToolbarDelegate* delegate = handle->delegate;

    // Install the delegate only after Go has populated every item. AppKit's
    // initial default-identifier request therefore sees the complete tree.
    toolbar.delegate = delegate;
    NSArray<NSToolbarItemIdentifier>* identifiers = delegate.orderedIdentifiers;
    if (@available(macOS 15.0, *)) {
        toolbar.itemIdentifiers = identifiers;
    } else {
        for (NSToolbarItemIdentifier identifier in identifiers) {
            [toolbar insertItemWithItemIdentifier:identifier atIndex:toolbar.items.count];
        }
    }

    for (NSToolbarItem* item in delegate.itemsByIdentifier.allValues) {
        WailsToolbarShareTarget* shareTarget = objc_getAssociatedObject(item, WailsToolbarShareTargetAssociationKey);
        if (shareTarget != nil) shareTarget.window = window;
    }

    window.toolbar = toolbar;
    handle->window = window;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
    if (@available(macOS 11.0, *)) {
        window.toolbarStyle = style;
    }
#endif
    [toolbar validateVisibleItems];
    toolbar.visible = YES;
}

void toolbarRelease(void* handlePtr) {
    WailsToolbarHandle* handle = toolbarHandle(handlePtr);
    if (handle == NULL) return;
    handle->window = nil;
    [handle->toolbar release];
    handle->toolbar = nil;
    handle->delegate = nil;
    free(handle);
}

void toolbarDetach(void* nsWindow) {
    NSWindow* window = (NSWindow*)nsWindow;
    window.toolbar = nil;
}

void toolbarSetDisplayMode(void* handlePtr, int displayMode) {
    WailsToolbarHandle* handle = toolbarHandle(handlePtr);
    if (handle == NULL || handle->toolbar == nil) return;
    if (displayMode < NSToolbarDisplayModeDefault || displayMode > NSToolbarDisplayModeLabelOnly) return;
    handle->toolbar.displayMode = (NSToolbarDisplayMode)displayMode;
}

// Returns a +1 retained item. A caller must transfer it to an owning
// collection and release this reference.
void* toolbarBuildButtonItemStandalone(const char* identifier, unsigned int itemID,
    const char* label, const char* symbolName, const char* tooltip,
    bool bordered, bool disabled, bool hidden) {
    NSString* identifierString = [NSString stringWithUTF8String:identifier];
    WailsToolbarItem* item = [[WailsToolbarItem alloc] initWithItemIdentifier:identifierString];
    item.itemID = itemID;
    item.label = [NSString stringWithUTF8String:label];
    item.target = item;
    item.action = @selector(handleClick);
    item.image = toolbarSymbolImage(symbolName, item.label);
    applyCommonItemStyle(item, tooltip, bordered, false, disabled, hidden,
        false, 0, 0, 0, 0, 0);
    return item;
}

void* toolbarAddButtonItem(void* handlePtr, const char* identifier, unsigned int itemID,
    const char* label, const char* symbolName, const char* tooltip,
    bool bordered, bool prominent, bool disabled, bool hidden,
    bool hasTint, double tintR, double tintG, double tintB, double tintA, int badgeCount) {
    WailsToolbarDelegate* delegate = toolbarDelegate(handlePtr);
    if (delegate == nil) return NULL;

    WailsToolbarItem* item = (WailsToolbarItem*)toolbarBuildButtonItemStandalone(
        identifier, itemID, label, symbolName, tooltip, bordered, disabled, hidden);
    applyCommonItemStyle(item, tooltip, bordered, prominent, disabled, hidden,
        hasTint, tintR, tintG, tintB, tintA, badgeCount);

    NSString* identifierString = [NSString stringWithUTF8String:identifier];
    [delegate.orderedIdentifiers addObject:identifierString];
    delegate.itemsByIdentifier[identifierString] = item;
    [item release];
    return item;
}

void* toolbarAddGroupItem(void* handlePtr, const char* identifier,
    const char* label, void** memberItems, int memberCount, int selectionMode, int selectedIndex) {
    WailsToolbarDelegate* delegate = toolbarDelegate(handlePtr);
    if (delegate == nil || memberCount <= 0) return NULL;

    NSMutableArray<NSToolbarItem*>* sourceItems = [NSMutableArray arrayWithCapacity:memberCount];
    for (int i = 0; i < memberCount; i++) {
        NSToolbarItem* memberItem = (NSToolbarItem*)memberItems[i];
        [sourceItems addObject:memberItem];
    }

    NSString* identifierString = [NSString stringWithUTF8String:identifier];
    NSToolbarItemGroup* group = nil;
    if (@available(macOS 10.15, *)) {
        NSToolbarItemGroupSelectionMode nativeSelectionMode;
        switch (selectionMode) {
            case 1: nativeSelectionMode = NSToolbarItemGroupSelectionModeMomentary; break;
            case 2: nativeSelectionMode = NSToolbarItemGroupSelectionModeSelectAny; break;
            default: nativeSelectionMode = NSToolbarItemGroupSelectionModeSelectOne; break;
        }

        NSMutableArray<NSString*>* titles = [NSMutableArray arrayWithCapacity:memberCount];
        NSMutableArray<NSString*>* labels = [NSMutableArray arrayWithCapacity:memberCount];
        NSMutableArray<NSImage*>* images = [NSMutableArray arrayWithCapacity:memberCount];
        BOOL hasAllImages = YES;
        NSMutableArray<NSNumber*>* itemIDs = [NSMutableArray arrayWithCapacity:memberCount];
        for (WailsToolbarItem* source in sourceItems) {
            [titles addObject:source.label ?: @""];
            [labels addObject:source.label ?: @""];
            [itemIDs addObject:@(source.itemID)];
            if (source.image != nil) {
                [images addObject:source.image];
            } else {
                hasAllImages = NO;
            }
        }

        WailsToolbarGroupTarget* target = [[WailsToolbarGroupTarget alloc] init];
        target.itemIDs = itemIDs;
        if (@available(macOS 11.0, *)) {
            if (hasAllImages) {
                group = [[NSToolbarItemGroup groupWithItemIdentifier:identifierString
                    images:images selectionMode:nativeSelectionMode labels:labels
                    target:target action:@selector(handleClick:)] retain];
            }
        }
        if (group == nil) {
            group = [[NSToolbarItemGroup groupWithItemIdentifier:identifierString
                titles:titles selectionMode:nativeSelectionMode labels:labels
                target:target action:@selector(handleClick:)] retain];
        }

        objc_setAssociatedObject(group, WailsToolbarGroupTargetAssociationKey, target, OBJC_ASSOCIATION_RETAIN);
        [target release];

        // The convenience constructor creates the actual segmented subitems.
        // Map Wails' private identifiers to those items so every live setter
        // continues to address the correct segment.
        for (int i = 0; i < memberCount; i++) {
            NSToolbarItem* source = sourceItems[i];
            NSToolbarItem* actual = group.subitems[i];
            actual.label = source.label;
            actual.image = source.image;
            actual.toolTip = source.toolTip;
            actual.autovalidates = NO;
            actual.enabled = source.enabled;
            actual.hidden = source.hidden;
            objc_setAssociatedObject(actual, WailsToolbarGroupAssociationKey,
                [NSValue valueWithNonretainedObject:group], OBJC_ASSOCIATION_RETAIN);
            objc_setAssociatedObject(actual, WailsToolbarGroupIndexAssociationKey,
                @(i), OBJC_ASSOCIATION_RETAIN);
            delegate.itemsByIdentifier[source.itemIdentifier] = actual;
            [source release];
        }

        if (selectedIndex >= -1 && selectedIndex < memberCount) group.selectedIndex = selectedIndex;
        group.controlRepresentation = NSToolbarItemGroupControlRepresentationAutomatic;
    } else {
        // Older AppKit still supports grouped toolbar items, but not the
        // segmented selection API. Each source item keeps its own target and
        // callback, preserving functional behavior without newer selectors.
        group = [[NSToolbarItemGroup alloc] initWithItemIdentifier:identifierString];
        group.subitems = sourceItems;
        for (NSToolbarItem* source in sourceItems) {
            delegate.itemsByIdentifier[source.itemIdentifier] = source;
            [source release];
        }
    }
    group.autovalidates = NO;
    group.label = [NSString stringWithUTF8String:label];

    [delegate.orderedIdentifiers addObject:identifierString];
    delegate.itemsByIdentifier[identifierString] = group;
    [group release];
    return group;
}

void* toolbarAddSearchItem(void* handlePtr, const char* identifier, unsigned int itemID,
    const char* label, const char* tooltip, bool disabled, bool hidden) {
    WailsToolbarDelegate* delegate = toolbarDelegate(handlePtr);
    if (delegate == nil) return NULL;

    NSString* identifierString = [NSString stringWithUTF8String:identifier];
    // Use a regular toolbar item rather than NSSearchToolbarItem. The latter
    // intentionally grows to consume available unified-toolbar space and does
    // not reliably honour maxSize on recent macOS releases. A view-backed
    // search field stays compact until the application explicitly replaces it.
    NSToolbarItem* item = [[NSToolbarItem alloc] initWithItemIdentifier:identifierString];
    NSSearchField* field = [[NSSearchField alloc] initWithFrame:NSMakeRect(0, 0, 160, 22)];
    item.view = field;
    [field release];

    field.controlSize = NSControlSizeSmall;
    field.frameSize = NSMakeSize(160, 22);
    item.minSize = NSMakeSize(120, 22);
    item.maxSize = NSMakeSize(160, 22);

    WailsToolbarSearchTarget* target = [[WailsToolbarSearchTarget alloc] init];
    target.itemID = itemID;
    objc_setAssociatedObject(item, WailsToolbarSearchTargetAssociationKey, target, OBJC_ASSOCIATION_RETAIN);
    [target release];

    item.label = [NSString stringWithUTF8String:label];
    item.toolTip = tooltip != NULL && strlen(tooltip) > 0
        ? [NSString stringWithUTF8String:tooltip]
        : nil;
    item.enabled = !disabled;
    item.hidden = hidden;
    field.target = target;
    field.action = @selector(handleSearch:);
    field.sendsWholeSearchString = YES;
    field.enabled = !disabled;

    [delegate.orderedIdentifiers addObject:identifierString];
    delegate.itemsByIdentifier[identifierString] = item;
    [item release];
    return item;
}

void* toolbarAddTitleItem(void* handlePtr, const char* identifier,
    const char* label, bool draggable, bool hidden) {
    WailsToolbarDelegate* delegate = toolbarDelegate(handlePtr);
    if (delegate == nil) return NULL;

    NSString* identifierString = [NSString stringWithUTF8String:identifier];
    NSString* title = [NSString stringWithUTF8String:label];
    NSToolbarItem* item = [[NSToolbarItem alloc] initWithItemIdentifier:identifierString];
    WailsToolbarTitleField* field = [WailsToolbarTitleField labelWithString:title];
    field.windowDraggable = draggable;
    field.font = [NSFont systemFontOfSize:15 weight:NSFontWeightSemibold];
    field.lineBreakMode = NSLineBreakByTruncatingTail;
    field.maximumNumberOfLines = 1;
    field.frame = NSMakeRect(0, 0, 180, 24);
    item.view = field;
    item.label = title;
    item.minSize = NSMakeSize(80, 24);
    item.maxSize = NSMakeSize(180, 24);
    item.visibilityPriority = NSToolbarItemVisibilityPriorityHigh;
    item.hidden = hidden;

    [delegate.orderedIdentifiers addObject:identifierString];
    delegate.itemsByIdentifier[identifierString] = item;
    [item release];
    return item;
}

static NSError* toolbarShareError(NSString* message) {
    return [NSError errorWithDomain:@"WailsMacShareError" code:1
        userInfo:@{NSLocalizedDescriptionKey: message ?: @"Unable to provide share data"}];
}

static void toolbarShareTargetSetProvider(WailsToolbarShareTarget* target, const char* providerJSON) {
    target.items = @[];
    target.subject = @"";
    if (providerJSON == NULL || strlen(providerJSON) == 0) return;

    NSData* data = [[NSString stringWithUTF8String:providerJSON] dataUsingEncoding:NSUTF8StringEncoding];
    if (data == nil) return;
    id decoded = [NSJSONSerialization JSONObjectWithData:data options:0 error:nil];
    if (![decoded isKindOfClass:[NSDictionary class]]) return;
    NSDictionary* descriptor = (NSDictionary*)decoded;
    id providerIDValue = descriptor[@"providerID"];
    unsigned int providerID = [providerIDValue isKindOfClass:[NSNumber class]]
        ? ((NSNumber*)providerIDValue).unsignedIntValue
        : 0;
    id subject = descriptor[@"subject"];
    if ([subject isKindOfClass:[NSString class]]) target.subject = subject;

    id encodedRepresentations = descriptor[@"representations"];
    if (![encodedRepresentations isKindOfClass:[NSArray class]] ||
        ((NSArray*)encodedRepresentations).count == 0) return;

    NSItemProvider* provider = [[NSItemProvider alloc] init];
    if (providerID > 0) {
        WailsToolbarShareProviderLifetime* lifetime = [[WailsToolbarShareProviderLifetime alloc] init];
        lifetime.providerID = providerID;
        objc_setAssociatedObject(provider, WailsToolbarShareProviderLifetimeAssociationKey,
            lifetime, OBJC_ASSOCIATION_RETAIN);
        [lifetime release];
    }
    id suggestedName = descriptor[@"suggestedName"];
    if (@available(macOS 10.14, *)) {
        if ([suggestedName isKindOfClass:[NSString class]] &&
            ((NSString*)suggestedName).length > 0) {
            provider.suggestedName = suggestedName;
        }
    }
    __block BOOL registeredRepresentation = NO;
    for (id encodedRepresentation in (NSArray*)encodedRepresentations) {
        if (![encodedRepresentation isKindOfClass:[NSDictionary class]]) continue;
        id contentTypeValue = encodedRepresentation[@"contentType"];
        if (![contentTypeValue isKindOfClass:[NSString class]] ||
            ((NSString*)contentTypeValue).length == 0) continue;

        NSString* contentType = [(NSString*)contentTypeValue copy];
        [provider registerDataRepresentationForTypeIdentifier:contentType
            visibility:NSItemProviderRepresentationVisibilityAll
            loadHandler:^NSProgress* (void (^completionHandler)(NSData*, NSError*)) {
                dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
                    char* responseJSON = processToolbarShareData(providerID, (char*)contentType.UTF8String);
                    if (responseJSON == NULL) {
                        completionHandler(nil, toolbarShareError(@"The share provider returned no response"));
                        return;
                    }

                    NSData* responseData = [[NSString stringWithUTF8String:responseJSON]
                        dataUsingEncoding:NSUTF8StringEncoding];
                    free(responseJSON);
                    NSError* decodeError = nil;
                    id response = responseData == nil ? nil :
                        [NSJSONSerialization JSONObjectWithData:responseData options:0 error:&decodeError];
                    if (![response isKindOfClass:[NSDictionary class]]) {
                        completionHandler(nil, toolbarShareError(
                            decodeError.localizedDescription ?: @"The share provider returned an invalid response"));
                        return;
                    }

                    id errorMessage = response[@"error"];
                    if ([errorMessage isKindOfClass:[NSString class]] &&
                        ((NSString*)errorMessage).length > 0) {
                        completionHandler(nil, toolbarShareError(errorMessage));
                        return;
                    }

                    id encodedData = response[@"data"];
                    if (![encodedData isKindOfClass:[NSString class]]) {
                        completionHandler(nil, toolbarShareError(@"The share provider returned no data"));
                        return;
                    }
                    NSData* providedData = [[NSData alloc]
                        initWithBase64EncodedString:encodedData options:0];
                    if (providedData == nil) {
                        completionHandler(nil, toolbarShareError(@"The share provider returned invalid data"));
                        return;
                    }
                    completionHandler(providedData, nil);
                    [providedData release];
                });
                return nil;
            }];
        registeredRepresentation = YES;
        [contentType release];
    }
    if (registeredRepresentation) target.items = @[provider];
    [provider release];
}

void* toolbarAddShareItem(void* handlePtr, const char* identifier, unsigned int itemID,
    const char* label, const char* symbolName, const char* tooltip,
    bool disabled, bool hidden, const char* providerJSON) {
    WailsToolbarDelegate* delegate = toolbarDelegate(handlePtr);
    if (delegate == nil) return NULL;

    NSString* identifierString = [NSString stringWithUTF8String:identifier];
    WailsToolbarShareTarget* target = [[WailsToolbarShareTarget alloc] init];
    target.itemID = itemID;
    toolbarShareTargetSetProvider(target, providerJSON);

    NSToolbarItem* item = nil;
    if (@available(macOS 10.15, *)) {
        NSSharingServicePickerToolbarItem* sharingItem =
            [[NSSharingServicePickerToolbarItem alloc] initWithItemIdentifier:identifierString];
        sharingItem.delegate = target;
        item = sharingItem;
    } else {
        item = [[NSToolbarItem alloc] initWithItemIdentifier:identifierString];
        item.target = target;
        item.action = @selector(showSharePicker:);
    }

    item.label = [NSString stringWithUTF8String:label];
    NSImage* symbol = toolbarSymbolImage(symbolName, item.label);
    if (symbol != nil) item.image = symbol;
    applyCommonItemStyle(item, tooltip, false, false, disabled || target.items.count == 0, hidden,
        false, 0, 0, 0, 0, 0);
    objc_setAssociatedObject(item, WailsToolbarShareTargetAssociationKey, target, OBJC_ASSOCIATION_RETAIN);
    [target release];

    [delegate.orderedIdentifiers addObject:identifierString];
    delegate.itemsByIdentifier[identifierString] = item;
    [item release];
    return item;
}

void toolbarAddFlexibleSpaceIdentifier(void* handlePtr) {
    WailsToolbarDelegate* delegate = toolbarDelegate(handlePtr);
    if (delegate != nil) {
        [delegate.orderedIdentifiers addObject:NSToolbarFlexibleSpaceItemIdentifier];
    }
}

void toolbarAddSeparatorIdentifier(void* handlePtr) {
    WailsToolbarDelegate* delegate = toolbarDelegate(handlePtr);
    if (delegate == nil) return;
    [delegate.orderedIdentifiers addObject:NSToolbarSeparatorItemIdentifier];
}

void toolbarAddSidebarToggleIdentifier(void* handlePtr) {
    WailsToolbarDelegate* delegate = toolbarDelegate(handlePtr);
    if (delegate != nil) {
        [delegate.orderedIdentifiers addObject:NSToolbarToggleSidebarItemIdentifier];
    }
}

void toolbarAddSidebarTrackingSeparatorIdentifier(void* handlePtr) {
    WailsToolbarDelegate* delegate = toolbarDelegate(handlePtr);
    if (delegate == nil) {
        return;
    }
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
    if (@available(macOS 11.0, *)) {
        [delegate.orderedIdentifiers addObject:NSToolbarSidebarTrackingSeparatorItemIdentifier];
    }
#endif
}

void toolbarAddInspectorToggleItem(void* handlePtr, const char* identifier, unsigned int itemID) {
    WailsToolbarDelegate* delegate = toolbarDelegate(handlePtr);
    if (delegate == nil) return;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 140000
    if (@available(macOS 14.0, *)) {
        [delegate.orderedIdentifiers addObject:NSToolbarToggleInspectorItemIdentifier];
        return;
    }
#endif

    WailsToolbarItem* item = (WailsToolbarItem*)toolbarBuildButtonItemStandalone(
        identifier, itemID, "Inspector", "sidebar.trailing", "Show or hide the inspector",
        true, false, false);
    if (item.image == nil) item.image = [NSImage imageNamed:NSImageNameInfo];
    NSString* identifierString = [NSString stringWithUTF8String:identifier];
    [delegate.orderedIdentifiers addObject:identifierString];
    delegate.itemsByIdentifier[identifierString] = item;
    [item release];
}

void toolbarAddInspectorTrackingSeparatorIdentifier(void* handlePtr) {
    WailsToolbarDelegate* delegate = toolbarDelegate(handlePtr);
    if (delegate == nil) return;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 140000
    if (@available(macOS 14.0, *)) {
        [delegate.orderedIdentifiers addObject:NSToolbarInspectorTrackingSeparatorItemIdentifier];
    }
#endif
}

static NSToolbarItem* toolbarItemForIdentifier(void* handlePtr, const char* identifier) {
    WailsToolbarDelegate* delegate = toolbarDelegate(handlePtr);
    if (delegate == nil || identifier == NULL) return nil;
    return delegate.itemsByIdentifier[[NSString stringWithUTF8String:identifier]];
}

static NSSegmentedControl* toolbarSegmentedControlInView(NSView* view) {
    if ([view isKindOfClass:[NSSegmentedControl class]]) {
        return (NSSegmentedControl*)view;
    }
    for (NSView* subview in view.subviews) {
        NSSegmentedControl* control = toolbarSegmentedControlInView(subview);
        if (control != nil) return control;
    }
    return nil;
}

void toolbarShareItemSetProvider(void* handlePtr, const char* identifier, const char* providerJSON) {
    NSToolbarItem* item = toolbarItemForIdentifier(handlePtr, identifier);
    WailsToolbarShareTarget* target = objc_getAssociatedObject(item, WailsToolbarShareTargetAssociationKey);
    if (target != nil) toolbarShareTargetSetProvider(target, providerJSON);
}

void toolbarItemSetLabel(void* handlePtr, const char* identifier, const char* label) {
    NSToolbarItem* item = toolbarItemForIdentifier(handlePtr, identifier);
    if (item != nil) {
        NSString* value = [NSString stringWithUTF8String:label];
        item.label = value;
        if ([item.view isKindOfClass:[NSTextField class]]) {
            ((NSTextField*)item.view).stringValue = value;
            [item.view setNeedsLayout:YES];
        }
    }
}

void toolbarItemSetSymbol(void* handlePtr, const char* identifier, const char* symbolName) {
    NSToolbarItem* item = toolbarItemForIdentifier(handlePtr, identifier);
    if (item != nil) item.image = toolbarSymbolImage(symbolName, item.label);
}

void toolbarItemSetTooltip(void* handlePtr, const char* identifier, const char* tooltip) {
    NSToolbarItem* item = toolbarItemForIdentifier(handlePtr, identifier);
    if (item != nil) {
        item.toolTip = tooltip != NULL && strlen(tooltip) > 0
            ? [NSString stringWithUTF8String:tooltip]
            : nil;
    }
}

void toolbarItemSetBordered(void* handlePtr, const char* identifier, bool bordered) {
    NSToolbarItem* item = toolbarItemForIdentifier(handlePtr, identifier);
    if (@available(macOS 10.15, *)) {
        if (item != nil) item.bordered = bordered;
    }
}

void toolbarItemSetProminent(void* handlePtr, const char* identifier, bool prominent) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 260000
    NSToolbarItem* item = toolbarItemForIdentifier(handlePtr, identifier);
    if (@available(macOS 26.0, *)) {
        if (item != nil) item.style = prominent ? NSToolbarItemStyleProminent : NSToolbarItemStylePlain;
    }
#endif
}

void toolbarItemSetTintColor(void* handlePtr, const char* identifier, bool hasTint,
    double tintR, double tintG, double tintB, double tintA) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 260000
    NSToolbarItem* item = toolbarItemForIdentifier(handlePtr, identifier);
    if (@available(macOS 26.0, *)) {
        if (item != nil) {
            item.backgroundTintColor = hasTint
                ? [NSColor colorWithRed:tintR green:tintG blue:tintB alpha:tintA]
                : nil;
        }
    }
#endif
}

void toolbarItemSetEnabled(void* handlePtr, const char* identifier, bool enabled) {
    NSToolbarItem* item = toolbarItemForIdentifier(handlePtr, identifier);
    if (item != nil) {
        WailsToolbarShareTarget* shareTarget = objc_getAssociatedObject(item, WailsToolbarShareTargetAssociationKey);
        if (shareTarget != nil) enabled = enabled && shareTarget.items.count > 0;
        item.enabled = enabled;
        if ([item.view isKindOfClass:[NSControl class]]) {
            ((NSControl*)item.view).enabled = enabled;
        }
        NSValue* groupValue = objc_getAssociatedObject(item, WailsToolbarGroupAssociationKey);
        NSNumber* segmentIndex = objc_getAssociatedObject(item, WailsToolbarGroupIndexAssociationKey);
        NSToolbarItemGroup* group = groupValue.nonretainedObjectValue;
        NSSegmentedControl* control = toolbarSegmentedControlInView(group.view);
        if (control != nil && segmentIndex != nil) {
            [control setEnabled:enabled forSegment:segmentIndex.integerValue];
        }
    }
}

void toolbarItemSetHidden(void* handlePtr, const char* identifier, bool hidden) {
    NSToolbarItem* item = toolbarItemForIdentifier(handlePtr, identifier);
    if (item != nil) item.hidden = hidden;
}

void toolbarItemSetBadgeCount(void* handlePtr, const char* identifier, int badgeCount) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 260000
    NSToolbarItem* item = toolbarItemForIdentifier(handlePtr, identifier);
    if (@available(macOS 26.0, *)) {
        if (item != nil) item.badge = badgeCount > 0 ? [NSItemBadge badgeWithCount:badgeCount] : nil;
    }
#endif
}

void toolbarItemSetDraggable(void* handlePtr, const char* identifier, bool draggable) {
    NSToolbarItem* item = toolbarItemForIdentifier(handlePtr, identifier);
    if ([item.view isKindOfClass:[WailsToolbarTitleField class]]) {
        ((WailsToolbarTitleField*)item.view).windowDraggable = draggable;
    }
}

void toolbarGroupSetSelectedIndex(void* handlePtr, const char* identifier, int index) {
    NSToolbarItem* item = toolbarItemForIdentifier(handlePtr, identifier);
    if ([item isKindOfClass:[NSToolbarItemGroup class]]) {
        NSToolbarItemGroup* group = (NSToolbarItemGroup*)item;
        if (index >= -1 && index < (int)group.subitems.count) group.selectedIndex = index;
    }
}

void toolbarGroupSetSelectionMode(void* handlePtr, const char* identifier, int selectionMode) {
    NSToolbarItem* item = toolbarItemForIdentifier(handlePtr, identifier);
    if (@available(macOS 10.15, *)) {
        if ([item isKindOfClass:[NSToolbarItemGroup class]]) {
            NSToolbarItemGroup* group = (NSToolbarItemGroup*)item;
            switch (selectionMode) {
                case 1: group.selectionMode = NSToolbarItemGroupSelectionModeMomentary; break;
                case 2: group.selectionMode = NSToolbarItemGroupSelectionModeSelectAny; break;
                default: group.selectionMode = NSToolbarItemGroupSelectionModeSelectOne; break;
            }
        }
    }
}
