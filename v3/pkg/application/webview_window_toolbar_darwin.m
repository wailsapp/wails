//go:build darwin && !ios && !server

#import "webview_window_toolbar_darwin.h"
#import <dispatch/dispatch.h>
#import <objc/runtime.h>
#import <stdlib.h>
#import <string.h>

typedef struct {
    NSToolbar* toolbar;
    WailsToolbarDelegate* delegate;
} WailsToolbarHandle;

static const void* WailsToolbarDelegateAssociationKey = &WailsToolbarDelegateAssociationKey;
static const void* WailsToolbarSearchTargetAssociationKey = &WailsToolbarSearchTargetAssociationKey;
static const void* WailsToolbarGroupTargetAssociationKey = &WailsToolbarGroupTargetAssociationKey;
static const void* WailsToolbarShareTargetAssociationKey = &WailsToolbarShareTargetAssociationKey;
static const void* WailsToolbarShareProviderLifetimeAssociationKey = &WailsToolbarShareProviderLifetimeAssociationKey;

@implementation WailsToolbarItem

- (void)handleClick {
    processToolbarItemClick(self.itemID);
}

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
        _knownIdentifiers = [[NSMutableSet alloc] init];
    }
    return self;
}

- (void)dealloc {
    [_orderedIdentifiers release];
    [_allowedIdentifiers release];
    [_knownIdentifiers release];
    [_itemsByIdentifier release];
    [super dealloc];
}

- (NSArray<NSToolbarItemIdentifier>*)toolbarDefaultItemIdentifiers:(NSToolbar*)toolbar {
    return self.orderedIdentifiers;
}

- (NSArray<NSToolbarItemIdentifier>*)toolbarAllowedItemIdentifiers:(NSToolbar*)toolbar {
    return self.allowedIdentifiers ?: self.orderedIdentifiers;
}

- (NSToolbarItem*)toolbar:(NSToolbar*)toolbar
        itemForItemIdentifier:(NSToolbarItemIdentifier)itemIdentifier
    willBeInsertedIntoToolbar:(BOOL)flag {
    return self.itemsByIdentifier[itemIdentifier];
}

@end

static NSToolbarItem* toolbarItemForIdentifier(void* handlePtr, const char* identifier);

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

void* toolbarCreate(const char* identifier, bool customizable) {
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
    // Both flags are set before the delegate so AppKit restores a saved
    // layout for a customizable toolbar when it is attached.
    toolbar.allowsUserCustomization = customizable;
    toolbar.autosavesConfiguration = customizable;
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
    if (toolbar.autosavesConfiguration) {
        // AppKit restores the user's saved layout (or asks the delegate for
        // the default one) when the toolbar joins the window. Forcing the
        // default identifiers here would discard that saved layout.
    } else if (@available(macOS 15.0, *)) {
        toolbar.itemIdentifiers = identifiers;
    } else {
        for (NSToolbarItemIdentifier identifier in identifiers) {
            [toolbar insertItemWithItemIdentifier:identifier atIndex:toolbar.items.count];
        }
    }
    [delegate.knownIdentifiers addObjectsFromArray:identifiers];
    if (delegate.allowedIdentifiers != nil) {
        [delegate.knownIdentifiers addObjectsFromArray:delegate.allowedIdentifiers];
    }

    for (NSToolbarItem* item in delegate.itemsByIdentifier.allValues) {
        WailsToolbarShareTarget* shareTarget = objc_getAssociatedObject(item, WailsToolbarShareTargetAssociationKey);
        if (shareTarget != nil) shareTarget.window = window;
    }

    window.toolbar = toolbar;
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

void toolbarSetCustomizable(void* handlePtr, bool customizable) {
    WailsToolbarHandle* handle = toolbarHandle(handlePtr);
    if (handle == NULL || handle->toolbar == nil) return;
    handle->toolbar.allowsUserCustomization = customizable;
    handle->toolbar.autosavesConfiguration = customizable;
}

void toolbarRunCustomizationPalette(void* handlePtr) {
    WailsToolbarHandle* handle = toolbarHandle(handlePtr);
    if (handle == NULL || handle->toolbar == nil) return;
    if (!handle->toolbar.allowsUserCustomization || handle->toolbar.customizationPaletteIsRunning) return;
    [handle->toolbar runCustomizationPalette:nil];
}

bool toolbarSupportsCenteredItems(void) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 130000
    if (@available(macOS 13.0, *)) return true;
#endif
    return false;
}

void toolbarForgetItem(void* handlePtr, const char* identifier) {
    WailsToolbarDelegate* delegate = toolbarDelegate(handlePtr);
    if (delegate == nil || identifier == NULL) return;
    [delegate.itemsByIdentifier removeObjectForKey:[NSString stringWithUTF8String:identifier]];
}

// Resolves one layout entry to the identifier AppKit should receive, or nil
// when the entry is unavailable on this system (a standard identifier that
// does not exist yet, or a Wails item that was not built natively).
static NSToolbarItemIdentifier toolbarResolveLayoutEntry(WailsToolbarDelegate* delegate, id entry) {
    if (![entry isKindOfClass:[NSDictionary class]]) return nil;
    id identifier = ((NSDictionary*)entry)[@"id"];
    id kind = ((NSDictionary*)entry)[@"kind"];
    if (![kind isKindOfClass:[NSString class]]) kind = @"item";

    if ([kind isEqualToString:@"space"]) return NSToolbarSpaceItemIdentifier;
    if ([kind isEqualToString:@"flexibleSpace"]) return NSToolbarFlexibleSpaceItemIdentifier;
    if ([kind isEqualToString:@"sidebarToggle"]) return NSToolbarToggleSidebarItemIdentifier;
    if ([kind isEqualToString:@"sidebarSeparator"]) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
        if (@available(macOS 11.0, *)) return NSToolbarSidebarTrackingSeparatorItemIdentifier;
#endif
        return nil;
    }
    if ([kind isEqualToString:@"inspectorSeparator"]) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 140000
        if (@available(macOS 14.0, *)) return NSToolbarInspectorTrackingSeparatorItemIdentifier;
#endif
        return nil;
    }
    if ([kind isEqualToString:@"inspectorToggle"]) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 140000
        if (@available(macOS 14.0, *)) return NSToolbarToggleInspectorItemIdentifier;
#endif
        // Fall through to the Wails-owned fallback item registered under id.
    }
    if (![identifier isKindOfClass:[NSString class]]) return nil;
    return delegate.itemsByIdentifier[identifier] != nil ? identifier : nil;
}

static NSArray<NSToolbarItemIdentifier>* toolbarResolveLayoutEntries(WailsToolbarDelegate* delegate, id entries) {
    NSMutableArray<NSToolbarItemIdentifier>* result = [NSMutableArray array];
    if (![entries isKindOfClass:[NSArray class]]) return result;
    for (id entry in (NSArray*)entries) {
        NSToolbarItemIdentifier identifier = toolbarResolveLayoutEntry(delegate, entry);
        if (identifier != nil) [result addObject:identifier];
    }
    return result;
}

static BOOL toolbarIdentifierIsSpace(NSToolbarItemIdentifier identifier) {
    return [identifier isEqualToString:NSToolbarSpaceItemIdentifier] ||
        [identifier isEqualToString:NSToolbarFlexibleSpaceItemIdentifier];
}

// A live item matches a layout identifier when the identifiers agree and,
// for a Wails-owned item, the live instance is the one currently registered
// (a rebuilt group registers a new instance under the same identifier).
static BOOL toolbarLiveItemMatches(WailsToolbarDelegate* delegate, NSToolbarItem* live, NSToolbarItemIdentifier identifier) {
    if (![live.itemIdentifier isEqualToString:identifier]) return NO;
    NSToolbarItem* canonical = delegate.itemsByIdentifier[identifier];
    return canonical == nil || canonical == live;
}

static NSInteger toolbarLiveIndexOfIdentifier(NSToolbar* toolbar, NSToolbarItemIdentifier identifier) {
    NSInteger index = 0;
    for (NSToolbarItem* live in toolbar.items) {
        if ([live.itemIdentifier isEqualToString:identifier]) return index;
        index++;
    }
    return NSNotFound;
}

// Makes the live toolbar mirror the desired identifier list with the minimum
// of removals and insertions: matching prefixes are kept, an item found later
// is exposed by removing what precedes it, and a missing item is inserted.
static void toolbarApplyExactOrder(NSToolbar* toolbar, WailsToolbarDelegate* delegate,
    NSArray<NSToolbarItemIdentifier>* desired) {
    NSUInteger index = 0;
    for (NSToolbarItemIdentifier identifier in desired) {
        if (index < toolbar.items.count && toolbarLiveItemMatches(delegate, toolbar.items[index], identifier)) {
            index++;
            continue;
        }
        NSUInteger found = NSNotFound;
        for (NSUInteger candidate = index; candidate < toolbar.items.count; candidate++) {
            if (toolbarLiveItemMatches(delegate, toolbar.items[candidate], identifier)) {
                found = candidate;
                break;
            }
        }
        if (found != NSNotFound) {
            while (found > index) {
                [toolbar removeItemAtIndex:index];
                found--;
            }
            index++;
            continue;
        }
        [toolbar insertItemWithItemIdentifier:identifier atIndex:index];
        if (index < toolbar.items.count && toolbarLiveItemMatches(delegate, toolbar.items[index], identifier)) {
            index++;
        }
    }
    while (toolbar.items.count > index) {
        [toolbar removeItemAtIndex:toolbar.items.count - 1];
    }
}

// Position for identifier in the live toolbar that preserves the relative
// order of the default layout: the number of live items that precede it in
// the default list (items unknown to the default list do not count).
static NSUInteger toolbarLivePositionForIdentifier(NSToolbar* toolbar, NSArray<NSToolbarItemIdentifier>* defaults,
    NSToolbarItemIdentifier identifier, NSToolbarItem* ignoring) {
    NSUInteger rank = [defaults indexOfObject:identifier];
    if (rank == NSNotFound) return toolbar.items.count;
    NSUInteger position = 0;
    NSUInteger scanned = 0;
    for (NSToolbarItem* live in toolbar.items) {
        if (live == ignoring) continue;
        scanned++;
        NSUInteger liveRank = [defaults indexOfObject:live.itemIdentifier];
        if (liveRank != NSNotFound && liveRank < rank) position = scanned;
    }
    return position;
}

// Preserves the user's layout while applying structural changes from Go.
static void toolbarApplyCustomizableLayout(NSToolbar* toolbar, WailsToolbarDelegate* delegate,
    NSArray<NSToolbarItemIdentifier>* defaults, NSArray<NSToolbarItemIdentifier>* allowed,
    NSToolbarItemIdentifier moved) {
    NSSet<NSToolbarItemIdentifier>* allowedSet = [NSSet setWithArray:allowed];

    // Remove items that are no longer offered and replace stale instances of
    // rebuilt items in place.
    for (NSInteger index = (NSInteger)toolbar.items.count - 1; index >= 0; index--) {
        NSToolbarItem* live = toolbar.items[index];
        NSToolbarItemIdentifier identifier = live.itemIdentifier;
        if (![allowedSet containsObject:identifier] && !toolbarIdentifierIsSpace(identifier)) {
            [toolbar removeItemAtIndex:index];
            continue;
        }
        NSToolbarItem* canonical = delegate.itemsByIdentifier[identifier];
        if (canonical != nil && canonical != live) {
            [toolbar removeItemAtIndex:index];
            [toolbar insertItemWithItemIdentifier:identifier atIndex:index];
        }
    }

    // Insert items Go added since the last sync. Identifiers already known
    // but absent were removed by the user and stay removed.
    for (NSToolbarItemIdentifier identifier in defaults) {
        if ([delegate.knownIdentifiers containsObject:identifier]) continue;
        if (toolbarLiveIndexOfIdentifier(toolbar, identifier) != NSNotFound) continue;
        NSUInteger position = toolbarLivePositionForIdentifier(toolbar, defaults, identifier, nil);
        [toolbar insertItemWithItemIdentifier:identifier atIndex:MIN(position, toolbar.items.count)];
    }

    if (moved.length > 0) {
        NSInteger current = toolbarLiveIndexOfIdentifier(toolbar, moved);
        if (current != NSNotFound) {
            NSToolbarItem* live = toolbar.items[current];
            NSUInteger position = toolbarLivePositionForIdentifier(toolbar, defaults, moved, live);
            if (position != (NSUInteger)current) {
                [toolbar removeItemAtIndex:current];
                [toolbar insertItemWithItemIdentifier:moved atIndex:MIN(position, toolbar.items.count)];
            }
        }
    }
}

void toolbarSync(void* handlePtr, const char* layoutJSON) {
    WailsToolbarHandle* handle = toolbarHandle(handlePtr);
    WailsToolbarDelegate* delegate = toolbarDelegate(handlePtr);
    if (handle == NULL || handle->toolbar == nil || delegate == nil || layoutJSON == NULL) return;

    NSData* data = [[NSString stringWithUTF8String:layoutJSON] dataUsingEncoding:NSUTF8StringEncoding];
    id decoded = data == nil ? nil : [NSJSONSerialization JSONObjectWithData:data options:0 error:nil];
    if (![decoded isKindOfClass:[NSDictionary class]]) return;
    NSDictionary* layout = (NSDictionary*)decoded;

    BOOL customizable = [layout[@"customizable"] isKindOfClass:[NSNumber class]] &&
        ((NSNumber*)layout[@"customizable"]).boolValue;
    NSArray<NSToolbarItemIdentifier>* defaults = toolbarResolveLayoutEntries(delegate, layout[@"default"]);
    NSMutableArray<NSToolbarItemIdentifier>* allowed =
        [toolbarResolveLayoutEntries(delegate, layout[@"allowed"]) mutableCopy];
    NSArray<NSToolbarItemIdentifier>* centered = toolbarResolveLayoutEntries(delegate, layout[@"centered"]);
    id movedValue = layout[@"moved"];
    NSToolbarItemIdentifier moved = [movedValue isKindOfClass:[NSString class]] ? movedValue : nil;

    if (customizable) {
        // The palette always offers the standard spacers.
        if (![allowed containsObject:NSToolbarSpaceItemIdentifier]) [allowed addObject:NSToolbarSpaceItemIdentifier];
        if (![allowed containsObject:NSToolbarFlexibleSpaceItemIdentifier]) {
            [allowed addObject:NSToolbarFlexibleSpaceItemIdentifier];
        }
    }

    delegate.orderedIdentifiers = [[defaults mutableCopy] autorelease];
    delegate.allowedIdentifiers = allowed;
    [allowed release];

    NSToolbar* toolbar = handle->toolbar;
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 130000
    if (@available(macOS 13.0, *)) {
        toolbar.centeredItemIdentifiers = [NSSet setWithArray:centered];
    }
#endif

    // Detached: the lists above seed the initial layout in toolbarAttach.
    if (toolbar.delegate == nil) return;

    if (customizable) {
        toolbarApplyCustomizableLayout(toolbar, delegate, defaults, delegate.allowedIdentifiers, moved);
    } else {
        toolbarApplyExactOrder(toolbar, delegate, defaults);
    }
    [delegate.knownIdentifiers addObjectsFromArray:delegate.allowedIdentifiers];
    [toolbar validateVisibleItems];
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
            actual.enabled = source.enabled;
            actual.hidden = source.hidden;
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
    NSToolbarItem* item = nil;
    NSSearchField* field = nil;

    if (@available(macOS 11.0, *)) {
        Class searchToolbarItemClass = NSClassFromString(@"NSSearchToolbarItem");
        item = [[searchToolbarItemClass alloc] initWithItemIdentifier:identifierString];
        field = [item valueForKey:@"searchField"];
    } else {
        item = [[NSToolbarItem alloc] initWithItemIdentifier:identifierString];
        field = [[NSSearchField alloc] initWithFrame:NSMakeRect(0, 0, 220, 24)];
        item.view = field;
        [field release];
    }

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

void toolbarAddSpaceIdentifier(void* handlePtr) {
    WailsToolbarDelegate* delegate = toolbarDelegate(handlePtr);
    if (delegate != nil) {
        [delegate.orderedIdentifiers addObject:NSToolbarSpaceItemIdentifier];
    }
}

void* toolbarAddMenuItem(void* handlePtr, const char* identifier, const char* label,
    const char* symbolName, const char* tooltip, bool bordered, bool disabled, bool hidden,
    void* nsMenu, bool showsIndicator) {
    WailsToolbarDelegate* delegate = toolbarDelegate(handlePtr);
    if (delegate == nil) return NULL;

    NSString* identifierString = [NSString stringWithUTF8String:identifier];
    if (@available(macOS 10.15, *)) {
        NSMenuToolbarItem* item = [[NSMenuToolbarItem alloc] initWithItemIdentifier:identifierString];
        item.label = [NSString stringWithUTF8String:label];
        item.menu = (NSMenu*)nsMenu;
        item.showsIndicator = showsIndicator;
        NSImage* symbol = toolbarSymbolImage(symbolName, item.label);
        if (symbol != nil) {
            item.image = symbol;
        } else if (@available(macOS 11.0, *)) {
            // A menu item without an image renders as text; supply the
            // system chevron image so icon-only toolbars still show it.
            item.image = [NSImage imageWithSystemSymbolName:@"ellipsis.circle"
                                   accessibilityDescription:item.label];
        }
        applyCommonItemStyle(item, tooltip, bordered, false, disabled, hidden,
            false, 0, 0, 0, 0, 0);
        [delegate.orderedIdentifiers addObject:identifierString];
        delegate.itemsByIdentifier[identifierString] = item;
        [item release];
        return item;
    }
    NSLog(@"[Wails] toolbar menu item %s requires macOS 10.15 or newer and was omitted", label);
    return NULL;
}

void toolbarMenuItemSetShowsIndicator(void* handlePtr, const char* identifier, bool showsIndicator) {
    NSToolbarItem* item = toolbarItemForIdentifier(handlePtr, identifier);
    if (@available(macOS 10.15, *)) {
        if ([item isKindOfClass:[NSMenuToolbarItem class]]) {
            ((NSMenuToolbarItem*)item).showsIndicator = showsIndicator;
        }
    }
}

void toolbarItemSetVisibilityPriority(void* handlePtr, const char* identifier, int priority) {
    NSToolbarItem* item = toolbarItemForIdentifier(handlePtr, identifier);
    if (item != nil) item.visibilityPriority = priority;
}

void toolbarItemSetNavigational(void* handlePtr, const char* identifier, bool navigational) {
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
    NSToolbarItem* item = toolbarItemForIdentifier(handlePtr, identifier);
    if (@available(macOS 11.0, *)) {
        if (item != nil) item.navigational = navigational;
    }
#endif
}

static NSSearchField* toolbarSearchFieldForItem(NSToolbarItem* item) {
    if (item == nil) return nil;
    if (@available(macOS 11.0, *)) {
        Class searchToolbarItemClass = NSClassFromString(@"NSSearchToolbarItem");
        if (searchToolbarItemClass != nil && [item isKindOfClass:searchToolbarItemClass]) {
            return [item valueForKey:@"searchField"];
        }
    }
    if ([item.view isKindOfClass:[NSSearchField class]]) return (NSSearchField*)item.view;
    return nil;
}

void toolbarSearchItemSetPlaceholder(void* handlePtr, const char* identifier, const char* placeholder) {
    NSSearchField* field = toolbarSearchFieldForItem(toolbarItemForIdentifier(handlePtr, identifier));
    if (field == nil) return;
    field.placeholderString = placeholder != NULL && strlen(placeholder) > 0
        ? [NSString stringWithUTF8String:placeholder]
        : nil;
}

void toolbarSearchItemSetIncremental(void* handlePtr, const char* identifier, bool incremental) {
    NSSearchField* field = toolbarSearchFieldForItem(toolbarItemForIdentifier(handlePtr, identifier));
    if (field == nil) return;
    field.sendsSearchStringImmediately = incremental;
    field.sendsWholeSearchString = !incremental;
}

void toolbarSearchItemSetRecents(void* handlePtr, const char* identifier, const char* autosaveName, int maximumRecents) {
    NSSearchField* field = toolbarSearchFieldForItem(toolbarItemForIdentifier(handlePtr, identifier));
    if (field == nil) return;
    if (autosaveName == NULL || strlen(autosaveName) == 0) {
        field.recentsAutosaveName = nil;
        field.maximumRecents = 0;
        return;
    }
    field.recentsAutosaveName = [NSString stringWithUTF8String:autosaveName];
    field.maximumRecents = maximumRecents > 0 ? maximumRecents : 10;
}

// Appends the standard recent-searches section to menu unless it is already
// present. AppKit recognises the entries by tag.
static void toolbarAppendRecentsSection(NSMenu* menu) {
    if ([menu itemWithTag:NSSearchFieldRecentsMenuItemTag] != nil) return;
    if (menu.numberOfItems > 0) {
        NSMenuItem* separator = [NSMenuItem separatorItem];
        separator.tag = NSSearchFieldRecentsTitleMenuItemTag;
        [menu addItem:separator];
    }
    NSMenuItem* title = [[NSMenuItem alloc] initWithTitle:@"Recent Searches" action:nil keyEquivalent:@""];
    title.tag = NSSearchFieldRecentsTitleMenuItemTag;
    [menu addItem:title];
    [title release];

    NSMenuItem* recents = [[NSMenuItem alloc] initWithTitle:@"Recents" action:nil keyEquivalent:@""];
    recents.tag = NSSearchFieldRecentsMenuItemTag;
    [menu addItem:recents];
    [recents release];

    NSMenuItem* noRecents = [[NSMenuItem alloc] initWithTitle:@"No Recent Searches" action:nil keyEquivalent:@""];
    noRecents.tag = NSSearchFieldNoRecentsMenuItemTag;
    [menu addItem:noRecents];
    [noRecents release];

    NSMenuItem* clear = [[NSMenuItem alloc] initWithTitle:@"Clear Recent Searches" action:nil keyEquivalent:@""];
    clear.tag = NSSearchFieldClearRecentsMenuItemTag;
    [menu addItem:clear];
    [clear release];
}

void toolbarSearchItemSetMenu(void* handlePtr, const char* identifier, void* nsMenu, bool includeRecents) {
    NSSearchField* field = toolbarSearchFieldForItem(toolbarItemForIdentifier(handlePtr, identifier));
    if (field == nil) return;
    NSMenu* menu = (NSMenu*)nsMenu;
    if (menu == nil) {
        if (!includeRecents) {
            field.searchMenuTemplate = nil;
            return;
        }
        menu = [[[NSMenu alloc] initWithTitle:@""] autorelease];
    }
    if (includeRecents) toolbarAppendRecentsSection(menu);
    field.searchMenuTemplate = menu;
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

void toolbarShareItemSetProvider(void* handlePtr, const char* identifier, const char* providerJSON) {
    NSToolbarItem* item = toolbarItemForIdentifier(handlePtr, identifier);
    WailsToolbarShareTarget* target = objc_getAssociatedObject(item, WailsToolbarShareTargetAssociationKey);
    if (target != nil) toolbarShareTargetSetProvider(target, providerJSON);
}

void toolbarItemSetLabel(void* handlePtr, const char* identifier, const char* label) {
    NSToolbarItem* item = toolbarItemForIdentifier(handlePtr, identifier);
    if (item != nil) item.label = [NSString stringWithUTF8String:label];
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
