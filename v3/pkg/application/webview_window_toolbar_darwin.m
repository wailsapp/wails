//go:build darwin && !ios && !server

#import "webview_window_toolbar_darwin.h"
#import <objc/runtime.h>
#import <string.h>

@implementation WailsToolbarItem

- (void)handleClick {
    processToolbarItemClick(self.itemID);
}

@end

@implementation WailsSearchToolbarItem

- (void)handleSearch:(id)sender {
    NSSearchField* field = (NSSearchField*)sender;
    char* query = (char*)[field.stringValue UTF8String];
    processToolbarSearch(self.itemID, query);
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

static void applyCommonItemStyle(NSToolbarItem* item, const char* tooltip,
    bool bordered, bool prominent, bool disabled, bool hidden,
    bool hasTint, double tintR, double tintG, double tintB, double tintA, int badgeCount) {
    if (tooltip != NULL && strlen(tooltip) > 0) item.toolTip = [NSString stringWithUTF8String:tooltip];
    item.enabled = !disabled;
    item.hidden = hidden;
    if (@available(macOS 10.15, *)) {
        item.bordered = bordered;
    }
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 260000
    if (@available(macOS 26.0, *)) {
        if (prominent) {
            item.style = NSToolbarItemStyleProminent;
        }
        if (hasTint) {
            item.backgroundTintColor = [NSColor colorWithRed:tintR green:tintG blue:tintB alpha:tintA];
        }
        if (badgeCount > 0) {
            item.badge = [NSItemBadge badgeWithCount:badgeCount];
        }
    }
#endif
}

void* toolbarNewAndAttach(void* nsWindow) {
    NSWindow* window = (NSWindow*)nsWindow;
    WailsToolbarDelegate* delegate = [[[WailsToolbarDelegate alloc] init] autorelease];
    NSToolbar* toolbar = [[[NSToolbar alloc] initWithIdentifier:@"wails.toolbar"] autorelease];
    toolbar.delegate = delegate;
    toolbar.displayMode = NSToolbarDisplayModeIconAndLabel;
    // A toolbar can be restored as hidden from AppKit's per-window toolbar
    // state. Wails owns this toolbar configuration, so always make it
    // visible when attaching it.
    toolbar.allowsUserCustomization = NO;
    toolbar.autosavesConfiguration = NO;
    toolbar.visible = YES;
    // NSToolbar.delegate is a weak/assign reference: without this, the
    // delegate has no owner once this function returns and gets
    // deallocated, silently breaking every item lookup. Associate it with
    // RETAIN so it lives exactly as long as the toolbar does.
    objc_setAssociatedObject(toolbar, "wailsToolbarDelegate", delegate, OBJC_ASSOCIATION_RETAIN);
    window.toolbar = toolbar;
    toolbar.visible = YES;
    return (void*)delegate;
}

void toolbarReload(void* nsWindow, void* delegatePtr, int style) {
    NSWindow* window = (NSWindow*)nsWindow;
    NSToolbar* toolbar = window.toolbar;
    if (toolbar == nil) return;

    // AppKit may ask for defaultItemIdentifiers while the toolbar is first
    // attached. Reinstalling the delegate after the item tree is complete
    // makes it re-read the now-populated identifiers.
    toolbar.delegate = nil;
    toolbar.delegate = (WailsToolbarDelegate*)delegatePtr;
    WailsToolbarDelegate* delegate = (WailsToolbarDelegate*)delegatePtr;
    NSArray<NSToolbarItemIdentifier>* identifiers = delegate.orderedIdentifiers;
    if (@available(macOS 15.0, *)) {
        // Explicitly set the live item list. Default identifiers are only
        // consulted when a toolbar is first configured, which is too early
        // for Wails because the Go item tree is built after attachment.
        toolbar.itemIdentifiers = identifiers;
    } else {
        while (toolbar.items.count > 0) {
            [toolbar removeItemAtIndex:toolbar.items.count - 1];
        }
        for (NSToolbarItemIdentifier identifier in identifiers) {
            [toolbar insertItemWithItemIdentifier:identifier atIndex:toolbar.items.count];
        }
    }
    [toolbar retain];
    window.toolbar = nil;
    window.toolbar = toolbar;
    [toolbar release];
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 110000
    if (@available(macOS 11.0, *)) {
        window.toolbarStyle = style;
    }
#endif
    [toolbar validateVisibleItems];
    toolbar.visible = YES;
}

void toolbarDetach(void* nsWindow) {
    NSWindow* window = (NSWindow*)nsWindow;
    window.toolbar = nil;
}

// Returns a +1 retained item: the caller (toolbarAddButtonItem, or
// toolbarAddGroupItem via the member array) is responsible for handing it
// to something that retains it (an NSDictionary/NSArray) and then
// releasing this +1.
void* toolbarBuildButtonItemStandalone(const char* identifier, unsigned int itemID,
    const char* label, const char* symbolName, const char* tooltip,
    bool bordered, bool disabled, bool hidden) {
    NSString* identifierStr = [NSString stringWithUTF8String:identifier];
    WailsToolbarItem* item = [[WailsToolbarItem alloc] initWithItemIdentifier:identifierStr];
    item.itemID = itemID;
    item.label = [NSString stringWithUTF8String:label];
    item.target = item;
    item.action = @selector(handleClick);
    item.enabled = !disabled;
    item.hidden = hidden;
    if (tooltip != NULL && strlen(tooltip) > 0) item.toolTip = [NSString stringWithUTF8String:tooltip];
    if (@available(macOS 10.15, *)) {
        item.bordered = bordered;
    }
    if (symbolName != NULL && strlen(symbolName) > 0) {
        if (@available(macOS 11.0, *)) {
            NSString* symbolStr = [NSString stringWithUTF8String:symbolName];
            item.image = [NSImage imageWithSystemSymbolName:symbolStr accessibilityDescription:item.label];
        }
    }
    return (void*)item; // +1, intentionally not autoreleased, see callers
}

void* toolbarAddButtonItem(void* delegatePtr, const char* identifier, unsigned int itemID,
    const char* label, const char* symbolName, const char* tooltip,
    bool bordered, bool prominent, bool disabled, bool hidden,
    bool hasTint, double tintR, double tintG, double tintB, double tintA, int badgeCount) {
    WailsToolbarDelegate* delegate = (WailsToolbarDelegate*)delegatePtr;
    WailsToolbarItem* item = (WailsToolbarItem*)toolbarBuildButtonItemStandalone(
        identifier, itemID, label, symbolName, tooltip, bordered, disabled, hidden); // +1
    applyCommonItemStyle(item, tooltip, bordered, prominent, disabled, hidden,
        hasTint, tintR, tintG, tintB, tintA, badgeCount);

    NSString* identifierStr = [NSString stringWithUTF8String:identifier];
    [delegate.orderedIdentifiers addObject:identifierStr];
    delegate.itemsByIdentifier[identifierStr] = item; // dictionary takes its own retain
    [item release]; // release our +1; the dictionary now owns the only reference
    return (void*)item; // safe: dictionary keeps it alive; caller must not release
}

void* toolbarAddGroupItem(void* delegatePtr, const char* identifier,
    const char* label, void** memberItems, int memberCount, int selectionMode, int selectedIndex) {
    WailsToolbarDelegate* delegate = (WailsToolbarDelegate*)delegatePtr;

    NSMutableArray<NSToolbarItem*>* subitems = [NSMutableArray arrayWithCapacity:memberCount];
    for (int i = 0; i < memberCount; i++) {
        NSToolbarItem* memberItem = (NSToolbarItem*)memberItems[i]; // +1 from the standalone builder
        [subitems addObject:memberItem]; // array takes its own retain
        delegate.itemsByIdentifier[memberItem.itemIdentifier] = memberItem;
        [memberItem release]; // release our +1
    }

    NSString* identifierStr = [NSString stringWithUTF8String:identifier];
    NSToolbarItemGroup* group = [[NSToolbarItemGroup alloc] initWithItemIdentifier:identifierStr];
    // subitems is a plain property (there's no subitems-taking initializer);
    // each subitem already carries its own target/action from
    // toolbarBuildButtonItemStandalone, so clicking a segment dispatches
    // that segment's own itemID directly, no group-level target/action.
    group.subitems = subitems;
    group.label = [NSString stringWithUTF8String:label];
    if (@available(macOS 10.15, *)) {
        switch (selectionMode) {
            case 1: group.selectionMode = NSToolbarItemGroupSelectionModeMomentary; break;
            case 2: group.selectionMode = NSToolbarItemGroupSelectionModeSelectAny; break;
            default: group.selectionMode = NSToolbarItemGroupSelectionModeSelectOne; break;
        }
        if (selectedIndex >= 0 && selectedIndex < memberCount) group.selectedIndex = selectedIndex;
        group.controlRepresentation = NSToolbarItemGroupControlRepresentationAutomatic;
    }

    [delegate.orderedIdentifiers addObject:identifierStr];
    delegate.itemsByIdentifier[identifierStr] = group; // dictionary retains
    [group release]; // release our +1 from alloc
    return (void*)group;
}

void* toolbarAddSearchItem(void* delegatePtr, const char* identifier, unsigned int itemID,
    const char* label, const char* tooltip, bool disabled, bool hidden) {
    WailsToolbarDelegate* delegate = (WailsToolbarDelegate*)delegatePtr;
    NSString* identifierStr = [NSString stringWithUTF8String:identifier];

    WailsSearchToolbarItem* item = [[WailsSearchToolbarItem alloc] initWithItemIdentifier:identifierStr];
    item.itemID = itemID;
    item.label = [NSString stringWithUTF8String:label];
    item.searchField.target = item;
    item.searchField.action = @selector(handleSearch:);
    item.searchField.sendsWholeSearchString = YES;
    item.enabled = !disabled;
    item.hidden = hidden;
    if (tooltip != NULL && strlen(tooltip) > 0) item.toolTip = [NSString stringWithUTF8String:tooltip];

    [delegate.orderedIdentifiers addObject:identifierStr];
    delegate.itemsByIdentifier[identifierStr] = item; // dictionary retains
    [item release]; // release our +1 from alloc
    return (void*)item;
}

void toolbarAddFlexibleSpaceIdentifier(void* delegatePtr) {
    WailsToolbarDelegate* delegate = (WailsToolbarDelegate*)delegatePtr;
    [delegate.orderedIdentifiers addObject:NSToolbarFlexibleSpaceItemIdentifier];
}

void toolbarAddSidebarToggleIdentifier(void* delegatePtr, const char* identifier) {
    WailsToolbarDelegate* delegate = (WailsToolbarDelegate*)delegatePtr;
    NSString* identifierStr = [NSString stringWithUTF8String:identifier];

    WailsToolbarItem* item = [[WailsToolbarItem alloc] initWithItemIdentifier:identifierStr];
    item.itemID = 0; // not user-dispatched; standard AppKit action below
    item.label = @"Sidebar";
    if (@available(macOS 11.0, *)) {
        item.image = [NSImage imageWithSystemSymbolName:@"sidebar.left" accessibilityDescription:@"Toggle Sidebar"];
    }
    item.target = nil; // let the responder chain find toggleSidebar: on the NSSplitViewController
    item.action = @selector(toggleSidebar:);

    [delegate.orderedIdentifiers addObject:identifierStr];
    delegate.itemsByIdentifier[identifierStr] = item; // dictionary retains
    [item release]; // release our +1 from alloc
}

void toolbarAddInspectorToggleIdentifier(void* delegatePtr, const char* identifier) {
    WailsToolbarDelegate* delegate = (WailsToolbarDelegate*)delegatePtr;
    NSString* identifierStr = [NSString stringWithUTF8String:identifier];
    WailsToolbarItem* item = [[WailsToolbarItem alloc] initWithItemIdentifier:identifierStr];
    item.itemID = 0;
    item.label = @"Inspector";
    if (@available(macOS 11.0, *)) {
        item.image = [NSImage imageWithSystemSymbolName:@"sidebar.right" accessibilityDescription:@"Toggle Inspector"];
    }
    item.target = nil;
    item.action = @selector(toggleInspector:);
    [delegate.orderedIdentifiers addObject:identifierStr];
    delegate.itemsByIdentifier[identifierStr] = item;
    [item release];
}

static NSToolbarItem* toolbarItemForIdentifier(void* delegatePtr, const char* identifier) {
    WailsToolbarDelegate* delegate = (WailsToolbarDelegate*)delegatePtr;
    return delegate.itemsByIdentifier[[NSString stringWithUTF8String:identifier]];
}

void toolbarItemSetLabel(void* delegatePtr, const char* identifier, const char* label) {
    NSToolbarItem* item = toolbarItemForIdentifier(delegatePtr, identifier);
    if (item != nil) item.label = [NSString stringWithUTF8String:label];
}

void toolbarItemSetEnabled(void* delegatePtr, const char* identifier, bool enabled) {
    NSToolbarItem* item = toolbarItemForIdentifier(delegatePtr, identifier);
    if (item != nil) item.enabled = enabled;
}

void toolbarItemSetHidden(void* delegatePtr, const char* identifier, bool hidden) {
    NSToolbarItem* item = toolbarItemForIdentifier(delegatePtr, identifier);
    if (item != nil) item.hidden = hidden;
}

void toolbarItemSetBadgeCount(void* delegatePtr, const char* identifier, int badgeCount) {
    NSToolbarItem* item = toolbarItemForIdentifier(delegatePtr, identifier);
#if MAC_OS_X_VERSION_MAX_ALLOWED >= 260000
    if (@available(macOS 26.0, *)) {
        item.badge = badgeCount > 0 ? [NSItemBadge badgeWithCount:badgeCount] : nil;
    }
#endif
}

void toolbarGroupSetSelectedIndex(void* delegatePtr, const char* identifier, int index) {
    NSToolbarItem* item = toolbarItemForIdentifier(delegatePtr, identifier);
    if ([item isKindOfClass:[NSToolbarItemGroup class]]) {
        NSToolbarItemGroup* group = (NSToolbarItemGroup*)item;
        if (index >= 0 && index < (int)group.subitems.count) group.selectedIndex = index;
    }
}
